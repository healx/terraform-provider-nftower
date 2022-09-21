package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/healx/terraform-provider-nftower/internal/template"
)

func TestAccResourceWorkspaceParticipant_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_workspace_participant",
				Config:       template.ParseRandName(testAccResourceWorkspaceParticipant_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"nftower_workspace_participant.foo", "email", regexp.MustCompile("^tf-acceptance-[0-9]+@example.com")),
					resource.TestCheckResourceAttr(
						"nftower_workspace_participant.foo", "role", "view"),
				),
			},
		},
	})
}

const testAccResourceWorkspaceParticipant_basic = `
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

func TestAccResourceWorkspaceParticipant_email(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_workspace_participant",
				Config:       template.ParseRandName(testAccResourceWorkspaceParticipant_email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"nftower_workspace_participant.foo", "email", regexp.MustCompile("^tf-acceptance-[0-9]+@example.com")),
						resource.TestMatchResourceAttr(
							"nftower_workspace_participant.foo", "member_id", regexp.MustCompile("^[0-9]+$")),
					resource.TestCheckResourceAttr(
						"nftower_workspace_participant.foo", "role", "view"),
				),
			},
		},
	})
}

const testAccResourceWorkspaceParticipant_email = `
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
  email = nftower_organization_member.foo.email
}
`

func TestAccResourceWorkspaceParticipant_maintain(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_workspace_participant",
				Config:       template.ParseRandName(testAccResourceWorkspaceParticipant_maintain),
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

const testAccResourceWorkspaceParticipant_maintain = `
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

func TestAccResourceWorkspaceParticipant_updateRole(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_workspace_participant",
				Config:       template.ParseRandName(testAccResourceWorkspaceParticipant_updateRole_1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"nftower_workspace_participant.foo", "email", regexp.MustCompile("^tf-acceptance-[0-9]+@example.com")),
					resource.TestCheckResourceAttr(
						"nftower_workspace_participant.foo", "role", "launch"),
				),
			},
			{
				ResourceName: "nftower_workspace_participant",
				Config:       template.ParseRandName(testAccResourceWorkspaceParticipant_updateRole_2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"nftower_workspace_participant.foo", "role", "maintain"),
				),
			},
		},
	})
}

const testAccResourceWorkspaceParticipant_updateRole_1 = `
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

const testAccResourceWorkspaceParticipant_updateRole_2 = `
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
