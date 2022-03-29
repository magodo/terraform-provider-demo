package acctest

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/magodo/terraform-provider-demo/client"
	"github.com/magodo/terraform-provider-demo/demo"
)

func Providers() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"demo": func() (tfprotov6.ProviderServer, error) {
			return tfsdk.NewProtocol6Server(demo.New()), nil
		},
	}
}

func ProviderConfig() string {
	envFsWorkdir := os.Getenv(EnvFsWorkdir)
	envJsUrl := os.Getenv(EnvJsUrl)
	if envFsWorkdir != "" {
		return fmt.Sprintf(`
provider "demo" {
  filesystem = {
    workdir = "%s"
  }
}
`, envFsWorkdir)
	}
	return fmt.Sprintf(`
provider "demo" {
  jsonserver = {
    url = "%s"
  }
}
`, envJsUrl)
}

func PreCheck(t *testing.T, customChecker func()) {
	envFsWorkdir := os.Getenv(EnvFsWorkdir)
	envJsUrl := os.Getenv(EnvJsUrl)
	if (envFsWorkdir == "") == (envJsUrl == "") {
		t.Fatalf("One of either environment variable %q or %q has to be set", EnvFsWorkdir, EnvJsUrl)
	}
	if customChecker != nil {
		customChecker()
	}
}

func IsDestroy(rt string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		c, err := buildClient()
		if err != nil {
			return err
		}
		for _, resource := range s.RootModule().Resources {
			if resource.Type != rt {
				continue
			}

			if label, err := c.Read(resource.Primary.ID); err != client.ErrNotFound {
				return fmt.Errorf("reading %s.%s: %v", resource.Type, label, err)
			}
		}
		return nil
	}
}
