package identity_platform

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	idp "github.com/sigmadigitalza/identity-platform-client"
)

type Object map[string]interface{}
type TerraformListType []Object

func resourceConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConfigCreate,
		ReadContext:   resourceConfigRead,
		UpdateContext: resourceConfigUpdate,
		DeleteContext: resourceRecordDelete,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"password_required": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			"notification": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sendEmail": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"callbackUri": {
										Type:     schema.TypeString,
										Optional: true,
										Default:  true,
									},
								},
							},
						},
					},
				},
			},
			"phone_number": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
			"subtype": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "IDENTITY_PLATFORM",
			},
			"authorized_domains": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceConfigCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*idp.Service)
	projectId, newConfig := configFromResourceData(d)

	config, err := client.Project.UpdateConfig(ctx, projectId, newConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(config.Name)

	resourceConfigRead(ctx, d, m)

	return diag.Diagnostics{}
}

func resourceConfigRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*idp.Service)
	projectId := d.Get("project_id").(string)

	config, err := client.Project.GetConfig(ctx, projectId)
	if err != nil {
		return diag.FromErr(err)
	}

	return hydrate(diag.Diagnostics{}, config, d)
}

func resourceConfigUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*idp.Service)
	projectId, newConfig := configFromResourceData(d)

	_, err := client.Project.UpdateConfig(ctx, projectId, newConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceConfigRead(ctx, d, m)
}

func resourceRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*idp.Service)
	projectId := d.Get("project_id").(string)
	subtype := d.Get("subtype").(string)

	emptyConfig := &idp.Config{
		SignIn:            nil,
		Subtype:           subtype,
		AuthorizedDomains: nil,
	}

	_, err := client.Project.UpdateConfig(ctx, projectId, emptyConfig)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return diag.Diagnostics{}
}

func configFromResourceData(d *schema.ResourceData) (string, *idp.Config) {
	projectId := d.Get("project_id").(string)
	email := d.Get("email").([]interface{})
	phoneNumber := d.Get("phone_number").([]interface{})
	subtype := d.Get("subtype").(string)
	authorizedDomains := d.Get("authorized_domains").([]interface{})
	notification := d.Get("notification").([]interface{})

	sendEmail := extractProperties(notification)["sendEmail"].([]interface{})

	config := &idp.Config{
		SignIn: &idp.SignInConfig{
			Email: &idp.Email{
				Enabled:          extractProperties(email)["enabled"].(bool),
				PasswordRequired: extractProperties(email)["password_required"].(bool),
			},
			PhoneNumber: &idp.PhoneNumber{
				Enabled: extractProperties(phoneNumber)["enabled"].(bool),
			},
		},
		Notification: &idp.NotificationConfig{
			SendEmail: &idp.SendEmail{
				CallbackUri: extractProperties(sendEmail)["callbackUri"].(string),
			},
		},
		Subtype:           subtype,
		AuthorizedDomains: extractStringSlice(authorizedDomains),
	}

	return projectId, config
}

func extractStringSlice(getResult []interface{}) []string {
	slice := make([]string, len(getResult))

	for i := 0; i < len(getResult); i++ {
		slice[i] = getResult[i].(string)
	}

	return slice
}

func extractProperties(getResult []interface{}) Object {
	return getResult[0].(map[string]interface{})
}

func hydrate(diags diag.Diagnostics, config *idp.Config, d *schema.ResourceData) diag.Diagnostics {
	if err := d.Set("name", config.Name); err != nil {
		return diag.FromErr(err)
	}

	arr := TerraformListType{
		{
			"enabled":           config.SignIn.Email.Enabled,
			"password_required": config.SignIn.Email.PasswordRequired,
		},
	}

	if err := d.Set("email", arr); err != nil {
		diag.FromErr(err)
	}

	arr = TerraformListType{
		{
			"sendEmail": Object{
				"callbackUri": config.Notification.SendEmail.CallbackUri,
			},
		},
	}

	if err := d.Set("notification", arr); err != nil {
		diag.FromErr(err)
	}

	arr = TerraformListType{
		{
			"enabled": config.SignIn.PhoneNumber.Enabled,
		},
	}
	if err := d.Set("phone_number", arr); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("subtype", config.Subtype); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("authorized_domains", config.AuthorizedDomains); err != nil {
		return diag.FromErr(err)
	}

	return diags
}
