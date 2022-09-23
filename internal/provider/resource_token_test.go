package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/healx/terraform-provider-nftower/internal/template"
)

func TestAccResourceToken(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_token",
				Config:       template.ParseRandName(testAccResourceToken),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"nftower_token.foo", "token", regexp.MustCompile("^[A-Za-z0-9]{32,}$")),
					resource.TestMatchResourceAttr(
						"nftower_token.foo", "name", regexp.MustCompile("^tf-acceptance-[0-9]+$")),
					resource.TestMatchResourceAttr(
						"nftower_token.foo", "date_created", regexp.MustCompile("^[0-9-:TZ]+")),
				),
			},
		},
	})
}

const testAccResourceToken = `
resource "nftower_token" "foo" {
  name = "tf-acceptance-{{.randName}}"
}
`
