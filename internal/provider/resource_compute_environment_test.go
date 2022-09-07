package provider

import (
	"context"
	"reflect"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/healx/terraform-provider-nftower/internal/client"
	"github.com/healx/terraform-provider-nftower/internal/template"
)

func TestAccResourceComputeEnvironment(t *testing.T) {
	t.Skip("requires real AWS credentials")
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_compute_environment",
				Config:       template.ParseRandName(testAccResourceComputeEnvironmentAWS),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"nftower_compute_environment.foo", "name", "tf-acceptance-aws"),
					resource.TestCheckResourceAttr(
						"nftower_compute_environment.foo", "aws.0.region", "eu-west-1"),
					resource.TestCheckResourceAttr(
						"nftower_compute_environment.foo", "aws.0.compute_queue", "aws-nftower-tf-acc"),
					resource.TestCheckResourceAttr(
						"nftower_compute_environment.foo", "aws.0.head_queue", "aws-nftower-tf-acc"),
					resource.TestCheckResourceAttr(
						"nftower_compute_environment.foo", "aws.0.work_dir", "s3://somebucket/"),
					resource.TestMatchResourceAttr(
						"nftower_compute_environment.foo", "date_created", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestMatchResourceAttr(
						"nftower_compute_environment.foo", "last_updated", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestMatchResourceAttr(
						"nftower_compute_environment.foo", "status", regexp.MustCompile("^(CREATING|AVAILABLE|ERRORED)$")),
				),
			},
		},
	})
}

const testAccResourceComputeEnvironmentAWS = `
resource "nftower_workspace" "foo" {
  name        = "tf-acceptance-{{.randName}}"
  full_name   = "tf acceptance testing environments"

  description = "Created by the nftower terraform provider acceptance tests. Will be deleted shortly"
  visibility  = "PRIVATE"
}

resource "nftower_credentials" "foo" {
  name        = "tf-acceptance-envs-aws"
  description = "tf acceptance testing aws environments"
  workspace_id = nftower_workspace.foo.id

  aws {
	access_key      = "foo"
	secret_key      = "bar"
	assume_role_arn = "baz"
  }
}

resource "nftower_compute_environment" "foo" {
  name           = "tf-acceptance-aws"
  workspace_id   = nftower_workspace.foo.id
  credentials_id = nftower_credentials.foo.id

  aws_batch {
	region        = "eu-west-1"
	compute_queue = "aws-nftower-tf-acc"
	head_queue    = "aws-nftower-tf-acc"
	work_dir      = "s3://somebucket/"
  }
}
`

func TestFlattenEnvironmentVariables(t *testing.T) {
	actual := flattenComputeEnvironmentVariables([]*client.ComputeEnvConfigEnvVar{
		{
			Name:    "headOnly",
			Value:   "foo",
			Head:    true,
			Compute: false,
		},
		{
			Name:    "computeOnly",
			Value:   "bar",
			Head:    false,
			Compute: true,
		},
		{
			Name:    "headAndCompute",
			Value:   "baz",
			Head:    true,
			Compute: true,
		},
	})

	expected := []interface{}{
		map[string]interface{}{
			"name":       "headOnly",
			"value":      "foo",
			"visibility": "HEAD",
		},
		map[string]interface{}{
			"name":       "computeOnly",
			"value":      "bar",
			"visibility": "COMPUTE",
		},
		map[string]interface{}{
			"name":       "headAndCompute",
			"value":      "baz",
			"visibility": "BOTH",
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestFlattenComputeEnvironmentAWSBatchMinimal(t *testing.T) {
	ctx := context.Background()
	actual := flattenComputeEnvironmentAWSBatch(ctx, &client.ComputeEnvAWSBatchConfig{
		Region:       "eu-west-1",
		ComputeQueue: "test-queue-compute",
		HeadQueue:    "test-queue-head",
		CliPath:      "/opt/conda/bin/aws",
		WorkDir:      "s3://somebucket/",
		Environment: []*client.ComputeEnvConfigEnvVar{
			{
				Name:    "foo",
				Value:   "bar",
				Head:    true,
				Compute: true,
			},
		},
	})

	expected := []interface{}{
		map[string]interface{}{
			"region":        "eu-west-1",
			"compute_queue": "test-queue-compute",
			"head_queue":    "test-queue-head",
			"cli_path":      "/opt/conda/bin/aws",
			"work_dir":      "s3://somebucket/",
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestFlattenComputeEnvironmentAWSBatchComplete(t *testing.T) {
	ctx := context.Background()
	actual := flattenComputeEnvironmentAWSBatch(ctx, &client.ComputeEnvAWSBatchConfig{
		Region:       "eu-west-1",
		ComputeQueue: "test-queue-compute",
		HeadQueue:    "test-queue-head",
		CliPath:      "/opt/conda/bin/aws",
		WorkDir:      "s3://somebucket/",
		Environment: []*client.ComputeEnvConfigEnvVar{
			{
				Name:    "foo",
				Value:   "bar",
				Head:    true,
				Compute: true,
			},
		},
		HeadJobRole:     "some-role-head",
		ComputeJobRole:  "some-role-compute",
		ExecutionRole:   "some-exec-role",
		HeadJobCpus:     2,
		HeadJobMemoryMb: 2048,
		PreRunScript:    "echo \"foo\"",
		PostRunScript:   "echo \"bar\"",
	})

	expected := []interface{}{
		map[string]interface{}{
			"region":             "eu-west-1",
			"compute_queue":      "test-queue-compute",
			"head_queue":         "test-queue-head",
			"cli_path":           "/opt/conda/bin/aws",
			"work_dir":           "s3://somebucket/",
			"head_job_role":      "some-role-head",
			"compute_job_role":   "some-role-compute",
			"execution_role":     "some-exec-role",
			"head_job_cpus":      2,
			"head_job_memory_mb": 2048,
			"pre_run_script":     "echo \"foo\"",
			"post_run_script":    "echo \"bar\"",
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}
