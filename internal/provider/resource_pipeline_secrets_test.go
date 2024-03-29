package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/healx/terraform-provider-nftower/internal/template"
)

func TestAccResourcePipelineSecrets(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_pipeline_secrets",
				Config:       template.ParseRandName(testAccResourcePipelineSecrets),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"nftower_pipeline_secrets.foo", "name", "tf_acceptance_pipeline_secrets"),
					resource.TestCheckResourceAttr(
						"nftower_pipeline_secrets.foo", "value", "something secret"),
					resource.TestMatchResourceAttr(
						"nftower_pipeline_secrets.foo", "date_created", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestMatchResourceAttr(
						"nftower_pipeline_secrets.foo", "last_updated", regexp.MustCompile("^[0-9-:TZ]+")),
				),
			},
		},
	})
}

const testAccResourcePipelineSecrets = `
resource "nftower_workspace" "foo" {
  name        = "tf-acceptance-{{.randName}}"
  full_name   = "tf acceptance testing pipeline secrets"
	
  description = "Created by the nftower terraform provider acceptance tests. Will be deleted shortly"
  visibility  = "PRIVATE"
}

resource "nftower_pipeline_secrets" "foo" {
  name        = "tf_acceptance_pipeline_secrets"
  workspace_id = nftower_workspace.foo.id
  value = "something secret"
}
`
