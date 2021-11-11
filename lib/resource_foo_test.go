package lib

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var providerFactories = map[string]func() (*schema.Provider, error){
	"demo": func() (*schema.Provider, error) {
		return Provider(), nil
	},
}

func TestAccExampleService_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
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

func TestAccExampleService_updateJob(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExampleService_withJob(),
				Check:  resource.ComposeTestCheckFunc(),
			},
			{
				ResourceName:      "demo_resource_foo.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
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
		ProviderFactories: providerFactories,
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
}
`)
}

func testAccExampleService_withJob() string {
	return fmt.Sprintf(`
resource "demo_resource_foo" "foo" {
  name = "abc"
  job = "foo"
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
    country = "UK"
  }
}
`)
}
