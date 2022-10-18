package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/healx/terraform-provider-nftower/internal/client"
)

func resourceDataset() *schema.Resource {
	return &schema.Resource{
		Description: "A workspace inside a tower organization.",

		CreateContext: resourceDatasetCreate,
		ReadContext:   resourceDatasetRead,
		UpdateContext: resourceDatasetUpdate,
		DeleteContext: resourceDatasetDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description:  "The name of the dataset. Only alphanumeric characters and dashes are allowed.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"description": {
				Description: "The description of the dataset.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"workspace_id": {
				Description: "The id of the workspace in which to create the dataset.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"date_created": {
				Description: "The datetime the dataset was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_updated": {
				Description: "The last updated datetime of the dataset.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceDatasetCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	id, err := client.CreateDataset(
		ctx,
		d.Get("workspace_id").(string),
		d.Get("name").(string),
		d.Get("description").(string),
	)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	return resourceDatasetRead(ctx, d, meta)
}

func resourceDatasetRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	dataset, err := client.GetDataset(ctx, d.Get("workspace_id").(string), d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if dataset == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", dataset["name"].(string))

	if description, ok := dataset["description"].(string); ok {
		d.Set("description", description)
	} else {
		d.Set("description", nil)
	}

	d.Set("date_created", dataset["dateCreated"].(string))
	d.Set("last_updated", dataset["lastUpdated"].(string))

	return nil
}

func resourceDatasetUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	err := client.UpdateDataset(
		ctx,
		d.Get("workspace_id").(string),
		d.Id(),
		d.Get("name").(string),
		d.Get("description").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceDatasetRead(ctx, d, meta)
}

func resourceDatasetDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	err := client.DeleteDataset(
		ctx,
		d.Get("workspace_id").(string),
		d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
