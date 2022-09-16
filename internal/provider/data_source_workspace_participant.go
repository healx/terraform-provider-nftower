package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/healx/terraform-provider-nftower/internal/client"
)

func dataSourceWorkspaceParticipant() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "A member who has been granted access to a workspace.",

		ReadContext: dataSourceWorkspaceParticipantRead,

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Description: "The id of the workspace.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"email": {
				Description: "The email of the member.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"member_id": {
				Description: "The id of the member in the organization.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"role": {
				Description: "The role of the participant.",
				Type:        schema.TypeString,
				Computed:    true,
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

func dataSourceWorkspaceParticipantRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	participant, err := client.GetWorkspaceParticipant(
		ctx,
		d.Get("workspace_id").(string),
		d.Get("email").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", int64(participant["participantId"].(float64))))

	if v, ok := participant["firstName"].(string); ok {
		d.Set("first_name", v)
	}
	if v, ok := participant["lastName"].(string); ok {
		d.Set("last_name", v)
	}

	d.Set("role", participant["wspRole"].(string))
	d.Set("member_id", fmt.Sprintf("%d", int64(participant["memberId"].(float64))))

	return nil
}
