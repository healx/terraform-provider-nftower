package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/healx/terraform-provider-nftower/internal/client"
)

func dataSourcePipelineSecrets() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "A pipeline-secret for use by Tower.",

		ReadContext: dataSourcePipelineSecretsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the pipeline-secret.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"value": {
				Description: "The value of the pipeline-secret.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourcePipelineSecretsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	towerClient := meta.(*client.TowerClient)

	pipelineSecret, err := towerClient.GetPipelineSecretByName(ctx, d.Get("workspace_id").(string), d.Get("name").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	if pipelineSecret == nil {
		return diag.Errorf("unable to find pipeline-secret with name: %s", d.Get("name").(string))
	}

	d.SetId(pipelineSecret["id"].(string))
	d.Set("name", pipelineSecret["name"].(string))

	d.Set("date_used", pipelineSecret["lastUsed"].(string))

	d.Set("date_created", pipelineSecret["dateCreated"].(string))
	d.Set("last_updated", pipelineSecret["lastUpdated"].(string))

	return nil
}
