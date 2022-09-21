package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/healx/terraform-provider-nftower/internal/template"
)

func TestAccResourceDataset_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_dataset",
				Config:       template.ParseRandName(testAccResourceDataset_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"nftower_dataset.foo", "name", regexp.MustCompile("^tf-acceptance-[0-9]+$")),
					resource.TestCheckResourceAttr(
						"nftower_dataset.foo", "description", "tf acceptance testing dataset"),
					resource.TestMatchResourceAttr(
						"nftower_dataset.foo", "date_created", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestMatchResourceAttr(
						"nftower_dataset.foo", "last_updated", regexp.MustCompile("^[0-9-:TZ]+")),
				),
			},
		},
	})
}

const testAccResourceDataset_basic = `
resource "nftower_workspace" "foo" {
  name        = "tf-acceptance-{{.randName}}"
  full_name   = "tf acceptance testing credentials"
	
  description = "Created by the nftower terraform provider acceptance tests. Will be deleted shortly"
  visibility  = "PRIVATE"
}

resource "nftower_dataset" "foo" {
  name        = "tf-acceptance-{{.randName}}"
  description = "tf acceptance testing dataset"
  workspace_id = nftower_workspace.foo.id
}
`
