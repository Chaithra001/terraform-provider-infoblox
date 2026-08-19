// Fetch a specific DTC LBDN by name
data "infoblox_dtc_lbdn" "by_name" {
  filters = { name = "example_lbdn_1" }
}

// Fetch DTC LBDNs by extensible attribute
data "infoblox_dtc_lbdn" "by_ext_attrs" {
  ext_attr_filters = { Site = "location-1" }
}

// Fetch all DTC LBDNs
data "infoblox_dtc_lbdn" "all" {}
