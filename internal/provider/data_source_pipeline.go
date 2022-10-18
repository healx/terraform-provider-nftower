package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/healx/terraform-provider-nftower/internal/client"
)

func dataSourcePipeline() *schema.Resource {
	return &schema.Resource{
		Description: "A pipeline.",

		ReadContext: dataSourcePipelineRead,

		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Description: "The id of the workspace to which the pipeline belongs.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "The name of the pipeline.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "A description of the pipeline.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"compute_environment_id": {
				Description: "The id of the compute environment to use.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"pipeline": {
				Description: `A Git repository name or URL e.g., "nextflow-io/hello" or "https://github.com/nextflow-io/hello". Private repositories require credentials. Local repositories are supported using the "file:" prefix followed by the repository path`,
				Type:        schema.TypeString,
				Computed: true,
			},
			"work_dir": {
				Description: "The bucket path where the pipeline scratch data is stored. When only the bucket name is specified, Tower will automatically create a scratch sub-folder.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"revision": {
				Description: "A valid repository commit Id, tag or branch name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"pre_run_script": {
				Type:        schema.TypeString,
				Description: "This is an optional Bash script that's executed in the same environment where Nextflow runs just before the pipeline is launched. It can useful to stage input data or similar tasks.",
				Computed:    true,
			},
			"post_run_script": {
				Type:        schema.TypeString,
				Description: "This is an optional Bash script that's executed in the same environment where Nextflow runs immediately after the pipeline completion. The script is executed either the pipeline completes successfully or with an error condition. The error condition can be verified using the environment variable NXF_EXIT_STATUS. It can useful to copy result data or similar tasks.",
				Computed:    true,
			},
			"config_profiles": {
				Description: "A list of one or more configuration profile names you want to use for this pipeline execution. The profile must be defined in the nextflow.config file included in the pipeline repository.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"pipeline_parameters": {
				Description: "You can specify here any pipeline parameters using either JSON or YML formatted content. This equivalent to the Nextflow -params-file option.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"nextflow_config": {
				Description: "Additional Nextflow config settings can be provided in the above field. These settings will be included in the nextflow.config file for this execution.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"tower_config": {
				Description: "Additional Tower config settings can be provided in the above field. These settings will override the tower.yml file for this execution.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"main_script": {
				Description: "Specify the pipeline main script file if different from `main.nf`",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"workflow_entry_name": {
				Description: "Specify the main workflow name to be executed when using DLS2 syntax",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"schema_name": {
				Description: "Schema name",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"workspace_secrets": {
				Description: "A list of named pipeline secrets required by the pipeline execution. Those secrets must be defined in the launching workspace.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourcePipelineRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	pipeline, err := client.GetPipelineByName(ctx, d.Get("workspace_id").(string), d.Get("name").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	if pipeline == nil {
		return diag.Errorf("unable to find pipeline with name: %s", d.Get("name").(string))
	}

	d.SetId(fmt.Sprintf("%d", int64(pipeline["pipelineId"].(float64))))

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

	return nil

	return nil
}
