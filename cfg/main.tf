terraform {
  required_providers {
    demo = {
    source = "magodo/demo"
  }
  }
}

resource "demo_resource_foo" "magodo" {
  name = "magodo"
  job = "aaa"
  contact {
    phone = 123
    github = "xxx"
  }
  addr {
    country = "China"
  }
  addr {
    country = "uS"
  }
}
