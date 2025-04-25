# Range Template-record Data Source

Use the `infoblox_ipv4_range_template` data resource for the IPV4 Range Template record to retrieve the following information for a Range Template record:

* `name`: The name of the range template record. Example: `range-template1`.
* `number_of_addresses`: The number of addresses for this range. Example: `100`.
* `offset`: The start address offset for the range. Example: `30`.
* `use_options`: Use flag for options. Example: `true`.
* `cloud_api_compatible`: The flag controls whether this template can be used to create network objects in a cloud-computing deployment. Example: `true`.
* `options`: An array of DHCP option structs that lists the DHCP options associated with the object. The description of the fields of `options` is as follows:
  * `name`: The Name of the DHCP option. Example: `domain-name-servers`.
  * `num`: The code of the DHCP option. Example: `6`.
  * `value`: The value of the option. Example: `11.22.33.44`.
  * `vendor_class`: The name of the space this DHCP option is associated to. Default value is `DHCP`.
  * `use_option`:Only applies to special options that are displayed separately from other options and have a use flag. These options are `router`,
    `router-templates`, `domain-name-servers`, `domain-name`, `broadcast-address`, `broadcast-address-offset`, `dhcp-lease-time`, and `dhcp6.name-servers`.
```terraform
options { 
    name = "domain-name-servers"
    value = "11.22.33.44"
    num = 6
    use_option = true
  }
```
* `comment`: The description of the record. This is a regular comment. Example: `Temporary Ipv4 Shared Network`.
* `ext_attrs`: The set of extensible attributes of the record, if any. The content is formatted as string of JSON map. Example: `"{\"Site\":"Nagoya"}"`
* `server_association_type`: The type of server that is going to serve the range. Valid values are: `FAILOVER`, `MEMBER`, `MS_FAILOVER`, `MS_SERVER`, `NONE` .Example: `NONE`.
* `failover_association`: The name of the failover association: the server in this failover association will serve the IPv4 range in case the main server is out of service. Example: `test.com`.
* `ms_server`: The Microsoft server that will provide service for this range. `server_association_type` needs to be set to `MS_SERVER` if you want the server specified here to serve the range.
* `member`: The member that will provide service for this range. `server_association_type` needs to be set to `MEMBER` if you want the server specified here to serve the range. `member` has the following three fields `name`, `ipv4addr` and `ipv6addr`.The description of the fields of `member` is as follows:
  * `name`: The name of the Grid member. Example: `infoblox.localdomain`.
  * `ipv4addr`: The IPv4 Address of the Grid Member. Example: `11.10.1.0`.
  * `ipv6addr`: The IPv6 address of the member. Example: `2403:8600:80cf:e10c:3a00::1192`.

Example for `member`:
```terraform
member = { 
    name = "infoblox.localdomain"
    ipv4addr = "11.22.33.44"
    ipv6addr = "2403:8600:80cf:e10c:3a00::1192"
  }
```

For usage of filters, add the fields as keys and appropriate values to be passed to the keys like `name`, `view`, `zone`, `comment`, `target_name`, and `target_type`  corresponding to object.
From the below list of supported arguments for filters,  use only the searchable fields for retrieving the matching records.

### Supported Arguments for filters

-----
| Field                   | Alias                   | Type   | Searchable |
|-------------------------|-------------------------|--------|------------|
| name                    | name                    | string | yes        |
| failover_association    | failover_association    | string | yes        |
| comment                 | comment                 | string | yes        |
| server_association_type | server_association_type | string | yes        |

!> Any of the combination from searchable fields in supported arguments list for fields are allowed.

!> Please consider using only fields as the keys in terraform datasource filters, kindly don't use alias names as keys from the above table.

!>The search functionality using the filters argument is not supported for member and ms_server fields.

### Example for using the filters:
 ```hcl
data "infoblox_ipv4_range_template" "range_template_read" {
  filters = {
    name = "range-template1"
    comment = "Temporary range template"
  }
}
 ```

!> If `null` or empty filters are passed, then all the records or objects associated with datasource like here `infoblox_ipv4_range_template`, will be fetched in results.

### Example of the Alias-record Data Source Block

This example defines a data source of type `infoblox_ipv4_range_template` and the name "range_template_read", which is configured in a Terraform file.
You can reference this resource and retrieve information about it.

```hcl
resource "infoblox_ipv4_range_template" "range_template_record" {
  name = "range-template2"
  number_of_addresses = 40
  offset = 30
  comment = "Temporary Range Template"
  cloud_api_compatible = false
  use_options = true
  ext_attrs = jsonencode({
    "Site" = "Kobe"
  })
  options {
    name = "domain-name-servers"
    value = "11.22.33.44"
    vendor_class = "DHCP"
    num = 6
    use_option = true
  }
  member {
    ipv4addr = "10.197.81.146"
    ipv6addr = "2403:8600:80cf:e10c:3a00::1192"
    name = "infoblox.localdomain"
  }
  failover_association = "failover1"
  server_association_type = "FAILOVER"
}

data "infoblox_ipv4_range_template" "range_template_read"{
  filters = {
    name = infoblox_ipv4_range_template.range_template_record.name
    server_association_type = infoblox_ipv4_range_template.range_template_record.server_association_type
    failover_association = infoblox_ipv4_range_template.range_template_record.failover_association
  }
}

output "range_template_record_out" {
  value = data.infoblox_ipv4_range_template.range_template_read
}

// accessing individual field in results
output "range_template_offset_out" {
  value = data.infoblox_ipv4_range_template.range_template_read.results.0.offset //zero represents index of json object from results list
}

// accessing Alias-Record through EA's
data "infoblox_ipv4_range_template" "range_template_read_ea" {
  filters = {
    "*Site" = "Kobe"
  }
}

// throws matching Alias records with EA, if any
output "range_template_read_ea_out" {
  value = data.infoblox_ipv4_range_template.range_template_read_ea
}
```
