case "basic" {
  backend  = "uddi"
  parallel = true

  step {
    uddi {
      name   = "dtc-lbdn-{{random}}"
      method = "round_robin"
    }
    check = {
      "uddi.name"   = "dtc-lbdn-{{random}}"
      "uddi.method" = "round_robin"
    }
  }

}

case "disappears" {
  backend               = "uddi"
  disappears            = true
  expect_non_empty_plan = true
  parallel              = true

  step {
    uddi {
      name   = "dtc-lbdn-{{random}}"
      method = "round_robin"
    }
  }

}

case "comment" {
  backend  = "uddi"
  parallel = true

  step {
    uddi {
      name    = "dtc-lbdn-{{random}}"
      method  = "round_robin"
      comment = "Created by Terraform"
    }
    check = {
      "uddi.comment" = "Created by Terraform"
    }
  }

  step {
    uddi {
      name    = "dtc-lbdn-{{random}}"
      method  = "round_robin"
      comment = "Updated by Terraform"
    }
    check = {
      "uddi.comment" = "Updated by Terraform"
    }
  }

}

case "disabled" {
  backend  = "uddi"
  parallel = true

  step {
    uddi {
      name     = "dtc-lbdn-{{random}}"
      method   = "round_robin"
      disabled = true
    }
    check = {
      "uddi.disabled" = "true"
    }
  }

  step {
    uddi {
      name     = "dtc-lbdn-{{random}}"
      method   = "round_robin"
      disabled = false
    }
    check = {
      "uddi.disabled" = "false"
    }
  }

}

case "method" {
  backend  = "uddi"
  parallel = true

  step {
    uddi {
      name   = "dtc-lbdn-{{random}}"
      method = "global_availability"
    }
    check = {
      "uddi.method" = "global_availability"
    }
  }

  step {
    uddi {
      name   = "dtc-lbdn-{{random}}"
      method = "ratio"
    }
    check = {
      "uddi.method" = "ratio"
    }
  }

}

case "name" {
  backend  = "uddi"
  parallel = true

  step {
    uddi {
      name   = "dtc-lbdn-{{random}}"
      method = "round_robin"
    }
    check = {
      "uddi.name" = "dtc-lbdn-{{random}}"
    }
  }

  step {
    uddi {
      name   = "dtc-lbdn-{{random2}}"
      method = "round_robin"
    }
    check = {
      "uddi.name" = "dtc-lbdn-{{random2}}"
    }
  }

}

case "tags" {
  backend  = "uddi"
  parallel = true

  step {
    uddi {
      name   = "dtc-lbdn-{{random}}"
      method = "round_robin"
      tags   = { Site = "{{random2}}" }
    }
    check = {
      "uddi.tags.Site" = "{{random2}}"
    }
  }

  step {
    uddi {
      name   = "dtc-lbdn-{{random}}"
      method = "round_robin"
      tags   = { Site = "{{random3}}" }
    }
    check = {
      "uddi.tags.Site" = "{{random3}}"
    }
  }

}

case "ttl" {
  backend  = "uddi"
  parallel = true

  step {
    uddi {
      name   = "dtc-lbdn-{{random}}"
      method = "round_robin"
      ttl    = 60
    }
    check = {
      "uddi.ttl" = "60"
    }
  }

  step {
    uddi {
      name   = "dtc-lbdn-{{random}}"
      method = "round_robin"
      ttl    = 120
    }
    check = {
      "uddi.ttl" = "120"
    }
  }

}
