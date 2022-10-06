package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/healx/terraform-provider-nftower/internal/client"
)

func dataSourceCredentials() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "A set of credentials for use by Tower.",

		ReadContext: dataSourceCredentialsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the credentials. Only alphanumeric characters and dashes are allowed.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the environment.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Description: "The id of the workspace in which to the credentials live.",
				Required:    true,
			},
			"date_created": {
				Description: "The datetime the credentials were created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_updated": {
				Description: "The last updated datetime of the credentials.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"aws": {
				Description: "Stores an AWS IAM access key.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key": {
							Type:        schema.TypeString,
							Description: "The AWS access key.",
							Computed:    true,
						},
						"assume_role_arn": {
							Type:        schema.TypeString,
							Description: "Arn of a role to assume.",
							Computed:    true,
						},
					},
				},
			},
			"github": {
				Description: "Stores a github access token.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:        schema.TypeString,
							Description: "The name of the user that the token belongs.",
							Computed:    true,
						},
						"base_url": {
							Type:        schema.TypeString,
							Description: "The base url when connecting to github. Used for github enterprise on-prem.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceCredentialsRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	towerClient := meta.(*client.TowerClient)

	credentials, err := towerClient.GetCredentialsByName(ctx, d.Get("workspace_id").(string), d.Get("name").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	if credentials == nil {
		return diag.Errorf("unable to find credentials with name: %s", d.Get("name").(string))
	}

	d.SetId(credentials["id"].(string))
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
					"assume_role_arn": assumeRoleArn,
				},
			})
		} else {
			d.Set("aws", []interface{}{
				map[string]interface{}{
					"access_key": keys["accessKey"].(string),
				},
			})
		}
	case "github":
		if baseUrl, ok := credentials["baseUrl"].(string); ok {
			d.Set("github", []interface{}{
				map[string]interface{}{
					"username": keys["username"].(string),
					"base_url": baseUrl,
				},
			})
		} else {
			d.Set("github", []interface{}{
				map[string]interface{}{
					"username": keys["username"].(string),
				},
			})
		}
	default:
		return diag.Errorf("unsupported credentials type %s", credentials["provider"].(string))
	}

	return nil
}
