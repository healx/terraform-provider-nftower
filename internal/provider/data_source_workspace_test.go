package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/healx/terraform-provider-nftower/internal/template"
)

func TestAccDataSourceWorkspace(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: template.ParseRandName(testAccDataSourceWorkspace),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.nftower_workspace.foo", "name", regexp.MustCompile("^tf-acceptance-[0-9]+$")),
					resource.TestCheckResourceAttr(
						"data.nftower_workspace.foo", "full_name", "tf acceptance testing ds workspace"),
					resource.TestCheckResourceAttr(
						"data.nftower_workspace.foo", "visibility", "PRIVATE"),
					resource.TestMatchResourceAttr(
						"data.nftower_workspace.foo", "date_created", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestMatchResourceAttr(
						"data.nftower_workspace.foo", "last_updated", regexp.MustCompile("^[0-9-:TZ]+")),
				),
			},
		},
	})
}

const testAccDataSourceWorkspace = `
resource "nftower_workspace" "foo" {
  name        = "tf-acceptance-{{.randName}}"
  full_name   = "tf acceptance testing ds workspace"
  description = "Created by the nftower terraform provider acceptance tests. Will be deleted shortly"
  visibility  = "PRIVATE"
}

data "nftower_workspace" "foo" {
  name = nftower_workspace.foo.name
}
`
