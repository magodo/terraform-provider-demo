package lib

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"demo_resource_foo": resourceFoo(),
			"demo_resource_bar": resourceBar(),
		},
		ConfigureFunc: func(d *schema.ResourceData) (interface{}, error) {
			return &client{
				&ClientFoo{},
				&ClientBar{},
			}, nil
		},
	}
}

type client struct {
	*ClientFoo
	*ClientBar
}
