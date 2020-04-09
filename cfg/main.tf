resource "demo_resource_foo" "example" {
  name = "magodo"
  addr {
    country = "China"
  }
  addr {
    country = "US"
  }
#   addr {
#     country = "Germany"
#     city = "Wetzla"
#   }
}

# resource "demo_resource_bar" "example" {
#   name = "kinoko"
#   github = demo_resource_foo.example.contact.0.github
# }
