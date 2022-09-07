package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/healx/terraform-provider-nftower/internal/template"
)

func TestAccResourceCredentialsAWS(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_credentials",
				Config:       template.ParseRandName(testAccResourceCredentialsAWS),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "name", "tf-acceptance-credentials-aws"),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "description", "tf acceptance testing aws credentials"),
					resource.TestMatchResourceAttr(
						"nftower_credentials.foo", "date_created", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestMatchResourceAttr(
						"nftower_credentials.foo", "last_updated", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "aws.0.access_key", "foo"),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "aws.0.secret_key", "bar"),
				),
			},
		},
	})
}

const testAccResourceCredentialsAWS = `
resource "nftower_workspace" "foo" {
  name        = "tf-acceptance-{{.randName}}"
  full_name   = "tf acceptance testing credentials"
	
  description = "Created by the nftower terraform provider acceptance tests. Will be deleted shortly"
  visibility  = "PRIVATE"
}

resource "nftower_credentials" "foo" {
  name        = "tf-acceptance-credentials-aws"
  description = "tf acceptance testing aws credentials"
  workspace_id = nftower_workspace.foo.id

  aws {
	access_key      = "foo"
	secret_key      = "bar"
	assume_role_arn = "baz"
  }
}
`

func TestAccResourceCredentialsGithub(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_credentials",
				Config:       template.ParseRandName(testAccResourceCredentialsGithub),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "name", "tf-acceptance-credentials-github"),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "description", "tf acceptance testing github credentials"),
					resource.TestMatchResourceAttr(
						"nftower_credentials.foo", "date_created", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestMatchResourceAttr(
						"nftower_credentials.foo", "last_updated", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "github.0.username", "foo"),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "github.0.access_token", "bar"),
				),
			},
		},
	})
}

const testAccResourceCredentialsGithub = `
resource "nftower_workspace" "foo" {
  name        = "tf-acceptance-{{.randName}}"
  full_name   = "tf acceptance testing credentials"
	
  description = "Created by the nftower terraform provider acceptance tests. Will be deleted shortly"
  visibility  = "PRIVATE"
}

resource "nftower_credentials" "foo" {
  name        = "tf-acceptance-credentials-github"
  description = "tf acceptance testing github credentials"
  workspace_id = nftower_workspace.foo.id

  github {
	username     = "foo"
	access_token = "bar"
  }
}
`
