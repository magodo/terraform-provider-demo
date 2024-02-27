package lib

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceFoo() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceFooRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"age": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"job": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parameters": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"contact": {
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
						"other_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"other_map": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"addr": {
				Type:     schema.TypeSet,
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
			"tags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceFooRead(d *schema.ResourceData, m interface{}) error {
	id := d.Get("name").(string)
	d.SetId(id)

	client := m.(*client).ClientFoo
	resp, err := client.Get(id)
	if err != nil {
		return err
	}

	d.Set("age", resp.Age)
	d.Set("job", resp.Job)
	d.Set("metadata", resp.Metadata)
	d.Set("parameters", resp.Parameters)

	contact, err := flattenContact(resp.Contact)
	if err != nil {
		return err
	}
	if err := d.Set("contact", contact); err != nil {
		return err
	}
	d.Set("addr", flattenAddrs(resp.Addrs))
	d.Set("tags", FlattenStringMap(resp.Tags))

	return nil
}
