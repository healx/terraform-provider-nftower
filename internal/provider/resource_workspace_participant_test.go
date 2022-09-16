package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/healx/terraform-provider-nftower/internal/template"
)

func TestAccResourceWorkspaceParticipant(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_workspace_participant",
				Config:       template.ParseRandName(testAccResourceWorkspaceParticipant),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"nftower_workspace_participant.foo", "email", regexp.MustCompile("^tf-acceptance-[0-9]+@example.com")),
					resource.TestCheckResourceAttr(
						"nftower_workspace_participant.foo", "role", "launch"),
				),
			},
		},
	})
}

const testAccResourceWorkspaceParticipant = `
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
`

func TestAccResourceWorkspaceParticipantMaintain(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_workspace_participant",
				Config:       template.ParseRandName(testAccResourceWorkspaceParticipantMaintain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"nftower_workspace_participant.foo", "email", regexp.MustCompile("^tf-acceptance-[0-9]+@example.com")),
					resource.TestCheckResourceAttr(
						"nftower_workspace_participant.foo", "role", "maintain"),
				),
			},
		},
	})
}

const testAccResourceWorkspaceParticipantMaintain = `
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
  member_id    = nftower_organization_member.foo.id
  role         = "maintain"
}
`

func TestAccResourceWorkspaceParticipantUpdateRole(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_workspace_participant",
				Config:       template.ParseRandName(testAccResourceWorkspaceParticipantUpdateRole_1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"nftower_workspace_participant.foo", "email", regexp.MustCompile("^tf-acceptance-[0-9]+@example.com")),
					resource.TestCheckResourceAttr(
						"nftower_workspace_participant.foo", "role", "launch"),
				),
			},
			{
				ResourceName: "nftower_workspace_participant",
				Config:       template.ParseRandName(testAccResourceWorkspaceParticipantUpdateRole_2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"nftower_workspace_participant.foo", "role", "maintain"),
				),
			},
		},
	})
}

const testAccResourceWorkspaceParticipantUpdateRole_1 = `
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
  member_id    = nftower_organization_member.foo.id
  role         = "launch"
}
`

const testAccResourceWorkspaceParticipantUpdateRole_2 = `
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
  member_id    = nftower_organization_member.foo.id
  role         = "maintain"
}
`
