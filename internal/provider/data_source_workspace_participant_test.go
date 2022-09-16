package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/healx/terraform-provider-nftower/internal/template"
)

func TestAccDataSourceWorkspaceParticipant(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: template.ParseRandName(testAccDataSourceWorkspaceParticipant),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"nftower_workspace_participant.foo", "email", regexp.MustCompile("^tf-acceptance-[0-9]+@example.com")),
					resource.TestCheckResourceAttr(
						"nftower_workspace_participant.foo", "role", "launch"),
					resource.TestMatchResourceAttr(
						"nftower_workspace_participant.foo", "member_id", regexp.MustCompile("^[0-9]+$")),
				),
			},
		},
	})
}

const testAccDataSourceWorkspaceParticipant = `
resource "nftower_organization_member" "foo" {
  email = "tf-acceptance-{{.randName}}@example.com"
}

resource "nftower_workspace" "foo" {
  name        = "tf-acceptance-{{.randName}}"
  full_name   = "tf acceptance testing workspace"
  description = "Created by the nftower terraform provider acceptance tests. Will be deleted shortly"
  visibility  = "PRIVATE"
}

resource "nftower_workspace_participant" "foo" {
  workspace_id = nftower_workspace.foo.id
  member_id = nftower_organization_member.foo.id
}

data "nftower_workspace_participant" "foo" {
  workspace_id = nftower_workspace.foo.id
  email = nftower_workspace_participant.foo.email
}
`
