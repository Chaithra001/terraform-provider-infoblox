// Create a DTC LBDN (Policy) with basic required fields
resource "infoblox_dtc_lbdn" "lbdn_basic" {
  uddi = {
    name   = "example-policy-1"
    method = "round_robin"
  }
}

// Create a DTC LBDN (Policy) with pools and additional fields
resource "infoblox_dtc_lbdn" "lbdn_with_pools" {
  uddi = {
    name     = "example-policy-2"
    method   = "ratio"
    comment  = "Created by Terraform"
    disabled = false
    ttl      = 60
    tags     = { Site = "location-1" }
    pools = [
      {
        pool_id = "dtc/pool/example-pool-id"
        weight  = 2
      }
    ]
  }
}

// Create a DTC LBDN (Policy) with topology method
resource "infoblox_dtc_lbdn" "lbdn_topology" {
  uddi = {
    name   = "example-policy-topology"
    method = "topology"
    topology = {
      id = "dtc/topology/example-topology-id"
    }
  }
}
