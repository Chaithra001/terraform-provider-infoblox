package infoblox

import (
	"encoding/json"
	"fmt"
	"github.com/infobloxopen/infoblox-go-client/v2/utils"
	"net"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	ibclient "github.com/infobloxopen/infoblox-go-client/v2"
)

func testAccCheckAAAARecordDestroy(s *terraform.State) error {
	meta := testAccProvider.Meta()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "resource_aaaa_record" {
			continue
		}
		connector := meta.(ibclient.IBConnector)
		objMgr := ibclient.NewObjectManager(connector, "terraform_test", "test")
		rec, _ := objMgr.GetAAAARecordByRef(rs.Primary.ID)
		if rec != nil {
			return fmt.Errorf("record not found")
		}

	}
	return nil
}

func testAccAAAARecordCompare(
	t *testing.T,
	resPath string,
	expectedRec *ibclient.RecordAAAA,
	notExpectedIpAddr string,
	expectedCidr string,
	expectedFilterParams string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		res, found := s.RootModule().Resources[resPath]
		if !found {
			return fmt.Errorf("not found: %s", resPath)
		}
		internalId := res.Primary.Attributes["internal_id"]
		if internalId == "" {
			return fmt.Errorf("ID is not set")
		}

		ref, found := res.Primary.Attributes["ref"]
		if !found {
			return fmt.Errorf("'ref' attribute is not set")
		}

		connector := testAccProvider.Meta().(ibclient.IBConnector)
		objMgr := ibclient.NewObjectManager(
			connector,
			"terraform_test",
			"test")
		recAAAA, err := objMgr.SearchObjectByAltId("AAAA", ref, internalId, eaNameForInternalId)
		if err != nil {
			if isNotFoundError(err) {
				if expectedRec == nil {
					return nil
				}
				return fmt.Errorf("object with Terraform ID '%s' not found, but expected to exist", internalId)
			}
		}
		// Assertion of object type and error handling
		var rec *ibclient.RecordAAAA
		recJson, _ := json.Marshal(recAAAA)
		err = json.Unmarshal(recJson, &rec)

		if rec.Name == nil {
			return fmt.Errorf("'fqdn' is expected to be defined but it is not")
		}
		if *rec.Name != *expectedRec.Name {
			return fmt.Errorf(
				"'fqdn' does not match: got '%s', expected '%s'",
				*rec.Name,
				*expectedRec.Name)
		}
		if rec.Ipv6Addr == nil {
			return fmt.Errorf("'ipv6addr' is expected to be defined but it is not")
		}
		if notExpectedIpAddr != "" && notExpectedIpAddr == *rec.Ipv6Addr {
			return fmt.Errorf(
				"'ipv6_addr' field has value '%s' but that is not expected to happen",
				notExpectedIpAddr)
		}
		if expectedCidr != "" {
			_, parsedCidr, err := net.ParseCIDR(expectedCidr)
			if err != nil {
				panic(fmt.Sprintf("cannot parse CIDR '%s': %s", expectedCidr, err))
			}

			if !parsedCidr.Contains(net.ParseIP(*rec.Ipv6Addr)) {
				return fmt.Errorf(
					"IP address '%s' does not belong to the expected CIDR '%s'",
					*rec.Ipv6Addr, expectedCidr)
			}
		}
		if expectedRec.Ipv6Addr != nil {
			if *expectedRec.Ipv6Addr == "" {
				expectedRec.Ipv6Addr = utils.StringPtr(res.Primary.Attributes["ipv6_addr"])
			}
			if *rec.Ipv6Addr != *expectedRec.Ipv6Addr {
				return fmt.Errorf(
					"'ipv4address' does not match: got '%s', expected '%s'",
					*rec.Ipv6Addr, *expectedRec.Ipv6Addr)
			}
		}
		if rec.View != expectedRec.View {
			return fmt.Errorf(
				"'dns_view' does not match: got '%s', expected '%s'",
				rec.View, expectedRec.View)
		}
		if rec.UseTtl != nil {
			if expectedRec.UseTtl == nil {
				return fmt.Errorf("'use_ttl' is expected to be undefined but it is not")
			}
			if *rec.UseTtl != *expectedRec.UseTtl {
				return fmt.Errorf(
					"'use_ttl' does not match: got '%t', expected '%t'",
					*rec.UseTtl, *expectedRec.UseTtl)
			}
			if *rec.UseTtl {
				if *rec.Ttl != *expectedRec.Ttl {
					return fmt.Errorf(
						"'TTL' usage does not match: got '%d', expected '%d'",
						rec.Ttl, expectedRec.Ttl)
				}
			}
		}

		if rec.Comment != nil {
			if expectedRec.Comment == nil {
				return fmt.Errorf("'comment' is expected to be undefined but it is not")
			}
			if *rec.Comment != *expectedRec.Comment {
				return fmt.Errorf(
					"'comment' does not match: got '%s', expected '%s'",
					*rec.Comment, *expectedRec.Comment)
			}
		} else if expectedRec.Comment != nil {
			return fmt.Errorf("'comment' is expected to be defined but it is not")
		}

		actualFilterParams, exists := res.Primary.Attributes["filter_params"]
		if expectedFilterParams != "" {
			if !exists {
				return fmt.Errorf("'filter_params' is expected to be defined but it is not")
			}
			if actualFilterParams != expectedFilterParams {
				return fmt.Errorf(
					"'filter_params' does not match: got '%s', expected '%s'",
					actualFilterParams, expectedFilterParams)
			}
		} else if exists {
			return fmt.Errorf("'filter_params' is expected to be undefined but it is not")
		}

		return validateEAs(rec.Ea, expectedRec.Ea)
	}
}

var (
	regexpRequiredMissingIPv6    = regexp.MustCompile("any one of 'ipv6_addr', 'cidr' and 'filter_params' values is required")
	regexpCidrIpAddrConflictIPv6 = regexp.MustCompile("only one of 'ipv6_addr', 'cidr' and 'filter_params' values is allowed to be defined")
	regexpUpdateConflictIPv6     = regexp.MustCompile("only one of 'ipv6_addr' and 'cidr' values is allowed to update")
)

func TestAccResourceAAAARecord(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAAAARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "infoblox_zone_auth" "zone" {
						fqdn = "test.com"
						view = "default"
					}
					resource "infoblox_aaaa_record" "foo"{
						fqdn = "name1.test.com"
						ipv6_addr = "2000::1"
						cidr = "2000:1fde::/96"
                        network_view = "default"
					}`),
				ExpectError: regexpCidrIpAddrConflictIPv6,
			},
			{
				Config: fmt.Sprintf(`
					resource "infoblox_zone_auth" "zone" {
						fqdn = "test.com"
						view = "default"
					}
					resource "infoblox_aaaa_record" "foo"{
						fqdn = "name1.test.com"
					}`),
				ExpectError: regexpRequiredMissingIPv6,
			},
			{
				Config: fmt.Sprintf(`
					resource "infoblox_zone_auth" "zone" {
						fqdn = "test.com"
						view = "default"
					}
					resource "infoblox_aaaa_record" "foo"{
						fqdn = "name1.test.com"
						ipv6_addr = "2000::1"
						depends_on = [infoblox_zone_auth.zone]
					}`),
				Check: resource.ComposeTestCheckFunc(
					testAccAAAARecordCompare(t, "infoblox_aaaa_record.foo", &ibclient.RecordAAAA{
						Ipv6Addr: utils.StringPtr("2000::1"),
						Name:     utils.StringPtr("name1.test.com"),
						View:     "default",
						Ttl:      utils.Uint32Ptr(0),
						UseTtl:   utils.BoolPtr(false),
						Comment:  nil,
						Ea:       nil,
					}, "", "", ""),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "infoblox_dns_view" "view1" {
						name = "nondefault_view"
					}
					resource "infoblox_zone_auth" "zone1" {
						fqdn = "test.com"
						view = "nondefault_view"
						depends_on = [infoblox_dns_view.view1]
					}
					resource "infoblox_aaaa_record" "foo2"{
						fqdn = "name2.test.com"
						ipv6_addr = "2002::10"
						ttl = 10
						dns_view = "nondefault_view"
						comment = "test comment 1"
						ext_attrs = jsonencode({
						  "Location" = "New York"
						  "Site" = "HQ"
						})
						depends_on = [infoblox_zone_auth.zone1, infoblox_dns_view.view1]
					}`),
				Check: resource.ComposeTestCheckFunc(
					testAccAAAARecordCompare(t, "infoblox_aaaa_record.foo2", &ibclient.RecordAAAA{
						Ipv6Addr: utils.StringPtr("2002::10"),
						Name:     utils.StringPtr("name2.test.com"),
						View:     "nondefault_view",
						Ttl:      utils.Uint32Ptr(10),
						UseTtl:   utils.BoolPtr(true),
						Comment:  utils.StringPtr("test comment 1"),
						Ea: ibclient.EA{
							"Location": "New York",
							"Site":     "HQ",
						},
					}, "", "", ""),
				),
			},
			{
				Config: fmt.Sprintf(`
							resource "infoblox_dns_view" "view1" {
								name = "nondefault_view"
							}
							resource "infoblox_zone_auth" "zone1" {
								fqdn = "test.com"
								view = "nondefault_view"
								depends_on = [infoblox_dns_view.view1]
							}
							resource "infoblox_aaaa_record" "foo2"{
								fqdn = "name3.test.com"
								ipv6_addr = "2000::1"
								ttl = 155
								dns_view = "nondefault_view"
								comment = "test comment 2"
								depends_on = [infoblox_zone_auth.zone1, infoblox_dns_view.view1]
							}`),
				Check: resource.ComposeTestCheckFunc(
					testAccAAAARecordCompare(t, "infoblox_aaaa_record.foo2", &ibclient.RecordAAAA{
						Ipv6Addr: utils.StringPtr("2000::1"),
						Name:     utils.StringPtr("name3.test.com"),
						View:     "nondefault_view",
						Ttl:      utils.Uint32Ptr(155),
						UseTtl:   utils.BoolPtr(true),
						Comment:  utils.StringPtr("test comment 2"),
					}, "", "", ""),
				),
			},
			{
				Config: fmt.Sprintf(`
							resource "infoblox_dns_view" "view1" {
								name = "nondefault_view"
							}
							resource "infoblox_zone_auth" "zone1" {
								fqdn = "test.com"
								view = "nondefault_view"
								depends_on = [infoblox_dns_view.view1]
							}
							resource "infoblox_aaaa_record" "foo2"{
								fqdn = "name3.test.com"
								ipv6_addr = "2000::1"
								dns_view = "nondefault_view"
								depends_on = [infoblox_zone_auth.zone1, infoblox_dns_view.view1]
							}`),
				Check: resource.ComposeTestCheckFunc(
					testAccAAAARecordCompare(t, "infoblox_aaaa_record.foo2", &ibclient.RecordAAAA{
						Ipv6Addr: utils.StringPtr("2000::1"),
						Name:     utils.StringPtr("name3.test.com"),
						View:     "nondefault_view",
						UseTtl:   utils.BoolPtr(false),
					}, "", "", ""),
				),
			},
			{
				Config: fmt.Sprintf(`
							resource "infoblox_dns_view" "view1" {
								name = "nondefault_view"
							}
							resource "infoblox_zone_auth" "zone1" {
								fqdn = "test.com"
								view = "nondefault_view"
								depends_on = [infoblox_dns_view.view1]
							}
						    resource "infoblox_ipv6_network" "net1" {
							    cidr = "2000:1fde::/96"
							    network_view = "default"
						    }
							resource "infoblox_aaaa_record" "foo2"{
								fqdn = "name3.test.com"
				               cidr = infoblox_ipv6_network.net1.cidr
				               network_view = infoblox_ipv6_network.net1.network_view
								dns_view = "nondefault_view"
								depends_on = [infoblox_zone_auth.zone1, infoblox_dns_view.view1, infoblox_ipv6_network.net1]
							}`),
				Check: resource.ComposeTestCheckFunc(
					testAccAAAARecordCompare(t, "infoblox_aaaa_record.foo2", &ibclient.RecordAAAA{
						Name:   utils.StringPtr("name3.test.com"),
						View:   "nondefault_view",
						UseTtl: utils.BoolPtr(false),
					}, "2000::1", "2000:1fde::/96", ""),
				),
			},
			{
				Config: fmt.Sprintf(`
							resource "infoblox_dns_view" "view1" {
								name = "nondefault_view"
							}
							resource "infoblox_zone_auth" "zone1" {
								fqdn = "test.com"
								view = "nondefault_view"
								depends_on = [infoblox_dns_view.view1]
							}
				            resource "infoblox_ipv6_network" "netA" {
				                cidr = "2000:1fcc::/96"
				                network_view = "default"
				            }
							resource "infoblox_aaaa_record" "foo2"{
								fqdn = "name3.test.com"
				               	cidr = infoblox_ipv6_network.netA.cidr
								ipv6_addr = "2002::4"
				               	network_view = infoblox_ipv6_network.netA.network_view
								dns_view = "nondefault_view"
								depends_on = [infoblox_zone_auth.zone1, infoblox_dns_view.view1, infoblox_ipv6_network.netA]
							}`),
				ExpectError: regexpUpdateConflictIPv6,
			},
			{
				Config: fmt.Sprintf(`
							resource "infoblox_dns_view" "view1" {
								name = "nondefault_view"
							}
							resource "infoblox_zone_auth" "zone1" {
								fqdn = "test.com"
								view = "nondefault_view"
								depends_on = [infoblox_dns_view.view1]
							}
				            resource "infoblox_ipv6_network" "netA" {
				                cidr = "2000:1fcc::/96"
				                network_view = "default"
				            }
							resource "infoblox_aaaa_record" "foo2"{
								fqdn = "name3.test.com"
				                cidr = infoblox_ipv6_network.netA.cidr
				                network_view = infoblox_ipv6_network.netA.network_view
								dns_view = "nondefault_view"
								depends_on = [infoblox_zone_auth.zone1, infoblox_dns_view.view1, infoblox_ipv6_network.netA]
							}`),
				Check: resource.ComposeTestCheckFunc(
					testAccAAAARecordCompare(t, "infoblox_aaaa_record.foo2", &ibclient.RecordAAAA{
						Name:   utils.StringPtr("name3.test.com"),
						View:   "nondefault_view",
						UseTtl: utils.BoolPtr(false),
					}, "", "2000:1fcc::/96", ""),
				),
			},
			{
				Config: fmt.Sprintf(`
							resource "infoblox_network_view" "netview1" {
								name = "nondefault_netview"
							}
							resource "infoblox_dns_view" "view2" {
								name = "nondefault_view"
								network_view = "nondefault_view"
								depends_on = [infoblox_network_view.netview1]
			
							}
							resource "infoblox_zone_auth" "zone1" {
								fqdn = "test.com"
								view = "nondefault_view"
								depends_on = [infoblox_dns_view.view2]
							}
			
				            resource "infoblox_ipv6_network" "net3" {
				                cidr = "2000:1fcd::/96"
				                network_view = infoblox_network_view.netview1.name
				            }
							resource "infoblox_aaaa_record" "foo2"{
								fqdn = "name3.test.com"
				                cidr = infoblox_ipv6_network.net3.cidr
				                network_view = infoblox_ipv6_network.net3.network_view
								dns_view = "nondefault_view"
							}`),
				ExpectError: regexpNetviewUpdateNotAllowed,
			},
			{
				Config: fmt.Sprintf(`
							resource "infoblox_dns_view" "view1" {
								name = "nondefault_view"
							}
							resource "infoblox_zone_auth" "zone1" {
								fqdn = "test.com"
								view = "nondefault_view"
								depends_on = [infoblox_dns_view.view1]
							}
							resource "infoblox_aaaa_record" "foo2"{
								fqdn = "name3.test.com"
								ipv6_addr = "2000::2"
								dns_view = "nondefault_view"
								depends_on = [infoblox_zone_auth.zone1, infoblox_dns_view.view1]
							}`),
				Check: resource.ComposeTestCheckFunc(
					testAccAAAARecordCompare(t, "infoblox_aaaa_record.foo2", &ibclient.RecordAAAA{
						Ipv6Addr: utils.StringPtr("2000::2"),
						Name:     utils.StringPtr("name3.test.com"),
						View:     "nondefault_view",
						UseTtl:   utils.BoolPtr(false),
					}, "", "", ""),
				),
			},
			{
				Config: fmt.Sprintf(`
							resource "infoblox_zone_auth" "zone" {
								fqdn = "test.com"
								view = "default"
							}
							resource "infoblox_aaaa_record" "foo2"{
								fqdn = "name3.test.com"
								ipv6_addr = "2000::2"
								dns_view = "default"
								depends_on = [infoblox_zone_auth.zone]
							}`),
				ExpectError: regexpDnsviewUpdateNotAllowed,
			},
			{
				Config: fmt.Sprintf(`
					resource "infoblox_zone_auth" "zone1" {
						fqdn = "test1.com"
					}
					resource "infoblox_ipv6_network" "net2" {
						cidr = "2001:db8:abcd::/64"
						ext_attrs = jsonencode({
 							"Site" = "Blr"
						})
					}
					resource "infoblox_aaaa_record" "rec4" {
						fqdn = "dynamic.test1.com"
						filter_params = jsonencode({
						"*Site" = "Blr"})
						dns_view = "default"
					}`),
				Check: resource.ComposeTestCheckFunc(
					testAccAAAARecordCompare(t, "infoblox_aaaa_record.rec4", &ibclient.RecordAAAA{
						Name:   utils.StringPtr("dynamic.test1.com"),
						View:   "default",
						UseTtl: utils.BoolPtr(false),
					}, "", "", `{"*Site":"Blr"}`),
				),
			},
			{
				Config: fmt.Sprintf(`
					resource "infoblox_zone_auth" "zone1" {
						fqdn = "test1.com"
					}
					resource "infoblox_aaaa_record" "rec1" {
						fqdn = "missing_fields.test1.com"
						comment = "missing required fields"
						dns_view = "default"
						depends_on = [infoblox_zone_auth.zone1]
					}`),
				ExpectError: regexpRequiredMissingIPv6,
			},
		},
	})
}

func TestAcc_resourceAAAARecord_ea_inheritance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAAAARecordDestroy,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "infoblox_zone_auth" "zone" {
					fqdn = "test.com"
				}
				resource "infoblox_aaaa_record" "foo3"{
					dns_view = "default"
					fqdn = "testname2.test.com"
					ipv6_addr = "2002::4"
					comment = "test comment on AAAA record"
					ext_attrs = jsonencode({
						"Location" = "test AAAA location"
					})
					depends_on = [infoblox_zone_auth.zone]
				}`,
				Check: testAccAAAARecordCompare(t, "infoblox_aaaa_record.foo3", &ibclient.RecordAAAA{
					Ipv6Addr: utils.StringPtr("2002::4"),
					Name:     utils.StringPtr("testname2.test.com"),
					View:     "default",
					UseTtl:   utils.BoolPtr(false),
					Comment:  utils.StringPtr("test comment on AAAA record"),
					Ea: ibclient.EA{
						"Location": "test AAAA location",
					},
				}, "", "", ""),
			},
			// When extensible attributes are added by another tool,
			// terraform shouldn't remove those EAs
			{
				PreConfig: func() {
					conn := testAccProvider.Meta().(ibclient.IBConnector)

					n := &ibclient.RecordAAAA{}
					n.SetReturnFields(append(n.ReturnFields(), "extattrs"))

					qp := ibclient.NewQueryParams(
						false,
						map[string]string{
							"name":     "testname2.test.com",
							"ipv6addr": "2002::4",
						},
					)
					var res []ibclient.RecordAAAA
					err := conn.GetObject(n, "", qp, &res)
					if err != nil {
						panic(err)
					}

					res[0].View = ""
					res[0].Ea["Site"] = "Testing Site"

					_, err = conn.UpdateObject(&res[0], res[0].Ref)
					if err != nil {
						panic(err)
					}
				},
				Config: `
				resource "infoblox_zone_auth" "zone" {
					fqdn = "test.com"
				}
				resource "infoblox_aaaa_record" "foo3"{
					dns_view = "default"
					fqdn = "testname2.test.com"
					ipv6_addr = "2002::4"
					comment = "test comment on AAAA record"
					ext_attrs = jsonencode({
						"Location" = "test AAAA location"
					})
					depends_on = [infoblox_zone_auth.zone]
				}`,
				Check: resource.ComposeTestCheckFunc(
					// Resource object shouldn't have Site EA, since it's omitted by provider
					resource.TestCheckResourceAttr(
						"infoblox_aaaa_record.foo3", "ext_attrs",
						`{"Location":"test AAAA location"}`,
					),
					// Actual API object should have Site EA
					testAccAAAARecordCompare(t, "infoblox_aaaa_record.foo3", &ibclient.RecordAAAA{
						Ipv6Addr: utils.StringPtr("2002::4"),
						Name:     utils.StringPtr("testname2.test.com"),
						View:     "default",
						UseTtl:   utils.BoolPtr(false),
						Comment:  utils.StringPtr("test comment on AAAA record"),
						Ea: ibclient.EA{
							"Location": "test AAAA location",
							"Site":     "Testing Site",
						},
					}, "", "", ""),
				),
			},
			// Validate that inherited EA won't be removed if some field is updated in the resource
			{
				Config: `
				resource "infoblox_zone_auth" "zone" {
					fqdn = "test.com"
				}
				resource "infoblox_aaaa_record" "foo3"{
					dns_view = "default"
					fqdn = "testname2.test.com"
					ipv6_addr = "2002::4"
					comment = "updated comment on AAAA record"
					ext_attrs = jsonencode({
						"Location" = "test AAAA location"
					})
					depends_on = [infoblox_zone_auth.zone]
				}`,
				Check: testAccAAAARecordCompare(t, "infoblox_aaaa_record.foo3", &ibclient.RecordAAAA{
					Ipv6Addr: utils.StringPtr("2002::4"),
					Name:     utils.StringPtr("testname2.test.com"),
					View:     "default",
					UseTtl:   utils.BoolPtr(false),
					Comment:  utils.StringPtr("updated comment on AAAA record"),
					Ea: ibclient.EA{
						"Location": "test AAAA location",
						"Site":     "Testing Site",
					},
				}, "", "", ""),
			},
			// Validate that inherited EA can be updated
			{
				Config: `
				resource "infoblox_zone_auth" "zone" {
					fqdn = "test.com"
				}
				resource "infoblox_aaaa_record" "foo3"{
					dns_view = "default"
					fqdn = "testname2.test.com"
					ipv6_addr = "2002::4"
					comment = "test comment on AAAA record"
					ext_attrs = jsonencode({
						"Location" = "test AAAA location"
						"Site" = "New Testing Site"
					})
					depends_on = [infoblox_zone_auth.zone]
				}`,
				Check: testAccAAAARecordCompare(t, "infoblox_aaaa_record.foo3", &ibclient.RecordAAAA{
					Ipv6Addr: utils.StringPtr("2002::4"),
					Name:     utils.StringPtr("testname2.test.com"),
					View:     "default",
					UseTtl:   utils.BoolPtr(false),
					Comment:  utils.StringPtr("test comment on AAAA record"),
					Ea: ibclient.EA{
						"Location": "test AAAA location",
						"Site":     "New Testing Site",
					},
				}, "", "", ""),
			},
			// Validate that inherited EA can be removed, if updated
			{
				Config: `
				resource "infoblox_zone_auth" "zone" {
					fqdn = "test.com"
				}
				resource "infoblox_aaaa_record" "foo3"{
					dns_view = "default"
					fqdn = "testname2.test.com"
					ipv6_addr = "2002::4"
					comment = "test comment on AAAA record"
					ext_attrs = jsonencode({
						"Location" = "test AAAA location"
					})
					depends_on = [infoblox_zone_auth.zone]
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"infoblox_aaaa_record.foo3", "ext_attrs",
						`{"Location":"test AAAA location"}`,
					),
					func(s *terraform.State) error {
						conn := testAccProvider.Meta().(ibclient.IBConnector)

						res, found := s.RootModule().Resources["infoblox_aaaa_record.foo3"]
						if !found {
							return fmt.Errorf("not found: %s", "infoblox_aaaa_record.foo3")
						}

						id := res.Primary.ID
						if id == "" {
							return fmt.Errorf("ID is not set")
						}

						objMgr := ibclient.NewObjectManager(
							conn,
							"terraform_test",
							"terraform_test_tenant")
						qarec, err := objMgr.GetARecordByRef(id)
						if err != nil {
							if isNotFoundError(err) {
								return fmt.Errorf("object with ID '%s' not found, but expected to exist", id)
							}
						}

						if _, ok := qarec.Ea["Site"]; ok {
							return fmt.Errorf("Site EA should've been removed, but still present in the WAPI object")
						}
						return nil
					},
				),
			},
		},
	})
}
