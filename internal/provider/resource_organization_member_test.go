package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/healx/terraform-provider-nftower/internal/template"
)

func TestAccResourceOrganizationMember(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_organization_member",
				Config:       template.ParseRandName(testAccResourceOrganizationMember),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"nftower_organization_member.foo", "email", regexp.MustCompile("^tf-acceptance-[0-9]+@example.com")),
					resource.TestMatchResourceAttr(
						"nftower_organization_member.foo", "user_name", regexp.MustCompile("^tf-acceptance-[0-9]+$")),
					resource.TestCheckResourceAttr(
						"nftower_organization_member.foo", "role", "member"),
				),
			},
		},
	})
}

const testAccResourceOrganizationMember = `
resource "nftower_organization_member" "foo" {
  email = "tf-acceptance-{{.randName}}@example.com"
}
`

func TestAccResourceOrganizationMemberOwner(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_organization_member",
				Config:       template.ParseRandName(testAccResourceOrganizationMemberOwner),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"nftower_organization_member.foo", "email", regexp.MustCompile("^tf-acceptance-[0-9]+@example.com")),
					resource.TestMatchResourceAttr(
						"nftower_organization_member.foo", "user_name", regexp.MustCompile("^tf-acceptance-[0-9]+$")),
					resource.TestCheckResourceAttr(
						"nftower_organization_member.foo", "role", "owner"),
				),
			},
		},
	})
}

const testAccResourceOrganizationMemberOwner = `
resource "nftower_organization_member" "foo" {
  email = "tf-acceptance-{{.randName}}@example.com"
  role  = "owner"
}
`

func TestAccResourceOrganizationMemberUpdateRole(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_organization_member",
				Config:       template.ParseRandName(testAccResourceOrganizationMemberUpdateRole_1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"nftower_organization_member.foo", "email", regexp.MustCompile("^tf-acceptance-[0-9]+@example.com")),
					resource.TestMatchResourceAttr(
						"nftower_organization_member.foo", "user_name", regexp.MustCompile("^tf-acceptance-[0-9]+$")),
					resource.TestCheckResourceAttr(
						"nftower_organization_member.foo", "role", "member"),
				),
			},
			{
				ResourceName: "nftower_organization_member",
				Config:       template.ParseRandName(testAccResourceOrganizationMemberUpdateRole_2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"nftower_organization_member.foo", "role", "owner"),
				),
			},
		},
	})
}

const testAccResourceOrganizationMemberUpdateRole_1 = `
resource "nftower_organization_member" "foo" {
  email = "tf-acceptance-{{.randName}}@example.com"
  role  = "member"
}
`

const testAccResourceOrganizationMemberUpdateRole_2 = `
resource "nftower_organization_member" "foo" {
  email = "tf-acceptance-{{.randName}}@example.com"
  role  = "owner"
}
`
