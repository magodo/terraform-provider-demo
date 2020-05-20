package main

import (
	"log"

	"github.com/magodo/terraform-provider-demo/lib"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return lib.Provider()
		},
	})
}
