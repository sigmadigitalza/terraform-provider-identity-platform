package identity_platform

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	idp "github.com/sigmadigitalza/identity-platform-client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https://www.googleapis.com/auth/cloud-platform",
			},
		},
		ConfigureContextFunc: configureContext,
		ResourcesMap:         map[string]*schema.Resource{
			"identity_platform_config": resourceConfig(),
		},
	}
}

func configureContext(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	scope := d.Get("scope").(string)

	service, err := idp.New(ctx, scope)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return service, nil
}
