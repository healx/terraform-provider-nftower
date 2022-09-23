package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/healx/terraform-provider-nftower/internal/template"
)

func TestAccResourceDatasetVersion_csv(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_dataset_version",
				Config:       template.ParseRandName(testAccResourceDatasetVersion_csv),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"nftower_dataset_version.foo", "file_name", "foo.csv"),
					resource.TestCheckResourceAttrSet(
						"nftower_dataset_version.foo", "contents"),
					resource.TestCheckResourceAttr(
						"nftower_dataset_version.foo", "version", "1"),
					resource.TestCheckResourceAttr(
						"nftower_dataset_version.foo", "media_type", "text/csv"),
					resource.TestMatchResourceAttr(
						"nftower_dataset_version.foo", "last_updated", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestMatchResourceAttr(
						"nftower_dataset_version.foo", "url", regexp.MustCompile("^https://api.tower.nf/workspaces/[0-9]+/datasets/[0-9A-Za-z]+/v/[0-9]+/n/foo.csv$")),
				),
			},
		},
	})
}

const testAccResourceDatasetVersion_csv = `
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

resource "nftower_dataset_version" "foo" {
  dataset_id = nftower_dataset.foo.id
  workspace_id = nftower_workspace.foo.id
  file_name = "foo.csv"
  contents = <<EOF
one,two,three,four
1,2,3,4
EOF
  has_header = true
}
`

func TestAccResourceDatasetVersion_tsv(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_dataset_version",
				Config:       template.ParseRandName(testAccResourceDatasetVersion_tsv),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"nftower_dataset_version.foo", "file_name", "foo.tsv"),
					resource.TestCheckResourceAttrSet(
						"nftower_dataset_version.foo", "contents"),
					resource.TestCheckResourceAttr(
						"nftower_dataset_version.foo", "version", "1"),
					resource.TestCheckResourceAttr(
						"nftower_dataset_version.foo", "media_type", "text/tab-separated-values"),
					resource.TestMatchResourceAttr(
						"nftower_dataset_version.foo", "last_updated", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestMatchResourceAttr(
						"nftower_dataset_version.foo", "url", regexp.MustCompile("^https://api.tower.nf/workspaces/[0-9]+/datasets/[0-9A-Za-z]+/v/[0-9]+/n/foo.tsv$")),
				),
			},
		},
	})
}

const testAccResourceDatasetVersion_tsv = `
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

resource "nftower_dataset_version" "foo" {
  dataset_id = nftower_dataset.foo.id
  workspace_id = nftower_workspace.foo.id
  file_name = "foo.tsv"
  contents = <<EOF
one	two	three	four
1	2	3	4
EOF
  has_header = true
}
`
