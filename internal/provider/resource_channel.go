package provider

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/automato-io/binocs-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var supportedChannelTypes = []string{
	"email",
}

var unsupportedChannelTypes = []string{
	"telegram",
	"slack",
}

var validHandlePattern = map[string]string{
	"email": `^(?:[a-z0-9!#$%&'*+/=?^_{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\])$`,
}

func channelResource() *schema.Resource {
	return &schema.Resource{
		Description: "`binocs_channel` defines a notification channel",

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Create: channelCreate,
		Read:   channelRead,
		Exists: channelExists,
		Update: channelUpdate,
		Delete: channelDelete,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  fmt.Sprintf("The only supported channel is currently \"email\", and it requires e-mail address verification. All other notification channels (%s) currently require interactive creation using Binocs CLI. All notification channels can be imported to Terraform.", strings.Join(unsupportedChannelTypes, ", ")),
				ValidateFunc: validation.StringInSlice(supportedChannelTypes, false),
			},
			"handle": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The e-mail address for a channel of `type = email`.",
				ValidateFunc: validation.StringMatch(regexp.MustCompile(validHandlePattern["email"]), "expected a valid e-mail address"),
			},
			"alias": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "",
				Description:  "The alias (name) of this notification channel. Maximum length is 25 characters.",
				ValidateFunc: validation.StringLenBetween(0, 25),
			},
			"checks": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Description: "The checks to associate with this notifications channel.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Set: schema.HashString,
			},
		},
	}
}

func channelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*binocs.Client)
	payload, err := constructChannelPayload(d)
	if err != nil {
		return err
	}
	channel, err := client.Channels.Create(payload)
	if err != nil {
		return fmt.Errorf("unable to create Binocs channel: %s", err)
	}
	d.SetId(channel.Ident)

	if v, ok := d.GetOk("checks"); ok {
		checksSlice := v.(*schema.Set).List()
		for s := range checksSlice {
			err = client.Channels.Attach(channel.Ident, checksSlice[s].(string))
			if err != nil {
				return fmt.Errorf("unable to attach Binocs channel %q to check %q: %s", channel.Ident, checksSlice[s].(string), err)
			}
		}
	}

	return channelRead(d, meta)
}

func channelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*binocs.Client)
	channel, err := client.Channels.Read(d.Id())
	if err != nil {
		return fmt.Errorf("unable to read Binocs channel: %s", err)
	}
	for k, v := range map[string]interface{}{
		"type":   channel.Type,
		"handle": channel.Handle,
		"alias":  channel.Alias,
		"checks": channel.Checks,
	} {
		if err := d.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}

func channelExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	err := channelRead(d, meta)
	if err != nil && strings.Contains(err.Error(), "404") {
		d.SetId("")
		return false, nil
	}
	return err == nil, err
}

func channelUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*binocs.Client)
	payload, err := constructChannelPayload(d)
	if err != nil {
		return err
	}
	err = client.Channels.Update(d.Id(), payload)
	if err != nil {
		return fmt.Errorf("unable to update Binocs channel: %s", err)
	}
	if d.HasChange("checks") {
		o, n := d.GetChange("checks")
		if o == nil {
			o = new(schema.Set)
		}
		if n == nil {
			n = new(schema.Set)
		}
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		detach := os.Difference(ns).List()
		attach := ns.Difference(os).List()
		if len(detach) > 0 {
			for _, r := range detach {
				err = client.Channels.Detach(d.Id(), r.(string))
				if err != nil {
					return fmt.Errorf("unable to detach Binocs channel %q from check %q: %s", d.Id(), r.(string), err)
				}
			}
		}
		if len(attach) > 0 {
			for _, a := range attach {
				err = client.Channels.Attach(d.Id(), a.(string))
				if err != nil {
					return fmt.Errorf("unable to attach Binocs channel %q to check %q: %s", d.Id(), a.(string), err)
				}
			}
		}
	}
	return nil
}

func channelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*binocs.Client)
	err := client.Channels.Delete(d.Id())
	if err != nil {
		return fmt.Errorf("unable to remove Binocs channel: %s", err)
	}
	return nil
}

func constructChannelPayload(d *schema.ResourceData) (binocs.Channel, error) {
	payload := binocs.Channel{}

	if v, ok := d.GetOk("type"); ok {
		payload.Type = v.(string)
	}

	if v, ok := d.GetOk("handle"); ok {
		payload.Handle = v.(string)
	}

	if v, ok := d.GetOk("alias"); ok {
		payload.Alias = v.(string)
	}

	return payload, nil
}
