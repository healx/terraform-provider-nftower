package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/healx/terraform-provider-nftower/internal/client"
)

func resourceComputeEnvironment() *schema.Resource {
	return &schema.Resource{
		Description: "A workspace inside a tower organization.",

		CreateContext: resourceComputeEnvironmentCreate,
		ReadContext:   resourceComputeEnvironmentRead,
		DeleteContext: resourceComputeEnvironmentDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description:  "The name of the environment. Only alphanumeric characters and dashes are allowed.",
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 40),
			},
			"description": {
				Description:  "The description of the environment.",
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(0, 1000),
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Description: "The id of the workspace in which to create the environment.",
				Required:    true,
				ForceNew:    true,
			},
			"credentials_id": {
				Type:        schema.TypeString,
				Description: "The id of the credentials to use for the environment.",
				Required:    true,
				ForceNew:    true,
			},
			"status": {
				Description: "The status of the workspace. Can be CREATING, AVAILABLE, ERRORED or INVALID.",
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
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Description: "The AWS region name where the environment lives.",
							Required:    true,
							ForceNew:    true,
						},
						"compute_queue": {
							Type:        schema.TypeString,
							Description: "The default Batch queue to which Nextflow will submit job executions. This can be overwritten via the usual Nextflow config.",
							Required:    true,
							ForceNew:    true,
						},
						"head_queue": {
							Type:        schema.TypeString,
							Description: "The Batch queue that will run the Nextflow application. A queue that does not use spot instances is expected.",
							Required:    true,
							ForceNew:    true,
						},
						"cli_path": {
							Type:        schema.TypeString,
							Description: "Nextflow requires the AWS CLI tool to be installed in the Ec2 instances launched by Batch. Use this field to specify the path where the tool is located. It must start with a '/' and terminate with the '/bin/aws' suffix.",
							Optional:    true,
							ForceNew:    true,
							Default:     "/home/ec2-user/miniconda/bin/aws",
						},
						"work_dir": {
							Type:        schema.TypeString,
							Description: "Either an S3 bucket path, a FSx directory path or a EFS directory path. The S3 bucket should be located in the same region as the one chosen previously.",
							Required:    true,
							ForceNew:    true,
						},
						"compute_job_role": {
							Type:        schema.TypeString,
							Description: "IAM role to fine-grained control permissions for jobs submitted by Nextflow.",
							Optional:    true,
							ForceNew:    true,
						},
						"execution_role": {
							Type:        schema.TypeString,
							Description: "The execution role grants the Amazon ECS container used by Batch the permission to make API calls on your behalf. This field is only required if the pipeline launched with this compute environment needs to access secrets stored in this workspace. If you are not using secrets you can ignore this field. See \"Required IAM permissions for AWS Batch secrets\" documentation for more details.",
							Optional:    true,
							ForceNew:    true,
						},
						"head_job_role": {
							Type:        schema.TypeString,
							Description: "IAM role to fine-grained control permissions for the Nextflow runner job.",
							Optional:    true,
							ForceNew:    true,
						},
						"pre_run_script": {
							Type:        schema.TypeString,
							Description: "This is an optional Bash script that's executed in the same environment where Nextflow runs just before the pipeline is launched. It can useful to stage input data or similar tasks.",
							Optional:    true,
							ForceNew:    true,
						},
						"post_run_script": {
							Type:        schema.TypeString,
							Description: "This is an optional Bash script that's executed in the same environment where Nextflow runs immediately after the pipeline completion. The script is executed either the pipeline completes successfully or with an error condition. The error condition can be verified using the environment variable NXF_EXIT_STATUS. It can useful to copy result data or similar tasks.",
							Optional:    true,
							ForceNew:    true,
						},
						"head_job_cpus": {
							Type:        schema.TypeInt,
							Description: "The number of CPUs to be allocated for the Nextflow runner job.",
							Optional:    true,
							ForceNew:    true,
						},
						"head_job_memory_mb": {
							Type:        schema.TypeInt,
							Description: "The number of MiB of memory reserved for the Nextflow runner job.",
							Optional:    true,
							ForceNew:    true,
						},
					},
				},
			},
			"environment_variable": {
				Type:        schema.TypeList,
				Description: "A List of environment variables that can be included for head or compute jobs.",
				Optional:    true,
				ForceNew:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the environment variable must contain only alphanumeric, dash and underscore characters, and cannot begin with a number.",
							Required:    true,
							ForceNew:    true,
						},
						"value": {
							Type:        schema.TypeString,
							Description: "The value of the environment variable.",
							Required:    true,
							ForceNew:    true,
						},
						"visibility": {
							Type:         schema.TypeString,
							Description:  "Which jobs this environment variable should be available to, can be HEAD, COMPUTE or BOTH.",
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice([]string{"COMPUTE", "HEAD", "BOTH"}, false),
						},
					},
				},
			},
		},
	}
}

func resourceComputeEnvironmentCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	tower_client := meta.(*client.TowerClient)
	var err error
	var id string

	if _, ok := d.GetOk("aws_batch"); ok {
		id, err = tower_client.CreateAWSBatchComputeEnv(
			ctx,
			d.Get("workspace_id").(string),
			d.Get("name").(string),
			d.Get("description").(string),
			d.Get("credentials_id").(string),
			expandComputeEnvironmentAWSBatch(ctx, d),
		)
	} else if _, ok := d.GetOk("lsf_platform"); ok {
		id, err = tower_client.CreateLSFPlatformComputeEnv(
			ctx,
			d.Get("workspace_id").(string),
			d.Get("name").(string),
			d.Get("description").(string),
			d.Get("credentials_id").(string),
			expandComputeEnvironmentLSFPlatform(ctx, d),
		)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)

	return resourceComputeEnvironmentRead(ctx, d, meta)
}

func resourceComputeEnvironmentRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	towerClient := meta.(*client.TowerClient)

	computeEnv, err := towerClient.GetComputeEnv(ctx, d.Get("workspace_id").(string), d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if computeEnv == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", computeEnv["name"].(string))
	d.Set("description", computeEnv["description"].(string))
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

func resourceComputeEnvironmentDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*client.TowerClient)

	err := client.DeleteComputeEnv(ctx, d.Get("workspace_id").(string), d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func expandComputeEnvironmentAWSBatch(ctx context.Context, d *schema.ResourceData) *client.ComputeEnvAWSBatchConfig {
	awsBatchConfig := &client.ComputeEnvAWSBatchConfig{
		Region:       d.Get("aws_batch.0.region").(string),
		ComputeQueue: d.Get("aws_batch.0.compute_queue").(string),
		HeadQueue:    d.Get("aws_batch.0.head_queue").(string),
		CliPath:      d.Get("aws_batch.0.cli_path").(string),
		WorkDir:      d.Get("aws_batch.0.work_dir").(string),
		Environment:  expandComputeEnvironmentVariables(d),
	}

	if v, ok := d.GetOk("aws_batch.0.compute_job_role"); ok {
		awsBatchConfig.ComputeJobRole = v.(string)
	}

	if v, ok := d.GetOk("aws_batch.0.head_job_role"); ok {
		awsBatchConfig.HeadJobRole = v.(string)
	}

	if v, ok := d.GetOk("aws_batch.0.head_job_cpus"); ok {
		awsBatchConfig.HeadJobCpus = v.(int)
	}

	if v, ok := d.GetOk("aws_batch.0.head_job_memory_mb"); ok {
		awsBatchConfig.HeadJobMemoryMb = v.(int)
	}

	if v, ok := d.GetOk("aws_batch.0.execution_role"); ok {
		awsBatchConfig.ExecutionRole = v.(string)
	}

	if v, ok := d.GetOk("aws_batch.0.pre_run_script"); ok {
		awsBatchConfig.PreRunScript = v.(string)
	}

	if v, ok := d.GetOk("aws_batch.0.post_run_script"); ok {
		awsBatchConfig.PostRunScript = v.(string)
	}

	return awsBatchConfig
}

func flattenComputeEnvironmentAWSBatch(ctx context.Context, config *client.ComputeEnvAWSBatchConfig) []interface{} {

	flattened := map[string]interface{}{
		"region":        config.Region,
		"compute_queue": config.ComputeQueue,
		"head_queue":    config.HeadQueue,
		"cli_path":      config.CliPath,
		"work_dir":      config.WorkDir,
	}

	if config.ComputeJobRole != "" {
		flattened["compute_job_role"] = config.ComputeJobRole
	}

	if config.HeadJobRole != "" {
		flattened["head_job_role"] = config.HeadJobRole
	}

	if config.HeadJobCpus != 0 {
		flattened["head_job_cpus"] = config.HeadJobCpus
	}

	if config.HeadJobMemoryMb != 0 {
		flattened["head_job_memory_mb"] = config.HeadJobMemoryMb
	}

	if config.ExecutionRole != "" {
		flattened["execution_role"] = config.ExecutionRole
	}

	if config.PreRunScript != "" {
		flattened["pre_run_script"] = config.PreRunScript
	}

	if config.PostRunScript != "" {
		flattened["post_run_script"] = config.PostRunScript
	}

	v := make([]interface{}, 1)
	v[0] = flattened

	return v
}

func expandComputeEnvironmentLSFPlatform(ctx context.Context, d *schema.ResourceData) *client.ComputeEnvLSFPlatformConfig {
	lsfPlatformConfig := &client.ComputeEnvLSFPlatformConfig{
		WorkDir:                 d.Get("lsf_platform.0.workDir").(string),
		LaunchDir:               d.Get("lsf_platform.0.launchDir").(string),
		UserName:                d.Get("lsf_platform.0.userName").(string),
		HostName:                d.Get("lsf_platform.0.hostName").(string),
		HeadQueue:               d.Get("lsf_platform.0.headQueue").(string),
		ComputeQueue:            d.Get("lsf_platform.0.computeQueue").(string),
		HeadJobOptions:          d.Get("lsf_platform.0.headJobOptions").(string),
		PropagateHeadJobOptions: d.Get("lsf_platform.0.propagateHeadJobOptions").(bool),
		Environment:             expandComputeEnvironmentVariables(d),
	}

	// port
	if v, ok := d.GetOk("lsf_platform.0.port"); ok {
		lsfPlatformConfig.Port = v.(int)
	}

	// maxQueueSize
	if v, ok := d.GetOk("lsf_platform.0.maxQueueSize"); ok {
		lsfPlatformConfig.MaxQueueSize = v.(int)
	}

	// preRunScript
	if v, ok := d.GetOk("lsf_platform.0.preRunScript"); ok {
		lsfPlatformConfig.PreRunScript = v.(string)
	}

	// postRunScript
	if v, ok := d.GetOk("lsf_platform.0.postRunScript"); ok {
		lsfPlatformConfig.PostRunScript = v.(string)
	}

	// unitForLimits
	if v, ok := d.GetOk("lsf_platform.0.unitForLimits"); ok {
		lsfPlatformConfig.UnitForLimits = v.(string)
	}

	// perJobMemLimits
	if v, ok := d.GetOk("lsf_platform.0.perJobMemLimits"); ok {
		lsfPlatformConfig.PerJobMemLimit = v.(bool)
	}

	// perTaskReserve
	if v, ok := d.GetOk("lsf_platform.0.perTaskReserve"); ok {
		lsfPlatformConfig.PerTaskReserve = v.(bool)
	}

	return lsfPlatformConfig
}

func flattenComputeEnvironmentLSFPlatform(ctx context.Context, config *client.ComputeEnvLSFPlatformConfig) []interface{} {

	flattened := map[string]interface{}{
		"work_dir":                   config.WorkDir,
		"launch_dir":                 config.LaunchDir,
		"user_name":                  config.UserName,
		"host_name":                  config.HostName,
		"head_queue":                 config.HeadQueue,
		"compute_queue":              config.ComputeQueue,
		"head_job_options":           config.HeadJobOptions,
		"propagate_head_job_options": config.PropagateHeadJobOptions,
		"per_job_mem_limit":          config.PerJobMemLimit,
		"per_task_reserve":           config.PerTaskReserve,
	}

	if config.Port != 0 {
		flattened["port"] = config.Port
	}

	if config.MaxQueueSize != 0 {
		flattened["max_queue_size"] = config.MaxQueueSize
	}

	if config.PreRunScript != "" {
		flattened["pre_run_script"] = config.PreRunScript
	}

	if config.PostRunScript != "" {
		flattened["post_run_script"] = config.PostRunScript
	}

	if config.UnitForLimits != "" {
		flattened["unit_for_limits"] = config.UnitForLimits
	}

	v := make([]interface{}, 1)
	v[0] = flattened

	return v
}

func expandComputeEnvironmentVariables(d *schema.ResourceData) []*client.ComputeEnvConfigEnvVar {
	vars := d.Get("environment_variable").([]interface{})
	envVars := make([]*client.ComputeEnvConfigEnvVar, len(vars))
	for i, raw := range vars {
		data := raw.(map[string]interface{})
		visibility := data["visibility"].(string)
		envVars[i] = &client.ComputeEnvConfigEnvVar{
			Name:    data["name"].(string),
			Value:   data["value"].(string),
			Head:    visibility == "HEAD" || visibility == "BOTH",
			Compute: visibility == "COMPUTE" || visibility == "BOTH",
		}
	}
	return envVars
}

func flattenComputeEnvironmentVariables(vars []*client.ComputeEnvConfigEnvVar) []interface{} {
	flattened := make([]interface{}, len(vars))

	for i, v := range vars {
		visibility := "BOTH"

		if v.Head && !v.Compute {
			visibility = "HEAD"
		}

		if !v.Head && v.Compute {
			visibility = "COMPUTE"
		}

		flattened[i] = map[string]interface{}{
			"name":       v.Name,
			"value":      v.Value,
			"visibility": visibility,
		}
	}

	return flattened
}
