package lib

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	return nil
}

func resourceBarDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*client).ClientBar
	return client.Delete(d.Id())
}
