package demo_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/magodo/terraform-provider-demo/demo/acctest"
)

type Foo struct{}

func TestAccFoo_basic(t *testing.T) {
	foo := Foo{}
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: acctest.Providers(),
		PreCheck:                 func() { acctest.PreCheck(t, nil) },
		CheckDestroy:             acctest.IsDestroy("demo_foo"),
		Steps: []resource.TestStep{
			{
				Config: foo.basic(),
			},
			{
				ResourceName:      "demo_foo.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func (_ Foo) basic() string {
	return fmt.Sprintf(`%s

resource "demo_foo" "test" {
  string  = "str"
  int64   = 1
  float64 = 1.0
  number  = 1e2
  bool    = true
}
`, acctest.ProviderConfig())
}
