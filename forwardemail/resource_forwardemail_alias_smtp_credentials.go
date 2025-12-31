package forwardemail

import (
	"context"
	"fmt"

	"github.com/abagayev/go-forwardemail/forwardemail"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAliasSmtpCredentials() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Fully qualified domain name (FQDN).",
			},
			"alias": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Alias name.",
			},
			"new_password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Custom password for the alias. Leave blank for a randomly generated strong password.",
			},
			"is_override": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Override existing password. WARNING: This permanently deletes existing IMAP storage and resets the alias' SQLite email database.",
			},
			"emailed_instructions": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Email address to send password and setup instructions to.",
			},
			"username": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SMTP username (email address).",
			},
			"password": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The generated SMTP password. Only available if emailed_instructions is not set.",
			},
		},
		CreateContext: resourceAliasSmtpCredentialsCreate,
		ReadContext:   resourceAliasSmtpCredentialsRead,
		UpdateContext: resourceAliasSmtpCredentialsUpdate,
		DeleteContext: resourceAliasSmtpCredentialsDelete,
	}
}

func resourceAliasSmtpCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*forwardemail.Client)
	domain := d.Get("domain").(string)
	alias := d.Get("alias").(string)

	params := forwardemail.GeneratePasswordParameters{
		IsOverride: toBool(d.Get("is_override")),
	}

	if v, ok := d.GetOk("new_password"); ok {
		s := v.(string)
		params.NewPassword = &s
	}

	if v, ok := d.GetOk("emailed_instructions"); ok {
		s := v.(string)
		params.EmailedInstructions = &s
	}

	credentials, err := client.GenerateAliasPassword(domain, alias, params)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("username", credentials.Username); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("password", credentials.Password); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", domain, alias))

	return nil
}

func resourceAliasSmtpCredentialsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*forwardemail.Client)
	domain := d.Get("domain").(string)
	alias := d.Get("alias").(string)

	_, err := client.GetAlias(domain, alias)
	if err != nil {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceAliasSmtpCredentialsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChanges("new_password", "is_override", "emailed_instructions") {
		client := meta.(*forwardemail.Client)
		domain := d.Get("domain").(string)
		alias := d.Get("alias").(string)

		params := forwardemail.GeneratePasswordParameters{
			IsOverride: toBool(d.Get("is_override")),
		}

		if v, ok := d.GetOk("new_password"); ok {
			s := v.(string)
			params.NewPassword = &s
		}

		if v, ok := d.GetOk("emailed_instructions"); ok {
			s := v.(string)
			params.EmailedInstructions = &s
		}

		credentials, err := client.GenerateAliasPassword(domain, alias, params)
		if err != nil {
			return diag.FromErr(err)
		}

		if err := d.Set("username", credentials.Username); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("password", credentials.Password); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

func resourceAliasSmtpCredentialsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
