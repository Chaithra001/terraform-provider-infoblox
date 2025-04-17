# Range Resource

The `infoblox_ipv4_range` resource enables you to perform `create`, `update` and `delete` operations on Network Range in a NIOS appliance.
The resource represents the ‘range’ WAPI object in NIOS.

The following list describes the parameters you can define in the resource block of the Network Range object:

* `name`: optional, specifies the display name. Example: `network-range`.
* `comment`: optional, comment for the range, maximum 256 characters. Example: `test range`.
* `network`: optional, The network to which this range belongs, in IPv4 Address/CIDR format. Example: `21.20.2.0/24`.
* `network_view`: optional, The name of the network view in which this range resides. Example: `default`.
* `start_addr`: required, The IPv4 Address starting address of the range. Example: `21.20.2.20`.
* `end_addr`: required, The IPv4 Address end address of the range. Example: `21.20.2.40`
* `disable`: optional, Determines whether a range is disabled or not. When this is set to False, the range is enabled. Default value: `false`. 
* `extattrs`: optional, Extensible attributes associated with the object. Example: `"{\"*Site\":\"Antarctica\"}"`
* `failover_association`: optional, The name of the failover association: the server in this failover association will serve the IPv4 range in case the main server is out of service. `server_association_type` must be set to `FAILOVER` or `FAILOVER_MS` if you want the failover association specified here to serve the range.
* `server_association_type`: optional, The type of server that is going to serve the range. Valid values are `FAILOVER`,`MEMBER`,`MS_FAILOVER`,`MS_SERVER`,`NONE`. Default value: `NONE`.
* `options`: optional, An array of DHCP option structs that lists the DHCP options associated with the object.
```terraform
options {
    name         = "dhcp-lease-time"
    value        = "43200"
    vendor_class = "DHCP"
    num          = 51
    use_option   = true
  }
```
* `use_options`: optional, Use option is a flag that indicates whether the options field are used or not. The default value is false. Example: `false`
* `member`: optional, specifies the member that will provide service for this range. `server_association_type` needs to be set to `MEMBER` if you want the server specified here to serve the range. `member` has the following three fields `name`, `ipv4addr` and `ipv6addr`. At least one of `name`, `ipv4addr`, or `ipv6addr` is required in the `member` block.
  The description of the fields of `member` is as follows:
    * `name`: optional, specifies the name of the pool. Example: `infoblox.localdomain`.
    * `ipv4addr`: optional, specifies the weight of the pool. Example: `11.10.1.0`.
    * `ipv6addr`: optional, specifies the IPv6 address of the member. Example: `2403:8600:80cf:e10c:3a00::1192`.

Example for `member`:
```terraform
member = {
  name = "infoblox.localdomain"
  ipv4addr = "11.10.1.0"
  ipv6addr = "2403:8600:80cf:e10c:3a00::1192"
}
```

### Examples of a Network Range Block
```hcl
// creating a Network Range
resource "infoblox_ipv4_range" "range3" {
  start_addr = "17.0.0.221"
  end_addr   = "17.0.0.240"
  options {
    name         = "dhcp-lease-time"
    value        = "43200"
    vendor_class = "DHCP"
    num          = 51
    use_option   = true
  }
  network              = "17.0.0.0/24"
  network_view = "default"
  comment              = "test comment"
  name                 = "test_range"
  disable              = false
  member = {
    name = "infoblox.localdomain"
  }
  server_association_type= "MEMBER"
  ext_attrs = jsonencode({
    "Site" = "Blr"
  })
  use_options = true
}
```
