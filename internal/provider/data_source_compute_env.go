package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/healx/terraform-provider-nftower/internal/client"
)

func dataSourceComputeEnv() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "A compute environment.",

		ReadContext: dataSourceComputeEnvRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the environment. Only alphanumeric characters and dashes are allowed.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Description: "The id of the workspace in which to create the environment.",
				Required:    true,
			},
			"description": {
				Description: "The description of the environment.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"credentials_id": {
				Type:        schema.TypeString,
				Description: "The id of the credentials to use for this environment.",
				Computed:    true,
			},
			"status": {
				Description: "The status of the workspace. Can be CREATING, AVAILABLE or ERRORED.",
				Type:        schema.TypeString,
				Computed:    true,
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
			"aws_batch": {
				Description: "Configures an AWS Batch compute environment (manual only).",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Description: "The AWS region name where the environment lives.",
							Computed:    true,
						},
						"compute_queue": {
							Type:        schema.TypeString,
							Description: "The default Batch queue to which Nextflow will submit job executions. This can be overwritten via the usual Nextflow config.",
							Computed:    true,
						},
						"head_queue": {
							Type:        schema.TypeString,
							Description: "The Batch queue that will run the Nextflow application. A queue that does not use spot instances is expected.",
							Computed:    true,
						},
						"cli_path": {
							Type:        schema.TypeString,
							Description: "Nextflow requires the AWS CLI tool to be installed in the Ec2 instances launched by Batch. Use this field to specify the path where the tool is located. It must start with a '/' and terminate with the '/bin/aws' suffix.",
							Computed:    true,
						},
						"work_dir": {
							Type:        schema.TypeString,
							Description: "Either an S3 bucket path, a FSx directory path or a EFS directory path. The S3 bucket should be located in the same region as the one chosen previously.",
							Computed:    true,
						},
						"compute_job_role": {
							Type:        schema.TypeString,
							Description: "IAM role to fine-grained control permissions for jobs submitted by Nextflow.",
							Computed:    true,
						},
						"execution_role": {
							Type:        schema.TypeString,
							Description: "The execution role grants the Amazon ECS container used by Batch the permission to make API calls on your behalf. This field is only required if the pipeline launched with this compute environment needs to access secrets stored in this workspace. If you are not using secrets you can ignore this field. See \"Required IAM permissions for AWS Batch secrets\" documentation for more details.",
							Computed:    true,
						},
						"head_job_role": {
							Type:        schema.TypeString,
							Description: "IAM role to fine-grained control permissions for the Nextflow runner job.",
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
						"head_job_cpus": {
							Type:        schema.TypeInt,
							Description: "The number of CPUs to be allocated for the Nextflow runner job.",
							Computed:    true,
						},
						"head_job_memory_mb": {
							Type:        schema.TypeInt,
							Description: "The number of MiB of memory reserved for the Nextflow runner job.",
							Computed:    true,
						},
					},
				},
			},
			"environment_variable": {
				Type:        schema.TypeList,
				Description: "A List of environment variables that can be included for head or compute jobs.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the environment variable must contain only alphanumeric, dash and underscore characters, and cannot begin with a number.",
							Computed:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "The value of the environment variable.",
							Computed:    true,
						},
						"visiblity": {
							Type:        schema.TypeString,
							Description: "Which jobs this environment variable should be available to, can be HEAD, COMPUTE or BOTH.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceComputeEnvRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	towerClient := meta.(*client.TowerClient)

	computeEnv, err := towerClient.GetComputeEnvByName(ctx, d.Get("workspace_id").(string), d.Get("name").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	if computeEnv == nil {
		return diag.Errorf("unable to find compute environment with name: %s", d.Get("name").(string))
	}

	d.SetId(computeEnv["id"].(string))

	d.Set("name", computeEnv["name"].(string))

	if description, ok := computeEnv["description"].(string); ok {
		d.Set("description", description)
	} else {
		d.Set("description", nil)
	}

	d.Set("credentials_id", computeEnv["credentialsId"].(string))
	d.Set("date_created", computeEnv["dateCreated"].(string))
	d.Set("last_updated", computeEnv["lastUpdated"].(string))
	d.Set("status", computeEnv["status"].(string))

	switch computeEnv["platform"].(string) {
	case "aws-batch":
		config := computeEnv["config"].(client.ComputeEnvAWSBatchConfig)
		d.Set("aws_batch", flattenComputeEnvironmentAWSBatch(ctx, &config))
		d.Set("environment_variable", flattenComputeEnvironmentVariables(config.Environment))
	default:
		return diag.Errorf("unsupported platform type: %s", computeEnv["platform"].(string))
	}

	return nil
}
