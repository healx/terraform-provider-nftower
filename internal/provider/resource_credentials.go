package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/healx/terraform-provider-nftower/internal/client"
)

func resourceCredentials() *schema.Resource {
	return &schema.Resource{
		Description: "A set of credentials for use by Tower.",

		CreateContext: resourceCredentialsCreate,
		ReadContext:   resourceCredentialsRead,
		UpdateContext: resourceCredentialsUpdate,
		DeleteContext: resourceCredentialsDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description:  "The name of the credentials. Only alphanumeric characters and dashes are allowed.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 100),
			},
			"description": {
				Description:  "The description of the credentials.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1000),
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Description: "The id of the workspace in which to create the credentials.",
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
			"aws": {
				Description:   "Stores an AWS IAM access key.",
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				MaxItems:      1,
				ConflictsWith: []string{"github"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key": {
							Type:        schema.TypeString,
							Description: "The AWS access key.",
							Required:    true,
						},
						"secret_key": {
							Type:        schema.TypeString,
							Description: "The AWS secret key.",
							Required:    true,
							Sensitive:   true,
						},
						"assume_role_arn": {
							Type:        schema.TypeString,
							Description: "Arn of a role to assume.",
							Optional:    true,
						},
					},
				},
			},
			"github": {
				Description:   "Stores a github access token.",
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      true,
				MaxItems:      1,
				ConflictsWith: []string{"aws"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:        schema.TypeString,
							Description: "The name of the user that the token belongs.",
							Required:    true,
						},
						"access_token": {
							Type:        schema.TypeString,
							Description: "The personal access token to use to connect to github.",
							Required:    true,
							Sensitive:   true,
						},
						"base_url": {
							Type:        schema.TypeString,
							Description: "The base url when connecting to github. Used for github enterprise on-prem.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func resourceCredentialsCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	towerClient := meta.(*client.TowerClient)
	var err error
	var id string

	if _, ok := d.GetOk("aws"); ok {
		id, err = towerClient.CreateCredentialsAWS(
			ctx,
			d.Get("workspace_id").(string),
			d.Get("name").(string),
			d.Get("description").(string),
			d.Get("aws.0.access_key").(string),
			d.Get("aws.0.secret_key").(string),
			d.Get("aws.0.assume_role_arn").(string),
		)
	} else if _, ok := d.GetOk("github"); ok {
		id, err = towerClient.CreateCredentialsGithub(
			ctx,
			d.Get("workspace_id").(string),
			d.Get("name").(string),
			d.Get("description").(string),
			d.Get("github.0.base_url").(string),
			d.Get("github.0.username").(string),
			d.Get("github.0.access_token").(string),
		)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	return resourceCredentialsRead(ctx, d, meta)
}

func resourceCredentialsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	towerClient := meta.(*client.TowerClient)

	credentials, err := towerClient.GetCredentials(ctx, d.Get("workspace_id").(string), d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", credentials["name"].(string))
	d.Set("description", credentials["description"].(string))
	d.Set("date_created", credentials["dateCreated"].(string))
	d.Set("last_updated", credentials["lastUpdated"].(string))

	keys := credentials["keys"].(map[string]interface{})
	switch credentials["provider"].(string) {
	case "aws":
		if assumeRoleArn, ok := keys["assumeRoleArn"].(string); ok {
			d.Set("aws", []interface{}{
				map[string]interface{}{
					"access_key":      keys["accessKey"].(string),
					"secret_key":      d.Get("aws.0.secret_key").(string),
					"assume_role_arn": assumeRoleArn,
				},
			})
		} else {
			d.Set("aws", []interface{}{
				map[string]interface{}{
					"access_key": keys["accessKey"].(string),
					"secret_key": d.Get("aws.0.secret_key").(string),
				},
			})
		}
	case "github":
		if baseUrl, ok := credentials["baseUrl"].(string); ok {
			d.Set("github", []interface{}{
				map[string]interface{}{
					"username":     keys["username"].(string),
					"access_token": d.Get("github.0.access_token").(string),
					"base_url":     baseUrl,
				},
			})
		} else {
			d.Set("github", []interface{}{
				map[string]interface{}{
					"username":     keys["username"].(string),
					"access_token": d.Get("github.0.access_token").(string),
				},
			})
		}
	default:
		return diag.Errorf("unsupported credentials type %s", credentials["provider"].(string))
	}

	return nil
}

func resourceCredentialsUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	var err error
	towerClient := meta.(*client.TowerClient)

	if _, ok := d.GetOk("aws"); ok {
		err = towerClient.UpdateCredentialsAWS(
			ctx,
			d.Id(),
			d.Get("workspace_id").(string),
			d.Get("description").(string),
			d.Get("aws.0.access_key").(string),
			d.Get("aws.0.secret_key").(string),
			d.Get("aws.0.assume_role_arn").(string),
		)
	} else if _, ok := d.GetOk("github"); ok {
		err = towerClient.UpdateCredentialsGithub(
			ctx,
			d.Get("workspace_id").(string),
			d.Get("name").(string),
			d.Get("description").(string),
			d.Get("github.0.base_url").(string),
			d.Get("github.0.username").(string),
			d.Get("github.0.access_token").(string),
		)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCredentialsRead(ctx, d, meta)
}

func resourceCredentialsDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	workspaceId := d.Get("workspace_id").(string)
	err := client.DeleteCredentials(ctx, workspaceId, d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
