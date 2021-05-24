package main

import (
	"context"
	"flag"
	"log"

	"github.com/magodo/terraform-provider-demo/lib"

	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	var debugMode bool
	flag.BoolVar(&debugMode, "debuggable", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	if debugMode {
		err := plugin.Debug(context.Background(), "registry.terraform.io/hashicorp/demo",
			&plugin.ServeOpts{
				ProviderFunc: func() terraform.ResourceProvider {
					return lib.Provider()
				},
			})
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: func() terraform.ResourceProvider {
				return lib.Provider()
			},
		})
	}
}
