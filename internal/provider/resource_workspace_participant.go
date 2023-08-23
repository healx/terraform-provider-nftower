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
		Description: "Grants access to a tower workspace. The member must already be added to the organization.",

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
				Description:   "The id of the member in the organization. Specify either member_id or email but not both.",
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"email"},
			},
			"email": {
				Description:   "The email of the member. Specify either member_id or email but not both.",
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"member_id"},
			},
			"role": {
				Description: "The role of the participant.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "view",
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
		},
	}
}

func resourceWorkspaceParticipantCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	var memberId int64

	if v, ok := d.GetOk("member_id"); ok {
		id, err := strconv.ParseInt(v.(string), 10, 64)
		memberId = id

		if err != nil {
			return diag.Errorf("member_id must be a number, got %s", v.(string))
		}
	}

	if v, ok := d.GetOk("email"); ok {
		email := v.(string)
		member, err := client.GetOrganizationMember(ctx, email)

		if err != nil {
			return diag.FromErr(err)
		}

		if member == nil {
			return diag.Errorf("no member found in organization with email %s", email)
		}

		memberId = int64(member["memberId"].(float64))
	}

	if memberId == 0 {
		return diag.Errorf("either member_id or email must be specified.")
	}

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
	d.Set("member_id", fmt.Sprintf("%d", memberId))

	return resourceWorkspaceParticipantRead(ctx, d, meta)
}

func resourceWorkspaceParticipantRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	participant, err := client.GetWorkspaceParticipantByMemberEmail(ctx,
		d.Get("workspace_id").(string),
		d.Get("email").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	if participant == nil {
		d.SetId("")
		return nil
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
