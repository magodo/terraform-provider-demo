terraform {
  required_providers {
    demo = {
    source = "magodo/demo"
  }
  }
}

resource "demo_resource_foo" "magodo" {
  name = "magodo"
  addr {
    country = "China"
  }
}
