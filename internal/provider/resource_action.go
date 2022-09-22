package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/healx/terraform-provider-nftower/internal/client"
)

func resourceAction() *schema.Resource {
	return &schema.Resource{
		Description: "A workspace inside a tower organization.",

		CreateContext: resourceActionCreate,
		ReadContext:   resourceActionRead,
		UpdateContext: resourceActionUpdate,
		DeleteContext: resourceActionDelete,

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Description: "The id of the workspace in which the action should be created.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "The name of the action",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"source": {
				Description:  "The source of the event. Can be github or tower",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"github", "tower"}, false),
			},
			"compute_environment_id": {
				Description: "The id of the compute environment to use.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"pipeline": {
				Description: `A Git repository name or URL e.g., "nextflow-io/hello" or "https://github.com/nextflow-io/hello". Private repositories require credentials. Local repositories are supported using the "file:" prefix followed by the repository path`,
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"work_dir": {
				Description: "The bucket path where the pipeline scratch data is stored. When only the bucket name is specified, Tower will automatically create a scratch sub-folder.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"revision": {
				Description: "A valid repository commit Id, tag or branch name",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"pre_run_script": {
				Type:        schema.TypeString,
				Description: "This is an optional Bash script that's executed in the same environment where Nextflow runs just before the pipeline is launched. It can useful to stage input data or similar tasks.",
				Optional:    true,
			},
			"post_run_script": {
				Type:        schema.TypeString,
				Description: "This is an optional Bash script that's executed in the same environment where Nextflow runs immediately after the pipeline completion. The script is executed either the pipeline completes successfully or with an error condition. The error condition can be verified using the environment variable NXF_EXIT_STATUS. It can useful to copy result data or similar tasks.",
				Optional:    true,
			},
			"config_profiles": {
				Description: "Select one or more configuration profile names you want to use for this pipeline execution. The profile must be defined in the nextflow.config file included in the pipeline repository.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"pipeline_parameters": {
				Description: "You can specify here any pipeline parameters using either JSON or YML formatted content. This equivalent to the Nextflow -params-file option.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"nextflow_config": {
				Description: "Additional Nextflow config settings can be provided in the above field. These settings will be included in the nextflow.config file for this execution.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"tower_config": {
				Description: "Additional Tower config settings can be provided in the above field. These settings will override the tower.yml file for this execution.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"main_script": {
				Description: "Specify the pipeline main script file if different from `main.nf`",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"workflow_entry_name": {
				Description: "Specify the main workflow name to be executed when using DLS2 syntax",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"schema_name": {
				Description: "Schema name",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"workspace_secrets": {
				Description: "A list of named pipeline secrets required by the pipeline execution. Those secrets must be defined in the launching workspace.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"launch_id": {
				Description: "The id of the launch configuration for the action.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"status": {
				Description: "The status of the action. Can be ACTIVE or PAUSED",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"hook_url": {
				Description: "The url to trigger the action.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"date_created": {
				Description: "The datetime that the action was created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"last_updated": {
				Description: "The datetime that the action was last updated.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceActionCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client.TowerClient)

	id, err := c.CreateAction(
		ctx,
		d.Get("workspace_id").(string),
		d.Get("name").(string),
		d.Get("source").(string),
		d.Get("compute_environment_id").(string),
		d.Get("pipeline").(string),
		d.Get("work_dir").(string),
		d.Get("revision").(string),
		d.Get("pre_run_script").(string),
		d.Get("post_run_script").(string),
		d.Get("config_profiles").(string),
		d.Get("pipeline_parameters").(string),
		d.Get("nextflow_config").(string),
		d.Get("tower_config").(string),
		d.Get("main_script").(string),
		d.Get("workflow_entry_name").(string),
		d.Get("schema_name").(string),
		d.Get("workspace_secrets").([]interface{}))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	return resourceActionRead(ctx, d, meta)
}

func resourceActionRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client.TowerClient)

	action, err := c.GetAction(
		ctx,
		d.Get("workspace_id").(string),
		d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", action["name"].(string))
	d.Set("source", action["source"].(string))
	d.Set("status", action["status"].(string))
	d.Set("hook_url", action["hookUrl"].(string))
	d.Set("date_created", action["dateCreated"].(string))
	d.Set("last_updated", action["lastUpdated"].(string))

	launch := action["launch"].(map[string]interface{})
	d.Set("pipeline", launch["pipeline"].(string))
	d.Set("work_dir", launch["workDir"].(string))
	d.Set("launch_id", launch["id"].(string))

	computeEnv := launch["computeEnv"].(map[string]interface{})
	d.Set("compute_environment_id", computeEnv["id"].(string))

	if v, ok := launch["revision"].(string); ok {
		d.Set("revision", v)
	}

	if v, ok := launch["preRunScript"].(string); ok {
		d.Set("pre_run_script", v)
	}

	if v, ok := launch["postRunScript"].(string); ok {
		d.Set("post_run_script", v)
	}

	if v, ok := launch["configProfiles"].(string); ok {
		d.Set("config_profiles", v)
	}

	if v, ok := launch["paramsText"].(string); ok {
		d.Set("pipeline_parameters", v)
	}

	if v, ok := launch["configText"].(string); ok {
		d.Set("nextflow_config", v)
	}

	if v, ok := launch["towerConfig"].(string); ok {
		d.Set("tower_config", v)
	}

	if v, ok := launch["mainScript"].(string); ok {
		d.Set("main_script", v)
	}

	if v, ok := launch["entryName"].(string); ok {
		d.Set("workflow_entry_name", v)
	}

	if v, ok := launch["schema_name"].(string); ok {
		d.Set("schema_name", v)
	}

	return nil
}

func resourceActionUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client.TowerClient)

	err := c.UpdateAction(
		ctx,
		d.Get("workspace_id").(string),
		d.Id(),
		d.Get("pipeline").(string),
		d.Get("launch_id").(string),
		d.Get("compute_environment_id").(string),
		d.Get("work_dir").(string),
		d.Get("revision").(string),
		d.Get("pre_run_script").(string),
		d.Get("post_run_script").(string),
		d.Get("config_profiles").(string),
		d.Get("pipeline_parameters").(string),
		d.Get("nextflow_config").(string),
		d.Get("tower_config").(string),
		d.Get("main_script").(string),
		d.Get("workflow_entry_name").(string),
		d.Get("schema_name").(string),
		d.Get("workspace_secrets").([]interface{}))

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceActionRead(ctx, d, meta)
}

func resourceActionDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client.TowerClient)

	err := c.DeleteAction(ctx, d.Get("workspace_id").(string), d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
