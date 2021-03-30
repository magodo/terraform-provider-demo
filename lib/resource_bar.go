package lib

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBar() *schema.Resource {
	return &schema.Resource{
		Create: resourceBarCreateOrUpdate,
		Read:   resourceBarRead,
		Update: resourceBarCreateOrUpdate,
		Delete: resourceBarDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"job": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"phone": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"locations_deprecated": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"locations": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceBarCreateOrUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*client).ClientBar
	name := d.Get("name").(string)

	param := &ModelBar{
		Name: &name,
	}

	if job, ok := d.GetOk("job"); ok {
		param.Job = StringPtr(job.(string))
	}
	if phone, ok := d.GetOk("phone"); ok {
		param.Phone = IntPtr(phone.(int))
	}

	param.Locations = expandLocationSlicePtr(d.Get("locations").(*schema.Set).List())
	if param.Locations == nil {
		param.Locations = expandLocationSlicePtr(d.Get("locations_deprecated").(*schema.Set).List())
	}

	resp, err := client.CreateOrUpdate(param)
	if err != nil {
		return err
	}

	d.SetId(*resp.Id)
	return resourceBarRead(d, m)
}

func resourceBarRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*client).ClientBar
	resp, err := client.Get(d.Id())
	if err != nil {
		return err
	}

	d.Set("name", resp.Name)
	d.Set("job", resp.Job)
	d.Set("phone", resp.Phone)

	if err := d.Set("locations", flattenLocationSlicePtr(resp.Locations)); err != nil {
		return fmt.Errorf(`setting "locations": %v`, err)
	}

	if err := d.Set("locations_deprecated", flattenLocationSlicePtr(resp.Locations)); err != nil {
		return fmt.Errorf(`setting "locations_deprecated": %v`, err)
	}

	return nil
}

func resourceBarDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client).ClientBar
	return client.Delete(d.Id())
}

func expandLocationSlicePtr(input []interface{}) *[]Location {
	output := make([]Location, 0)
	for _, elem := range input {
		output = append(output, Location{Name: StringPtr(elem.(string))})
	}
	return &output
}

func flattenLocationSlicePtr(input *[]Location) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	output := make([]interface{}, 0)
	for _, elem := range *input {
		name := ""
		if elem.Name != nil {
			name = *elem.Name
		}
		output = append(output, name)
	}
	return output
}
