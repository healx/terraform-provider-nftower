package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/healx/terraform-provider-nftower/internal/template"
)

func TestAccDataSourceCredentialsAWS(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: template.ParseRandName(testAccDataSourceCredentialsAWS),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nftower_credentials.foo", "name", "tf-acceptance-credentials-ds-aws"),
					resource.TestCheckResourceAttr(
						"data.nftower_credentials.foo", "description", "tf acceptance testing aws ds credentials"),
					resource.TestMatchResourceAttr(
						"data.nftower_credentials.foo", "date_created", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestMatchResourceAttr(
						"data.nftower_credentials.foo", "last_updated", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestCheckResourceAttr(
						"data.nftower_credentials.foo", "aws.0.access_key", "foo"),
					resource.TestCheckNoResourceAttr("data.nftower_credentials.foo", "aws.0.secret_key"),
				),
			},
		},
	})
}

const testAccDataSourceCredentialsAWS = `
resource "nftower_workspace" "foo" {
  name        = "tf-acceptance-{{.randName}}"
  full_name   = "tf acceptance testing ds credentials"
	
  description = "Created by the nftower terraform provider acceptance tests. Will be deleted shortly"
  visibility  = "PRIVATE"
}

resource "nftower_credentials" "foo" {
  name        = "tf-acceptance-credentials-ds-aws"
  description = "tf acceptance testing aws ds credentials"
  workspace_id = nftower_workspace.foo.id

  aws {
	access_key      = "foo"
	secret_key      = "bar"
	assume_role_arn = "baz"
  }
}

data "nftower_credentials" "foo" {
	name = nftower_credentials.foo.name
	workspace_id = nftower_credentials.foo.workspace_id
}
`

func TestAccDataSourceCredentialsGithub(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: template.ParseRandName(testAccDataSourceCredentialsGithub),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nftower_credentials.foo", "name", "tf-acceptance-credentials-ds-github"),
					resource.TestCheckResourceAttr(
						"data.nftower_credentials.foo", "description", "tf acceptance testing github ds credentials"),
					resource.TestMatchResourceAttr(
						"data.nftower_credentials.foo", "date_created", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestMatchResourceAttr(
						"data.nftower_credentials.foo", "last_updated", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestCheckResourceAttr(
						"data.nftower_credentials.foo", "github.0.username", "foo"),
					resource.TestCheckNoResourceAttr(
						"data.nftower_credentials.foo", "github.0.access_token"),
				),
			},
		},
	})
}

const testAccDataSourceCredentialsGithub = `
resource "nftower_workspace" "foo" {
  name        = "tf-acceptance-{{.randName}}"
  full_name   = "tf acceptance testing ds credentials"
	
  description = "Created by the nftower terraform provider acceptance tests. Will be deleted shortly"
  visibility  = "PRIVATE"
}

resource "nftower_credentials" "foo" {
  name        = "tf-acceptance-credentials-ds-github"
  description = "tf acceptance testing github ds credentials"
  workspace_id = nftower_workspace.foo.id

  github {
	username     = "foo"
	access_token = "bar"
  }
}

data "nftower_credentials" "foo" {
	name = nftower_credentials.foo.name
	workspace_id = nftower_credentials.foo.workspace_id
}
`
