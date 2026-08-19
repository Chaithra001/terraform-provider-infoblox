// Authoritative DNS zones for LBDN association
resource "infoblox_zone_auth" "parent_zone" {
  nios = {
    fqdn = "wapi.com"
    view = "default"
    grid_primary = [
      {
        name = "infoblox.localdomain"
      }
    ]
  }
}

resource "infoblox_zone_auth" "parent_zone2" {
  nios = {
    fqdn = "info.com"
    view = "default"
    grid_primary = [
      {
        name = "infoblox.localdomain"
      }
    ]
  }
}

// Create a DTC LBDN with basic fields
resource "infoblox_dtc_lbdn" "lbdn_basic" {
  nios = {
    name      = "example_lbdn_1"
    lb_method = "SOURCE_IP_HASH"
  }
}

// Create a DTC LBDN with additional fields
resource "infoblox_dtc_lbdn" "lbdn_additional_fields" {
  nios = {
    name = "example_lbdn_2"
    auth_zones = [
      infoblox_zone_auth.parent_zone.id,
      infoblox_zone_auth.parent_zone2.id
    ]
    comment   = "lbdn with additional parameters"
    ext_attrs = { Site = "location-1" }
    lb_method = "RATIO"
    patterns  = ["*wapi.com", "info.com*"]
    pools = [
      {
        pool  = "dtc:pool/ZG5zLmlkbnNfcG9vbCRUZXN0UG9vbDE:TestPool1"
        ratio = 2
      }
    ]
    ttl         = 60
    disable     = false
    types       = ["A", "AAAA"]
    persistence = 100
    priority    = 1
  }
}
