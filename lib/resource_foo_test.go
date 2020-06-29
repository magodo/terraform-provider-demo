package lib

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccExampleService_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: map[string]terraform.ResourceProvider{
			"demo": Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccExampleService_basic(),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				ResourceName:      "demo_resource_foo.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccExampleService_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: map[string]terraform.ResourceProvider{
			"demo": Provider(),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccExampleService_before_update(),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				ResourceName:      "demo_resource_foo.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccExampleService_after_update(),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				ResourceName:      "demo_resource_foo.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccExampleService_basic() string {
	return fmt.Sprintf(`
resource "demo_resource_foo" "foo" {
  name = "abc"
  addr {
    country = "China"
  }
  contact {
    github = "Magodo"
  }
}
`)
}

func testAccExampleService_before_update() string {
	return fmt.Sprintf(`
resource "demo_resource_foo" "foo" {
  name = "abc"
  addr {
    country = "China"
  }
  addr {
    country = "UK"
  }
}
`)
}

func testAccExampleService_after_update() string {
	return fmt.Sprintf(`
resource "demo_resource_foo" "foo" {
  name = "abc"
  addr {
    country = "China"
  }
  addr {
    country = "US"
  }
}
`)
}
