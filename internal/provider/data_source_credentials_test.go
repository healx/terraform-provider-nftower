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

// TestAccDataSourceCredentialsContainerRegistry cannot be run withouth real credentials for accessing the container registry.
// Seqera Platform test registry connection during the creation of the credentials.
// For AWS ECR would be username: AWS_ACCESS_KEY_ID), password: AWS_SECRET_ACCESS_KEY, registry_server: AWS_ACCOUNT_ID.dkr.ecr.AWS_REGION.amazonaws.com
func TestAccDataSourceCredentialsContainerRegistry(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: template.ParseRandName(testAccDataSourceContainerRegistry),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nftower_credentials.foo", "name", "tf-acceptance-credentials-ds-container-registry"),
					resource.TestCheckResourceAttr(
						"data.nftower_credentials.foo", "description", "tf acceptance testing container registry ds credentials"),
					resource.TestMatchResourceAttr(
						"data.nftower_credentials.foo", "date_created", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestMatchResourceAttr(
						"data.nftower_credentials.foo", "last_updated", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestCheckResourceAttr(
						"data.nftower_credentials.foo", "container_registry.0.username", "<<AWS_ACCESS_KEY_ID>>"),
					resource.TestCheckResourceAttr(
						"data.nftower_credentials.foo", "container_registry.0.registry_server", "<<AWS_ACCOUNT_ID.dkr.ecr.AWS_REGION.amazonaws.com>>"),
					resource.TestCheckNoResourceAttr(
						"data.nftower_credentials.foo", "container_registry.0.password"),
				),
			},
		},
	})
}

const testAccDataSourceContainerRegistry = `
resource "nftower_workspace" "foo" {
  name        = "tf-acceptance-{{.randName}}"
  full_name   = "tf acceptance container testing ds credentials"

  description = "Created by the nftower terraform provider acceptance tests. Will be deleted shortly"
  visibility  = "PRIVATE"
}

resource "nftower_credentials" "foo" {
  name        = "tf-acceptance-credentials-ds-container-registry"
  description = "tf acceptance testing container registry ds credentials"
  workspace_id = nftower_workspace.foo.id

  container_registry {
	username      = "<<AWS_ACCESS_KEY_ID>>"
	password      = "<<AWS_SECRET_ACCESS_KEY>>"
	registry_server = "<<AWS_ACCOUNT_ID.dkr.ecr.AWS_REGION.amazonaws.com>>"
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

func TestAccDataSourceCredentialsGitlab(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: template.ParseRandName(testAccDataSourceCredentialsGitlab),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.nftower_credentials.foo", "name", "tf-acceptance-credentials-ds-gitlab"),
					resource.TestCheckResourceAttr(
						"data.nftower_credentials.foo", "description", "tf acceptance testing gitlab ds credentials"),
					resource.TestMatchResourceAttr(
						"data.nftower_credentials.foo", "date_created", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestMatchResourceAttr(
						"data.nftower_credentials.foo", "last_updated", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestCheckResourceAttr(
						"data.nftower_credentials.foo", "gitlab.0.username", "foo"),
					resource.TestCheckNoResourceAttr(
						"data.nftower_credentials.foo", "gitlab.0.password"),
					resource.TestCheckNoResourceAttr(
						"data.nftower_credentials.foo", "gitlab.0.token"),
				),
			},
		},
	})
}

const testAccDataSourceCredentialsGitlab = `
resource "nftower_workspace" "foo" {
  name        = "tf-acceptance-{{.randName}}"
  full_name   = "tf acceptance testing ds credentials"
	
  description = "Created by the nftower terraform provider acceptance tests. Will be deleted shortly"
  visibility  = "PRIVATE"
}

resource "nftower_credentials" "foo" {
  name        = "tf-acceptance-credentials-ds-gitlab"
  description = "tf acceptance testing gitlab ds credentials"
  workspace_id = nftower_workspace.foo.id

  gitlab {
	username     = "foo"
	password     = "bar"
	token        = "baz"
  }
}

data "nftower_credentials" "foo" {
	name = nftower_credentials.foo.name
	workspace_id = nftower_credentials.foo.workspace_id
}
`
