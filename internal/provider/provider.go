package provider

import (
	"github.com/automato-io/binocs-client-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// New returns a Terraform provider resource
func New() func() *schema.Provider {
	return func() *schema.Provider {
		return &schema.Provider{
			Schema: map[string]*schema.Schema{
				"access_key": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("BINOCS_ACCESS_KEY", ""),
					Description: "Access Key required to communicate with Binocs API. Get yours at [https://binocs.sh](https://binocs.sh)",
				},
				"secret_key": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("BINOCS_SECRET_KEY", ""),
					Description: "Secret Key required to communicate with Binocs API. Get yours at [https://binocs.sh](https://binocs.sh)",
				},
			},
			ConfigureFunc: configureProvider,
			ResourcesMap: map[string]*schema.Resource{
				"binocs_check":   checkResource(),
				"binocs_channel": channelResource(),
			},
		}
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	config := binocs.ClientConfig{
		AccessKey: d.Get("access_key").(string),
		SecretKey: d.Get("secret_key").(string),
	}
	binocs, err := binocs.New(config)
	if err != nil {
		return nil, err
	}
	return binocs, nil
}
