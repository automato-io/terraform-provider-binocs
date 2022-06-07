package provider

import (
	"fmt"

	"github.com/automato-io/binocs-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func checkResource() *schema.Resource {
	return &schema.Resource{
		Description: "`binocs_check` defines a check",

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Create: checkCreate,
		Read:   checkRead,
		Exists: checkExists,
		Update: checkUpdate,
		Delete: checkDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name (alias) of this check. Maximum length is 25 characters.",
				Default:     "",
			},
			"resource": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The resource to check; a URL in case of a HTTP(S) resource, or HOSTNAME:PORT in case of a TCP resource.",
				ForceNew:    true,
			},
			"method": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The HTTP method (one of GET, HEAD, POST, PUT or DELETE). Only required for HTTP(S) resources.",
				// RequiredWith: []string{"protocol=http|https"},
			},
			"interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "How often Binocs checks this resource, in seconds. Minimum is 5 and maximum is 900 seconds.",
				Default:     60,
			},
			"target": {
				Type:        schema.TypeFloat,
				Optional:    true,
				Description: "The response time that accommodates Apdex=1.0 (in seconds with up to 3 decimal places). Valid target is a value between 0.01 and 10.0 seconds.",
				Default:     1.2,
			},
			"regions": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "From where in the world Binocs checks this resource.",
				Elem: &schema.Schema{
					Type:     schema.TypeString,
					Default:  []string{"eu-west-1", "eu-central-1", "us-east-1", "us-west-1", "ap-southeast-1", "ap-northeast-1"},
					MinItems: 1,
				},
			},
			"up_codes": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The good (\"up\") HTTP(S) response codes, e.g. 2xx or `200-302`, or `200,301`",
				Default:     "200-302",
			},
			"up_confirmations_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "How many subsequent \"up\" responses need to occur before Binocs creates an incident and triggers notifications. Minimum is 1 and maximum is 10.",
				Default:     2,
			},
			"down_confirmations_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "How many subsequent \"down\" responses need to occur before Binocs closes an incident and triggers \"recovery\" notifications. Minimum is 1 and maximum is 10.",
				Default:     2,
			},
		},
	}
}

func checkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*binocs.Client)
	check, err := client.Checks.Create(constructCheckPayload(d))
	if err != nil {
		return fmt.Errorf("unable to create Binocs check: %s", err)
	}
	d.SetId(check.Ident)
	return checkRead(d, meta)
}

func checkRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*binocs.Client)
	check, err := client.Checks.Read(d.Id())
	if err != nil {
		return fmt.Errorf("unable to read Binocs check: %s", err)
	}
	for k, v := range map[string]interface{}{
		"name":                         check.Name,
		"resource":                     check.Resource,
		"method":                       check.Method,
		"interval":                     check.Interval,
		"target":                       check.Target,
		"regions":                      check.Regions,
		"up_codes":                     check.UpCodes,
		"up_confirmations_threshold":   check.UpConfirmationsThreshold,
		"down_confirmations_threshold": check.DownConfirmationsThreshold,
	} {
		if err := d.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

func checkExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	err := checkRead(d, meta)
	return err == nil, err
}

func checkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*binocs.Client)
	err := client.Checks.Update(d.Id(), constructCheckPayload(d))
	if err != nil {
		return fmt.Errorf("unable to update Binocs check: %s", err)
	}
	return nil
}

func checkDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*binocs.Client)
	err := client.Checks.Delete(d.Id())
	if err != nil {
		return fmt.Errorf("unable to remove Binocs check: %s", err)
	}
	return nil
}

func constructCheckPayload(d *schema.ResourceData) binocs.Check {
	payload := binocs.Check{}

	if v, ok := d.GetOk("name"); ok {
		payload.Name = v.(string)
	}

	if v, ok := d.GetOk("protocol"); ok {
		payload.Protocol = v.(string)
	}

	if v, ok := d.GetOk("resource"); ok {
		payload.Resource = v.(string)
	}

	if v, ok := d.GetOk("method"); ok {
		payload.Method = v.(string)
	}

	if v, ok := d.GetOk("interval"); ok {
		payload.Interval = v.(int)
	}

	if v, ok := d.GetOk("target"); ok {
		payload.Target = v.(float64)
	}

	if v, ok := d.GetOk("regions"); ok {
		interfaceSlice := v.(*schema.Set).List()
		var stringSlice []string
		for s := range interfaceSlice {
			stringSlice = append(stringSlice, interfaceSlice[s].(string))
		}
		payload.Regions = stringSlice
	}

	if v, ok := d.GetOk("up_codes"); ok {
		payload.UpCodes = v.(string)
	}

	if v, ok := d.GetOk("up_confirmations_threshold"); ok {
		payload.UpConfirmationsThreshold = v.(int)
	}

	if v, ok := d.GetOk("down_confirmations_threshold"); ok {
		payload.DownConfirmationsThreshold = v.(int)
	}

	return payload
}
