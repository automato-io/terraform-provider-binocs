package provider

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/automato-io/binocs-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	protocolHTTP  = "HTTP"
	protocolHTTPS = "HTTPS"
	protocolTCP   = "TCP"
)

const (
	validDNSNamePattern = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`
	validUpCodePattern  = `^([1-5]{1}[0-9]{2}-[1-5]{1}[0-9]{2}|([1-5]{1}(([0-9]{2}|[0-9]{1}x)|xx))){1}(,([1-5]{1}[0-9]{2}-[1-5]{1}[0-9]{2}|([1-5]{1}(([0-9]{2}|[0-9]{1}x)|xx))))*$`
)

var supportedRegions = []string{
	"af-south-1",
	"ap-east-1",
	"ap-northeast-1",
	"ap-south-1",
	"ap-southeast-1",
	"ap-southeast-2",
	"eu-central-1",
	"eu-west-1",
	"sa-east-1",
	"us-east-1",
	"us-west-1",
}

var supportedHTTPMethods = []string{"GET", "HEAD", "POST", "PUT", "DELETE"}

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
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The name (alias) of this check. Maximum length is 25 characters.",
				Default:      "",
				ValidateFunc: validation.StringLenBetween(0, 25),
			},
			"resource": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The resource to check; a valid URL in case of a HTTP(S) resource, or tcp://{HOSTNAME}:{PORT} in case of a TCP resource.",
				ForceNew:    true,
				ValidateFunc: func(i interface{}, s string) (_ []string, errors []error) {
					v, ok := i.(string)
					if !ok {
						errors = append(errors, fmt.Errorf("expected type of %q to be string", s))
						return
					}
					if strings.HasPrefix(v, "tcp://") {
						rc := strings.Split(v[6:], ":")
						if len(rc) != 2 {
							errors = append(errors, fmt.Errorf("expected %q tcp resource to contain host and port components", s))
							return
						}
						if !isHost(rc[0]) {
							errors = append(errors, fmt.Errorf("expected %q tcp resource to contain a valid host", s))
						}
						port, err := strconv.Atoi(rc[1])
						if err != nil {
							errors = append(errors, fmt.Errorf("expected %q tcp resource to contain a port number", s))
							return
						}
						if _, errs := validation.IsPortNumber(port, s); len(errs) > 0 {
							errors = append(errors, fmt.Errorf("expected %q tcp resource to contain a valid port number", s))
						}
						return
					}
					if strings.HasPrefix(v, "http") {
						_, err := validation.IsURLWithHTTPorHTTPS(v, s)
						if len(err) > 0 {
							errors = append(errors, err...)
						}
						return
					}
					errors = append(errors, fmt.Errorf("expected %q to be either TCP or HTTP(S) resource", s))
					return
				},
			},
			"method": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  fmt.Sprintf("The HTTP method (one of %s). Only used and required with HTTP(S) resources.", strings.Join(supportedHTTPMethods, ", ")),
				ValidateFunc: validation.StringInSlice(supportedHTTPMethods, false),
			},
			"interval": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "How often Binocs checks this resource, in seconds. Minimum is 5 and maximum is 900 seconds.",
				Default:      60,
				ValidateFunc: validation.IntBetween(5, 900),
			},
			"target": {
				Type:         schema.TypeFloat,
				Optional:     true,
				Description:  "The response time that accommodates Apdex=1.0 (in seconds with up to 3 decimal places). Valid target is a value between 0.01 and 10.0 seconds.",
				Default:      1.2,
				ValidateFunc: validation.FloatBetween(0.01, 10.0),
			},
			"regions": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: fmt.Sprintf("From where in the world Binocs checks this resource. At least one region is required. Valid values: %s.", strings.Join(supportedRegions, ", ")),
				MinItems:    1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: func(i interface{}, s string) (_ []string, errors []error) {
						v, ok := i.(string)
						if !ok {
							errors = append(errors, fmt.Errorf("expected type of %q to be string", s))
							return
						}
						if !stringInSlice(v, supportedRegions) {
							errors = append(errors, fmt.Errorf("expected %q to be any of %q", "regions", strings.Join(supportedRegions, ", ")))
						}
						return
					},
				},
			},
			"up_codes": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The good (\"up\") HTTP(S) response codes, e.g. 2xx or `200-302`, or `200,301`. Only used with HTTP(S) resources.",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(validUpCodePattern), fmt.Sprintf("up_codes must be of a format such as %q", "200, 2xx, 200-302, 200,301")),
			},
			"up_confirmations_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "How many subsequent \"up\" responses need to occur before Binocs creates an incident and triggers notifications. Minimum is 1 and maximum is 10.",
				Default:      2,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"down_confirmations_threshold": {
				Type:         schema.TypeInt,
				Optional:     true,
				Description:  "How many subsequent \"down\" responses need to occur before Binocs closes an incident and triggers \"recovery\" notifications. Minimum is 1 and maximum is 10.",
				Default:      2,
				ValidateFunc: validation.IntBetween(1, 10),
			},
		},
	}
}

func checkCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*binocs.Client)
	payload, err := constructCheckPayload(d)
	if err != nil {
		return err
	}
	check, err := client.Checks.Create(payload)
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
	if err != nil && strings.Contains(err.Error(), "404") {
		d.SetId("")
		return false, nil
	}
	return err == nil, err
}

func checkUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*binocs.Client)
	payload, err := constructCheckPayload(d)
	if err != nil {
		return err
	}
	err = client.Checks.Update(d.Id(), payload)
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

func constructCheckPayload(d *schema.ResourceData) (binocs.Check, error) {
	payload := binocs.Check{}

	if v, ok := d.GetOk("name"); ok {
		payload.Name = v.(string)
	}

	if v, ok := d.GetOk("resource"); ok {
		payload.Resource = v.(string)
		payload.Protocol = strings.ToUpper(strings.Split(v.(string), ":")[0])
	}

	if payload.Protocol == protocolHTTP || payload.Protocol == protocolHTTPS {
		if v, ok := d.GetOk("method"); ok {
			payload.Method = v.(string)
		} else {
			return payload, fmt.Errorf("expected \"method\" to be one of %s for a %s resource", strings.Join(supportedHTTPMethods, ", "), payload.Protocol)
		}
	} else {
		if _, ok := d.GetOk("method"); ok {
			return payload, fmt.Errorf("\"method\" cannot be used with a %s resource", payload.Protocol)
		}
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

	if payload.Protocol == protocolHTTP || payload.Protocol == protocolHTTPS {
		if v, ok := d.GetOk("up_codes"); ok {
			payload.UpCodes = v.(string)
		} else {
			return payload, fmt.Errorf("expected \"up_codes\" to be set for a %s resource", payload.Protocol)
		}
	} else {
		if _, ok := d.GetOk("up_codes"); ok {
			return payload, fmt.Errorf("\"up_codes\" cannot be used with a %s resource", payload.Protocol)
		}
	}

	if v, ok := d.GetOk("up_confirmations_threshold"); ok {
		payload.UpConfirmationsThreshold = v.(int)
	}

	if v, ok := d.GetOk("down_confirmations_threshold"); ok {
		payload.DownConfirmationsThreshold = v.(int)
	}

	return payload, nil
}

func isIP(str string) bool {
	return net.ParseIP(str) != nil
}

func isDNSName(str string) bool {
	if str == "" || len(strings.Replace(str, ".", "", -1)) > 255 {
		return false
	}
	return !isIP(str) && regexp.MustCompile(validDNSNamePattern).MatchString(str)
}

func isHost(str string) bool {
	return isIP(str) || isDNSName(str)
}

func stringInSlice(s string, slice []string) bool {
	for _, h := range slice {
		if s == h {
			return true
		}
	}
	return false
}
