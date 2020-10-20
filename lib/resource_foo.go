package lib

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceFoo() *schema.Resource {
	return &schema.Resource{
		Create: resourceFooCreateOrUpdate,
		Read:   resourceFooRead,
		Update: resourceFooCreateOrUpdate,
		Delete: resourceFooDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"age": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"job": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"contact": {
				Type:     schema.TypeList,
				Optional: true,
				//Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"phone": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"github": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"addr": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"country": {
							Type:     schema.TypeString,
							Required: true,
						},
						"city": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"roads": {
							Type:     schema.TypeMap,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"tag": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output_job": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_addr": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"country": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"city": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"roads": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"output_contact": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"phone": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"github": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceFooCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client).ClientFoo
	name := d.Get("name").(string)

	param := &ModelFoo{
		Name: &name,
	}

	if age, ok := d.GetOkExists("age"); ok {
		param.Age = IntPtr(age.(int))
	}

	if job, ok := d.GetOkExists("job"); ok {
		param.Job = StringPtr(job.(string))
	}

	param.Contact = expandContact(d.Get("contact").([]interface{}))
	param.Addrs = expandAddrs(d.Get("addr").(*schema.Set).List())
	param.Tags = ExpandStringMap(d.Get("tag").(map[string]interface{}))

	resp, err := client.CreateOrUpdate(param)
	if err != nil {
		return err
	}

	d.SetId(*resp.Id)
	return resourceFooRead(d, m)
}

func resourceFooRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client).ClientFoo
	resp, err := client.Get(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", resp.Name)
	// for root level property, below operations are different:
	// - not set at all: will results into a null in tf state
	// - set a nil: will results into a `""` in tf state
	if resp.Age != nil {
		d.Set("age", resp.Age)
	}
	if resp.Job != nil {
		d.Set("job", resp.Job)
	}

	d.Set("contact", flattenContact(resp.Contact))
	d.Set("addr", flattenAddrs(resp.Addrs))
	d.Set("tags", FlattenStringMap(resp.Tags))
	d.Set("output_contact", flattenContact(resp.Contact))
	d.Set("output_addr", flattenAddrs(resp.Addrs))
	job := ""
	if resp.Job != nil {
		job = *resp.Job
	}
	d.Set("output_job", job)

	return nil
}

func resourceFooDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client).ClientFoo
	return client.Delete(d.Id())
}

func expandContact(contact []interface{}) *ContactFoo {
	if len(contact) == 0 {
		return nil
	}

	v := contact[0].(map[string]interface{})

	// Looks like the nil check below is redundent.
	// If the map inside interface is not nil, then after get, the value of nested fields are always filled.

	// if v["github"] != nil {
	// 	contact.Github = StringPtr(v["github"].(string))
	// }
	// if v["phone"] != nil {
	// 	contact.Phone = IntPtr(v["phone"].(int))
	// }
	output := &ContactFoo{
		Github: StringPtr(v["github"].(string)),
		Phone:  IntPtr(v["phone"].(int)),
	}
	return output
}

func expandAddrs(addrs []interface{}) *[]*Addr {
	result := make([]*Addr, 0)
	for _, v := range addrs {
		if v != nil {
			m := v.(map[string]interface{})
			addr := &Addr{
				Country: StringPtr(m["country"].(string)),
			}
			if city, ok := m["city"].(string); ok {
				addr.City = &city
			}
			result = append(result, addr)
		}
	}

	return &result
}

func flattenContact(input *ContactFoo) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	// Given following cfg:
	//
	//----------------------------------------
	//resource "demo_resource_foo" "example" {
	//  name = "magodo"
	//  contact {
	//    phone = 123
	//  }
	//}
	//
	//resource "demo_resource_bar" "example" {
	//  name = "kinoko"
	//  github = demo_resource_foo.example.contact.0.github
	//}
	//----------------------------------------
	//
	// demo_resource_bar.example refers to a non specified field of
	// demo_resource_foo.example.
	// No matter what we specify for "flattenStrategy", it will always
	// ends up with following error when apply:
	// (ensure a clean apply with no prior state)
	//
	// 	When expanding the plan for demo_resource_bar.example to include new values
	// 	learned so far during apply, provider "registry.terraform.io/-/demo" produced
	// 	an invalid new value for .github: was null, but now cty.StringVal("").
	//
	// 	This is a bug in the provider, which should be reported in the provider's own
	// 	issue tracker.

	flattenStrategy := 3

	// Possible ways of handling:
	//
	// case 1 might not set a property at all, which results into null value
	// case 2 && case 3 always set properties (though might set nil), which results into at least default value of the correct type.
	switch flattenStrategy {
	// 1. only set it if non nil in response
	case 1:
		result := make(map[string]interface{})
		if input.Phone != nil {
			result["phone"] = *input.Phone
		}
		if input.Github != nil {
			result["github"] = *input.Github
		}
		return []interface{}{result}
	// 2. set pointer anyway
	//	  If setting a nil into map, it will ends up with the default value of that type. But this behavior is not guarateed
	//    to be supported in long term, so we'd better to set the default value by ourselves (i.e. use the 3rd way).
	case 2:
		result := make(map[string]interface{})
		result["phone"] = input.Phone
		result["github"] = input.Github
		return []interface{}{result}
		// 3. set it if non nil in response, otherwise, set a default value
	case 3:
		phone, github := 0, ""
		if input.Phone != nil {
			phone = *input.Phone
		}
		if input.Github != nil {
			github = *input.Github
		}
		return []interface{}{map[string]interface{}{
			"phone":  phone,
			"github": github,
		}}
	}
	panic("not supposed to reach here")
}

func flattenAddrs(input *[]*Addr) []interface{} {
	if input == nil {
		return nil
	}

	output := make([]interface{}, 0)

	for _, v := range *input {
		if v == nil {
			continue
		}

		var (
			country string
			city    string
		)

		if v.Country != nil {
			country = *v.Country
		}
		if v.City != nil {
			city = *v.City
		}
		output = append(output, map[string]interface{}{
			"country": country,
			"city":    city,
		})
	}

	return output
}
