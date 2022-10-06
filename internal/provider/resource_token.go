package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/healx/terraform-provider-nftower/internal/client"
)

func resourceToken() *schema.Resource {
	return &schema.Resource{
		Description: "An api token.",

		CreateContext: resourceTokenCreate,
		ReadContext:   resourceTokenRead,
		DeleteContext: resourceTokenDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the action",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"token": {
				Description: "The token.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"date_created": {
				Description: "The datetime that the token was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceTokenCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client.TowerClient)

	id, token, err := c.CreateToken(ctx, d.Get("name").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	d.Set("token", token)

	return resourceTokenRead(ctx, d, meta)
}

func resourceTokenRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client.TowerClient)

	token, err := c.GetToken(ctx, d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if token == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", token["name"].(string))
	d.Set("date_created", token["dateCreated"].(string))

	return nil
}

func resourceTokenDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client.TowerClient)

	err := c.DeleteToken(ctx, d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
