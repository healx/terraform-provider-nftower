package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/healx/terraform-provider-nftower/internal/client"
)

func dataSourceOrganizationMember() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "A member of a tower organization.",

		ReadContext: dataSourceOrganizationMemberRead,

		Schema: map[string]*schema.Schema{
			"email": {
				Description: "The email address of the member.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"role": {
				Description: "The role of the member. Can be owner or member",
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
			"user_name": {
				Description: "The username of the member.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceOrganizationMemberRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)
	
	email := d.Get("email").(string)
	member, err := client.GetOrganizationMember(ctx, email)

	if err != nil {
		return diag.FromErr(err)
	}

	if member == nil {
		return diag.Errorf("unable to find member with email: %s", email)
	}

	d.SetId(fmt.Sprintf("%d", int64(member["memberId"].(float64))))

	if v, ok := member["firstName"].(string); ok {
		d.Set("first_name", v)
	}
	if v, ok := member["lastName"].(string); ok {
		d.Set("last_name", v)
	}

	d.Set("role", member["role"].(string))
	d.Set("user_name", member["userName"].(string))

	return nil
}
