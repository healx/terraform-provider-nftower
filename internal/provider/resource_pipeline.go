package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/healx/terraform-provider-nftower/internal/client"
)

func resourcePipeline() *schema.Resource {
	return &schema.Resource{
		Description: "A workflow pipeline.",

		CreateContext: resourcePipelineCreate,
		ReadContext:   resourcePipelineRead,
		UpdateContext: resourcePipelineUpdate,
		DeleteContext: resourcePipelineDelete,

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Description: "The id of the workspace in which the pipeline should be created.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"name": {
				Description: "The name of the pipeline.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Description: "A description of the pipeline.",
				Type:        schema.TypeString,
				Optional:    true,
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
				Description: "A list of one or more configuration profile names you want to use for this pipeline execution. The profile must be defined in the nextflow.config file included in the pipeline repository.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
			"labels": {
				Description: "A set of labels to apply to the triggered pipeline run. Minimum 2 characters.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringLenBetween(2, 1000),
				},
			},
		},
	}
}

func resourcePipelineCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client.TowerClient)

	id, err := c.CreatePipeline(
		ctx,
		d.Get("workspace_id").(string),
		d.Get("name").(string),
		d.Get("description").(string),
		d.Get("compute_environment_id").(string),
		d.Get("pipeline").(string),
		d.Get("work_dir").(string),
		d.Get("revision").(string),
		d.Get("pre_run_script").(string),
		d.Get("post_run_script").(string),
		d.Get("config_profiles").([]interface{}),
		d.Get("pipeline_parameters").(string),
		d.Get("nextflow_config").(string),
		d.Get("tower_config").(string),
		d.Get("main_script").(string),
		d.Get("workflow_entry_name").(string),
		d.Get("schema_name").(string),
		d.Get("workspace_secrets").([]interface{}),
		expandLabels(d.Get("labels").(*schema.Set)))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", id))

	return resourcePipelineRead(ctx, d, meta)
}

func resourcePipelineRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client.TowerClient)

	pipeline, err := c.GetPipeline(ctx, d.Get("workspace_id").(string), d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if pipeline == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", pipeline["name"].(string))
	d.Set("pipeline", pipeline["pipeline"].(string))
	d.Set("work_dir", pipeline["workDir"].(string))

	computeEnv := pipeline["computeEnv"].(map[string]interface{})
	d.Set("compute_environment_id", computeEnv["id"].(string))

	if v, ok := pipeline["revision"].(string); ok {
		d.Set("revision", v)
	}

	if v, ok := pipeline["preRunScript"].(string); ok {
		d.Set("pre_run_script", v)
	}

	if v, ok := pipeline["postRunScript"].(string); ok {
		d.Set("post_run_script", v)
	}

	if v, ok := pipeline["configProfiles"].([]interface{}); ok {
		d.Set("config_profiles", v)
	}

	if v, ok := pipeline["paramsText"].(string); ok {
		d.Set("pipeline_parameters", v)
	}

	if v, ok := pipeline["configText"].(string); ok {
		d.Set("nextflow_config", v)
	}

	if v, ok := pipeline["towerConfig"].(string); ok {
		d.Set("tower_config", v)
	}

	if v, ok := pipeline["mainScript"].(string); ok {
		d.Set("main_script", v)
	}

	if v, ok := pipeline["entryName"].(string); ok {
		d.Set("workflow_entry_name", v)
	}

	if v, ok := pipeline["schema_name"].(string); ok {
		d.Set("schema_name", v)
	}

	if v, ok := pipeline["workspaceSecrets"].([]interface{}); ok {
		d.Set("workspace_secrets", v)
	}

	if v, ok := pipeline["labels"].([]interface{}); ok {
		d.Set("labels", flattenLabels(v))
	}

	return nil
}

func resourcePipelineUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client.TowerClient)

	err := c.UpdatePipeline(
		ctx,
		d.Get("workspace_id").(string),
		d.Id(),
		d.Get("description").(string),
		d.Get("compute_environment_id").(string),
		d.Get("pipeline").(string),
		d.Get("work_dir").(string),
		d.Get("revision").(string),
		d.Get("pre_run_script").(string),
		d.Get("post_run_script").(string),
		d.Get("config_profiles").([]interface{}),
		d.Get("pipeline_parameters").(string),
		d.Get("nextflow_config").(string),
		d.Get("tower_config").(string),
		d.Get("main_script").(string),
		d.Get("workflow_entry_name").(string),
		d.Get("schema_name").(string),
		d.Get("workspace_secrets").([]interface{}),
		expandLabels(d.Get("labels").(*schema.Set)))

	if err != nil {
		return diag.FromErr(err)
	}
	return resourcePipelineRead(ctx, d, meta)
}

func resourcePipelineDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	c := meta.(*client.TowerClient)

	err := c.DeletePipeline(ctx, d.Get("workspace_id").(string), d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
