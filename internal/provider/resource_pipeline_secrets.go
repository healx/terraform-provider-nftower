package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/healx/terraform-provider-nftower/internal/client"
)

func resourcePipelineSecrets() *schema.Resource {
	return &schema.Resource{
		Description: "A pipeline-secret for use by Tower.",

		CreateContext: resourcePipelineSecretsCreate,
		ReadContext:   resourcePipelineSecretsRead,
		UpdateContext: resourcePipelineSecretsUpdate,
		DeleteContext: resourcePipelineSecretsDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the pipeline-secret.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"value": {
				Description: "The value of the pipeline-secret.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Description: "The id of the workspace in which to create the pipeline-secret.",
				Required:    true,
				ForceNew:    true,
			},
			"date_created": {
				Description: "The datetime the workspace was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_updated": {
				Description: "The last updated datetime of the workspace.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourcePipelineSecretsCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	towerClient := meta.(*client.TowerClient)
	var err error
	var id string

	id, err = towerClient.CreatePipelineSecrets(
		ctx,
		d.Get("workspace_id").(string),
		d.Get("name").(string),
		d.Get("value").(string),
	)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	return resourcePipelineSecretsRead(ctx, d, meta)
}

func resourcePipelineSecretsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	towerClient := meta.(*client.TowerClient)

	pipelineSecret, err := towerClient.GetPipelineSecret(ctx, d.Get("workspace_id").(string), d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if pipelineSecret == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", pipelineSecret["name"].(string))

	d.Set("last_used", pipelineSecret["lastUsed"].(string))

	d.Set("date_created", pipelineSecret["dateCreated"].(string))
	d.Set("last_updated", pipelineSecret["lastUpdated"].(string))

	return nil
}

func resourcePipelineSecretsUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var err error
	towerClient := meta.(*client.TowerClient)

	err = towerClient.UpdatePipelineSecrets(
		ctx,
		d.Id(),
		d.Get("workspace_id").(string),
		d.Get("name").(string),
		d.Get("value").(string),
	)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePipelineSecretsRead(ctx, d, meta)
}

func resourcePipelineSecretsDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	workspaceId := d.Get("workspace_id").(string)
	err := client.DeletePipelineSecrets(ctx, workspaceId, d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
