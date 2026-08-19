// Fetch a specific DTC LBDN (Policy) by name
data "infoblox_dtc_lbdn" "by_name" {
  filters = { name = "example-policy-1" }
}

// Fetch DTC LBDNs by tag
data "infoblox_dtc_lbdn" "by_tags" {
  tag_filters = { Site = "location-1" }
}

// Fetch all DTC LBDNs
data "infoblox_dtc_lbdn" "all" {}
