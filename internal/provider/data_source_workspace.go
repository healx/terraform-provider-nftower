package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/healx/terraform-provider-nftower/internal/client"
)

func dataSourceWorkspace() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "A workspace.",

		ReadContext: dataSourceWorkspaceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the workspace.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"full_name": {
				Description: "The full name of the workspace.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "The description of the workspace.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"visibility": {
				Description: "The visiblity of the workspace. Can be PRIVATE or PUBLIC.",
				Type:        schema.TypeString,
				Computed:    true,
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

func dataSourceWorkspaceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	workspace, err := client.GetWorkspaceByName(ctx, d.Get("name").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", int64(workspace["id"].(float64))))

	d.Set("name", workspace["name"].(string))
	d.Set("full_name", workspace["fullName"].(string))

	if description, ok := workspace["description"].(string); ok {
		d.Set("description", description)
	} else {
		d.Set("description", nil)
	}

	d.Set("visibility", workspace["visibility"].(string))
	d.Set("date_created", workspace["dateCreated"].(string))
	d.Set("last_updated", workspace["lastUpdated"].(string))

	return nil
}
