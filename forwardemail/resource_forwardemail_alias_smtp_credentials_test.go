package forwardemail

import (
	"fmt"
	"testing"

	"github.com/abagayev/go-forwardemail/forwardemail"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccForwardemailAliasSmtpCredentials_basic(t *testing.T) {
	domain := fake.Internet().Domain()
	name := fake.Internet().User()
	recipient := fake.Internet().FreeEmail()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccForwardemailProviderFactories,
		CheckDestroy:      testAccCheckForwardemailAliasSmtpCredentialsDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckForwardemailAliasSmtpCredentialsConfig_basic, domain, name, recipient),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("forwardemail_alias_smtp_credentials.test", "domain", domain),
					resource.TestCheckResourceAttr("forwardemail_alias_smtp_credentials.test", "alias", name),
					resource.TestCheckResourceAttrSet("forwardemail_alias_smtp_credentials.test", "username"),
					resource.TestCheckResourceAttrSet("forwardemail_alias_smtp_credentials.test", "password"),
				),
			},
		},
	})
}

func TestAccForwardemailAliasSmtpCredentials_customPassword(t *testing.T) {
	domain := fake.Internet().Domain()
	name := fake.Internet().User()
	recipient := fake.Internet().FreeEmail()
	customPassword := fake.Internet().Password()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccForwardemailProviderFactories,
		CheckDestroy:      testAccCheckForwardemailAliasSmtpCredentialsDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckForwardemailAliasSmtpCredentialsConfig_customPassword, domain, name, recipient, customPassword),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("forwardemail_alias_smtp_credentials.test", "domain", domain),
					resource.TestCheckResourceAttr("forwardemail_alias_smtp_credentials.test", "alias", name),
					resource.TestCheckResourceAttrSet("forwardemail_alias_smtp_credentials.test", "username"),
					resource.TestCheckResourceAttr("forwardemail_alias_smtp_credentials.test", "password", customPassword),
				),
			},
		},
	})
}

func testAccCheckForwardemailAliasSmtpCredentialsDestroy(s *terraform.State) error {
	client := testAccForwardemailProvider.Meta().(*forwardemail.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "forwardemail_alias_smtp_credentials" {
			continue
		}

		if _, err := client.GetAlias(rs.Primary.Attributes["domain"], rs.Primary.Attributes["alias"]); err == nil {
			return fmt.Errorf("alias still exists")
		}
	}

	return nil
}

const testAccCheckForwardemailAliasSmtpCredentialsConfig_basic = `
	resource "forwardemail_domain" "test" {
		name = "%s"
	}

	resource "forwardemail_alias" "test" {
		name   = "%s"
		domain = forwardemail_domain.test.name
		recipients = ["%s"]
	}

	resource "forwardemail_alias_smtp_credentials" "test" {
		domain      = forwardemail_domain.test.name
		alias       = forwardemail_alias.test.name
		is_override = true
	}
`

const testAccCheckForwardemailAliasSmtpCredentialsConfig_customPassword = `
	resource "forwardemail_domain" "test" {
		name = "%s"
	}

	resource "forwardemail_alias" "test" {
		name   = "%s"
		domain = forwardemail_domain.test.name
		recipients = ["%s"]
	}

	resource "forwardemail_alias_smtp_credentials" "test" {
		domain       = forwardemail_domain.test.name
		alias        = forwardemail_alias.test.name
		new_password = "%s"
		is_override  = true
	}
`
