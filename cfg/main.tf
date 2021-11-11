resource "demo_resource_foo" "magodo" {
  name = "magodo"
  contact {
    other_string = jsonencode({
      allowedLocations = "abc"
    })
  }

  # lifecycle {
  #   ignore_changes = [age]
  # }
}
