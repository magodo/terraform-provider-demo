resource "demo_resource_foo" "magodo" {
  name = "magodo"
  job  = "teacher"
}

resource "demo_resource_foo" "rr" {
  name   = "rr"
  job = demo_resource_foo.magodo.output_job
  age  = 2
}
