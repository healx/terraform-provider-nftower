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

func resourceWorkspaceParticipant() *schema.Resource {
	return &schema.Resource{
		Description: "Grants access to a tower workspace.",

		CreateContext: resourceWorkspaceParticipantCreate,
		ReadContext:   resourceWorkspaceParticipantRead,
		UpdateContext: resourceWorkspaceParticipantUpdate,
		DeleteContext: resourceWorkspaceParticipantDelete,

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Description: "The id of the workspace to grant access to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"member_id": {
				Description: "The id of the member in the organization.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"role": {
				Description: "The role of the participant.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "launch",
				ValidateFunc: validation.StringInSlice(
					[]string{"owner", "admin", "maintain", "launch", "view"},
					false),
			},
			"first_name": {
				Description: "The first name of the member.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_name": {
				Description: "The last name of the member.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email": {
				Description: "The email of the member.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceWorkspaceParticipantCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	memberId, _ := strconv.ParseInt(d.Get("member_id").(string), 10, 64)

	id, email, err := client.CreateWorkspaceParticipant(
		ctx,
		d.Get("workspace_id").(string),
		memberId,
		d.Get("role").(string),
	)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", id))
	d.Set("email", email)

	return resourceWorkspaceParticipantRead(ctx, d, meta)
}

func resourceWorkspaceParticipantRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	participant, err := client.GetWorkspaceParticipant(ctx,
		d.Get("workspace_id").(string),
		d.Get("email").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	if v, ok := participant["firstName"].(string); ok {
		d.Set("first_name", v)
	}
	if v, ok := participant["lastName"].(string); ok {
		d.Set("last_name", v)
	}

	d.Set("role", participant["wspRole"].(string))

	return nil
}

func resourceWorkspaceParticipantUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	participantId, _ := strconv.ParseInt(d.Id(), 10, 64)
	err := client.UpdateWorkspaceParticipantRole(
		ctx,
		d.Get("workspace_id").(string),
		participantId,
		d.Get("role").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceWorkspaceParticipantRead(ctx, d, meta)
}

func resourceWorkspaceParticipantDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	participantId, _ := strconv.ParseInt(d.Id(), 10, 64)
	err := client.DeleteWorkspaceParticipant(
		ctx,
		d.Get("workspace_id").(string),
		participantId)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
