package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/healx/terraform-provider-nftower/internal/template"
)

func TestAccResourceWorkspace(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_workspace",
				Config:       template.ParseRandName(testAccResourceWorkspace),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"nftower_workspace.foo", "name", regexp.MustCompile("^tf-acceptance-[0-9]+$")),
					resource.TestCheckResourceAttr(
						"nftower_workspace.foo", "full_name", "tf acceptance testing workspace"),
					resource.TestCheckResourceAttr(
						"nftower_workspace.foo", "visibility", "PRIVATE"),
					resource.TestMatchResourceAttr(
						"nftower_workspace.foo", "date_created", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestMatchResourceAttr(
						"nftower_workspace.foo", "last_updated", regexp.MustCompile("^[0-9-:TZ]+")),
				),
			},
		},
	})
}

const testAccResourceWorkspace = `
resource "nftower_workspace" "foo" {
  name        = "tf-acceptance-{{.randName}}"
  full_name   = "tf acceptance testing workspace"
  description = "Created by the nftower terraform provider acceptance tests. Will be deleted shortly"
  visibility  = "PRIVATE"
}
`
