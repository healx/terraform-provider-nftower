package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/healx/terraform-provider-nftower/internal/client"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"api_url": {
					Type:     schema.TypeString,
					Optional: true,
					DefaultFunc: schema.MultiEnvDefaultFunc([]string{
						"NFTOWER_API_URL",
					}, "https://api.tower.nf"),
				},
				"api_key": {
					Type:     schema.TypeString,
					Required: true,
					DefaultFunc: schema.MultiEnvDefaultFunc([]string{
						"NFTOWER_API_KEY",
					}, nil),
				},
				"organization": {
					Type:     schema.TypeString,
					Required: true,
					DefaultFunc: schema.MultiEnvDefaultFunc([]string{
						"NFTOWER_ORGANIZATION",
					}, nil),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"nftower_workspace":             dataSourceWorkspace(),
				"nftower_compute_environment":   dataSourceComputeEnv(),
				"nftower_credentials":           dataSourceCredentials(),
				"nftower_organization_member":   dataSourceOrganizationMember(),
				"nftower_workspace_participant": dataSourceWorkspaceParticipant(),
				"nftower_pipeline":              dataSourcePipeline(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"nftower_workspace":             resourceWorkspace(),
				"nftower_compute_environment":   resourceComputeEnvironment(),
				"nftower_credentials":           resourceCredentials(),
				"nftower_organization_member":   resourceOrganizationMember(),
				"nftower_workspace_participant": resourceWorkspaceParticipant(),
				"nftower_dataset":               resourceDataset(),
				"nftower_dataset_version":       resourceDatasetVersion(),
				"nftower_action":                resourceAction(),
				"nftower_token":                 resourceToken(),
				"nftower_pipeline":              resourcePipeline(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (any, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (any, diag.Diagnostics) {
		c, err := client.NewTowerClient(ctx,
			p.UserAgent("terraform-provider-nftower", version),
			d.Get("api_key").(string),
			d.Get("api_url").(string),
			d.Get("organization").(string))

		if err != nil {
			return nil, diag.FromErr(err)
		}

		return c, nil
	}
}
