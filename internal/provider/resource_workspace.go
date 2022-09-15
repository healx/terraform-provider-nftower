package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/healx/terraform-provider-nftower/internal/client"
)

func resourceWorkspace() *schema.Resource {
	return &schema.Resource{
		Description: "A workspace inside a tower organization.",

		CreateContext: resourceWorkspaceCreate,
		ReadContext:   resourceWorkspaceRead,
		UpdateContext: resourceWorkspaceUpdate,
		DeleteContext: resourceWorkspaceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description:  "The name of the workspace. Only alphanumeric characters and dashes are allowed.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 40),
			},
			"full_name": {
				Description:  "The full name of the workspace. Spaces and other characters are allowed.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"description": {
				Description:  "The description of the workspace.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1000),
			},
			"visibility": {
				Description:  "The visiblity of the workspace. Can be PRIVATE, SHARED or PUBLIC.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "PRIVATE",
				ValidateFunc: validation.StringInSlice([]string{"PRIVATE", "SHARED", "PUBLIC"}, false),
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

func resourceWorkspaceCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	id, err := client.CreateWorkspace(
		ctx,
		d.Get("name").(string),
		d.Get("full_name").(string),
		d.Get("description").(string),
		d.Get("visibility").(string),
	)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", id))

	return resourceWorkspaceRead(ctx, d, meta)
}

func resourceWorkspaceRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	workspaceId, _ := strconv.ParseInt(d.Id(), 10, 64)
	workspace, err := client.GetWorkspace(ctx, workspaceId)

	if err != nil {
		return diag.FromErr(err)
	}

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

func resourceWorkspaceUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	workspaceId, _ := strconv.ParseInt(d.Id(), 10, 64)
	err := client.UpdateWorkspace(ctx, workspaceId, d.Get("full_name").(string), d.Get("description").(string), d.Get("visibility").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceWorkspaceRead(ctx, d, meta)
}

func resourceWorkspaceDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	workspaceId, _ := strconv.ParseInt(d.Id(), 10, 64)
	err := client.DeleteWorkspace(ctx, workspaceId)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
