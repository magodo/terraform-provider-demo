package acctest

import (
	"fmt"
	"os"

	"github.com/magodo/terraform-provider-demo/client"
)

const (
	EnvFsWorkdir = "DEMO_FS_WORKDIR"
	EnvJsUrl     = "DEMO_JS_URL"
)

func buildClient() (client.Client, error) {
	envFsWorkdir := os.Getenv(EnvFsWorkdir)
	envJsUrl := os.Getenv(EnvJsUrl)
	if (envFsWorkdir == "") == (envJsUrl == "") {
		return nil, fmt.Errorf("One of either environment variable %q or %q has to be set", EnvFsWorkdir, EnvJsUrl)
	}
	if envFsWorkdir != "" {
		return client.NewFsClient(envFsWorkdir)
	}
	return client.NewJSONServerClient(envJsUrl)
}
