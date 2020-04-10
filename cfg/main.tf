resource "demo_resource_foo" "example" {
  name = "magodo"
  addr {
    country = "China"
    city = "aaa"
  }
  addr {
    country = "US"
    city = "xxx"
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
