package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/healx/terraform-provider-nftower/internal/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceOrganizationMember() *schema.Resource {
	return &schema.Resource{
		Description: "A member of a tower organization.",

		CreateContext: resourceOrganizationMemberCreate,
		ReadContext:   resourceOrganizationMemberRead,
		UpdateContext: resourceOrganizationMemberUpdate,
		DeleteContext: resourceOrganizationMemberDelete,

		Schema: map[string]*schema.Schema{
			"email": {
				Description: "The email address of the user to add.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"role": {
				Description: "The role of the member. Can be owner or member",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "member",
				ValidateFunc: validation.StringInSlice([]string{"member", "owner"}, false),
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

func resourceOrganizationMemberCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	id, err := client.CreateOrganizationMember(
		ctx,
		d.Get("email").(string),
		d.Get("role").(string),
	)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", id))

	return resourceOrganizationMemberRead(ctx, d, meta)
}

func resourceOrganizationMemberRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	member, err := client.GetOrganizationMember(ctx, d.Get("email").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Trace(ctx, "reading new member", member)

	d.Set("email", member["email"].(string))
	d.Set("user_name", member["userName"].(string))

	if v, ok := member["firstName"].(string); ok {
		d.Set("first_name", v)
	}
	if v, ok := member["lastName"].(string); ok {
		d.Set("last_name", v)
	}

	d.Set("role", member["role"].(string))

	return nil
}

func resourceOrganizationMemberUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	memberId, _ := strconv.ParseInt(d.Id(), 10, 64)
	err := client.UpdateOrganizationMemberRole(ctx, memberId, d.Get("role").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceOrganizationMemberRead(ctx, d, meta)
}

func resourceOrganizationMemberDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	memberId, _ := strconv.ParseInt(d.Id(), 10, 64)
	err := client.DeleteOrganizationMember(ctx, memberId)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
