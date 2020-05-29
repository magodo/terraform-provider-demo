resource "demo_resource_foo" "example" {
  name = "magodo"
  job  = "doctor2"
  #   addr {
  #     country = "Germany"
  #     city = "Wetzla"
  #   }
}

resource "demo_resource_bar" "example" {
  name   = "magodo"
  github = demo_resource_foo.example.output_job
  phone  = 2
  #   addr {
  #     country = "Germany"
  #     city = "Wetzla"
  #   }
}


# resource "demo_resource_bar" "example" {
#   name = "kinoko"
#   github = demo_resource_foo.example.contact.0.github
# }
