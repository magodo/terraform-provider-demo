package lib

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/acctest"
	"os"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var once sync.Once

func init() {
	// require reattach testing is enabled
	os.Setenv("TF_ACCTEST_REATTACH", "1")

	once.Do(func() {
		acctest.UseBinaryDriver("demo", func() terraform.ResourceProvider {
			return Provider()
		})
	})
}

func TestAccExampleService_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]terraform.ResourceProviderFactory{
			"demo"	: func() (terraform.ResourceProvider, error) {
				return Provider(), nil
			},
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

func TestAccExampleService_updateJob(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]terraform.ResourceProviderFactory{
			"demo"	: func() (terraform.ResourceProvider, error) {
				return Provider(), nil
			},
		},
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
		ProviderFactories: map[string]terraform.ResourceProviderFactory{
			"demo"	: func() (terraform.ResourceProvider, error) {
				return Provider(), nil
			},
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
