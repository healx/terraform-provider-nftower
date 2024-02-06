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

func TestAccResourceComputeEnvironmentAWS(t *testing.T) {
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

func TestAccResourceComputeEnvironmentLSF(t *testing.T) {
	t.Skip("requires real ssh login node")
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_compute_environment",
				Config:       template.ParseRandName(testAccResourceComputeEnvironmentLSF),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "name", "tf-acceptance-lsf-ssh-cred"),
					resource.TestCheckResourceAttr(
						"nftower_compute_environment.foo", "name", "tf-acceptance-lsf"),
					resource.TestCheckResourceAttr(
						"nftower_compute_environment.foo", "lsf_platform.0.work_dir", "/nextflow/work"),
					resource.TestCheckResourceAttr(
						"nftower_compute_environment.foo", "lsf_platform.0.launch_dir", "/nextflow/launch"),
					resource.TestCheckResourceAttr(
						"nftower_compute_environment.foo", "lsf_platform.0.user_name", "nextflow"),
					resource.TestCheckResourceAttr(
						"nftower_compute_environment.foo", "lsf_platform.0.host_name", "localhost"),
					resource.TestCheckResourceAttr(
						"nftower_compute_environment.foo", "lsf_platform.0.head_queue", "head"),
					resource.TestCheckResourceAttr(
						"nftower_compute_environment.foo", "lsf_platform.0.compute_queue", "compute"),
					resource.TestCheckResourceAttr(
						"nftower_compute_environment.foo", "lsf_platform.0.head_job_options", "-x something"),
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

const testAccResourceComputeEnvironmentLSF = `
resource "nftower_workspace" "foo" {
  name        = "tf-acceptance-{{.randName}}"
  full_name   = "tf acceptance testing environments"

  description = "Created by the nftower terraform provider acceptance tests. Will be deleted shortly"
  visibility  = "PRIVATE"
}

resource "nftower_credentials" "foo" {
  name        = "tf-acceptance-lsf-ssh-cred"
  description = "tf acceptance testing lsf environments"
  workspace_id = nftower_workspace.foo.id

  ssh {
	private_key = <<EOF
	-----BEGIN OPENSSH PRIVATE KEY-----
	b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABlwAAAAdzc2gtcn
	NhAAAAAwEAAQAAAYEA1Ui6IQY+FdeCgGPxiK1kz1Smet1uiydviL4pSGzZkJamhIw3Zf/i
	ccUcH81Oas21fi/sSGKXyMEr1P3qPk3bs25MUpVqS8Mc/2u3grkjL+5BJin5DVpeX+Slzl
	xyrDTCV0PU+jb0vCGoo1Fuiea5/15IEswfWw+lEjJI05qvYNkX2hQA471FaZmQlbYokStQ
	OtJaDjLrbxCjNdjTz4vlRzPdpG8jvlSjT4LiP8M9gtPFRPf/MgXXx+fZUx+b/Ki8FeqTP4
	SbbDhfgtf7Kq5XeFrZOMHa2N7+iaUF3tTV/GNoywLco1sOkpx2pDy8rmoQCJaxcmuUIlji
	ZxHn3Tu8aNPtAJoWYfPtwPQpyHCcMW72vGNfNmnu4v+jfB6RDZcvThZJOKX0XSrEB+SodS
	rnBtwNJ1+yhP0wnbJ63fNvetbEq7LstGN2bVi4KWnRaDh07fyCwFL771gRKRK2g+8JhiML
	AOcqq672WLWZMZ/PRXnpU0HIPF6n7Pz2Yve8n7xnAAAFmP5TY4v+U2OLAAAAB3NzaC1yc2
	EAAAGBANVIuiEGPhXXgoBj8YitZM9Upnrdbosnb4i+KUhs2ZCWpoSMN2X/4nHFHB/NTmrN
	tX4v7Ehil8jBK9T96j5N27NuTFKVakvDHP9rt4K5Iy/uQSYp+Q1aXl/kpc5ccqw0wldD1P
	o29LwhqKNRbonmuf9eSBLMH1sPpRIySNOar2DZF9oUAOO9RWmZkJW2KJErUDrSWg4y628Q
	ozXY08+L5Ucz3aRvI75Uo0+C4j/DPYLTxUT3/zIF18fn2VMfm/yovBXqkz+Em2w4X4LX+y
	quV3ha2TjB2tje/omlBd7U1fxjaMsC3KNbDpKcdqQ8vK5qEAiWsXJrlCJY4mcR5907vGjT
	7QCaFmHz7cD0KchwnDFu9rxjXzZp7uL/o3wekQ2XL04WSTil9F0qxAfkqHUq5wbcDSdfso
	T9MJ2yet3zb3rWxKuy7LRjdm1YuClp0Wg4dO38gsBS++9YESkStoPvCYYjCwDnKquu9li1
	mTGfz0V56VNByDxep+z89mL3vJ+8ZwAAAAMBAAEAAAGBAJELsZDt5uEBu71GuqbBjLI3Fj
	SuTBQUUJSFBhw78kWTPlEb7jzOpRfL/ZFfFPorRUc4ng6oBiM/w2hI+bk/R68hzoPHGw/E
	8/58KcOb1mMtO18R4k6Da3T5UQ0i79VO1+9ysO8s2ojqtv3CTlM39rvFSWyHJrfNzuuuCL
	rnEmfhm4fyXJyERiVHiv1VcQcwlpI6JYZMeLICdYwUFg+qStV+XzgJYRx6AMn875J/W2CS
	VjDOGt3Q/Wr0sGYINBPCR1Ci1yXSaZCHtY0P9yCAh0Elh9htHuOtn6nA3nnhOUyvpVHl0S
	eQDEQI0jDycgKfsf0chHC5Abmc9rHFcRKEbfggWNjIFwl3S+Q3gCTtIAqTrtwQ2pi0giQR
	KOc4sT+g6bEmzVsgkEmuSuFAMDyp65Rb1zTB1GgCVCk594k0bRFX4AuA0/fTeSby4enD/Y
	knjWOEjK87umagvQqfp67Ur5fjQpCszrG9rWHnTTaHbCMzAeWVfplrQ8GLwn+r0GwRYQAA
	AMEA4YO9tRjgcf3L2rgGSPOLtyCBnoI+3LIPc87n8xYZkEp/HHXH3EnPIqfatZYYv8IYwD
	UJGCykixuoUpcwHrqQodNLTmm4Q9yf1slBpG08NjXA9ULYy0TZnzxufWSpxnV8JAgAPAJC
	QXkxqfnWWSNXTa+xhjroBjQoiBjLu1DCCuD26aXdGyZJKmmUIcxCskVHXXvMJw+hbb7/7/
	hrXpe7fbyDVWnMmkMGZyez7/UMLfyo3UwjiC2uzL+dzu2IQT9jAAAAwQD3igueobM0BH/K
	KQeJIxIn/K5TWV8XVIGtJzQXvExacS2iZFNP01opyeWYpzmbOpRUzuGJScwsQWWEu/UAYX
	9nf5aSxjU/XW2OHFEuiL0ha3ccBnGa+encV4CNXdzQdak2rflQr0nBGXou31cj9wRQCXUn
	tHeBwkupUA+ydCJ5UiNJYii3rRc/urAO18TKTan7bf9eIh8OZTXfvcWw3xNd21FE+XWBqW
	yimBiwhFmCW5FPZfwZoVku0xjflUzbiZEAAADBANyS83XvnEMzXLun4T37L7UNz3cFooFr
	OMJpPlPSvJfQdPfKIeHViRyIER5ONRtX8JMWzJPsoVR6M2XtghlAo0cmh7d58UQ717uUiu
	94mzwJWVyqwRU9/zLJE/HRCvJoWCJUXFMoHnrgPPGdoyylx7F3htOlu4/k3fApdSy/6ZMW
	DWPlUrsaXhWdBAoEYrv/WlIVef5XnP+sB8nC70hFMwie4YLQlX6vQ4joI/Qr1QEGhfaAc3
	krNGWJs7lOdXEqdwAAABx0aW1yaWNoYXJkc29uQEdFTC1DSjc0RDJYOU1GAQIDBAU=
	-----END OPENSSH PRIVATE KEY-----
	EOF
  }
}

resource "nftower_compute_environment" "foo" {
	name           = "tf-acceptance-lsf"
	workspace_id   = nftower_workspace.foo.id
	credentials_id = nftower_credentials.foo.id
  
	lsf_platform {
	  work_dir                 = "/nextflow/work"
	  launch_dir               = "/nextflow/launch"
	  user_name                = "nextflow"
	  host_name                = "example.com"
	  head_queue               = "head"
	  compute_queue            = "compute"
	  head_job_options         = "-x something"
	  propagate_head_job_options = true
	  per_job_mem_limit = false
	  per_task_reserve = false
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

func TestFlattenComputeEnvironmentLSFPlatfromMinimal(t *testing.T) {
	ctx := context.Background()
	actual := flattenComputeEnvironmentLSFPlatform(ctx, &client.ComputeEnvLSFPlatformConfig{
		WorkDir:                 "/nextflow/work",
		LaunchDir:               "/nextflow/launch",
		UserName:                "nextflow",
		HostName:                "localhost",
		HeadQueue:               "head",
		ComputeQueue:            "compute",
		HeadJobOptions:          "-x something",
		PropagateHeadJobOptions: true,
		PerJobMemLimit:          false,
		PerTaskReserve:          false,
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
			"work_dir":                   "/nextflow/work",
			"launch_dir":                 "/nextflow/launch",
			"user_name":                  "nextflow",
			"host_name":                  "localhost",
			"head_queue":                 "head",
			"compute_queue":              "compute",
			"head_job_options":           "-x something",
			"propagate_head_job_options": true,
			"per_job_mem_limit":          false,
			"per_task_reserve":           false,
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestFlattenComputeEnvironmentLSFPlatformComplete(t *testing.T) {
	ctx := context.Background()
	actual := flattenComputeEnvironmentLSFPlatform(ctx, &client.ComputeEnvLSFPlatformConfig{
		WorkDir:                 "/nextflow/work",
		LaunchDir:               "/nextflow/launch",
		UserName:                "nextflow",
		HostName:                "localhost",
		HeadQueue:               "head",
		ComputeQueue:            "compute",
		HeadJobOptions:          "-x something",
		PropagateHeadJobOptions: false,
		PerJobMemLimit:          true,
		PerTaskReserve:          true,
		Port:                    8000,
		MaxQueueSize:            100,
		PreRunScript:            "set -x",
		PostRunScript:           "echo done",
		UnitForLimits:           "GB",
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
			"work_dir":                   "/nextflow/work",
			"launch_dir":                 "/nextflow/launch",
			"user_name":                  "nextflow",
			"host_name":                  "localhost",
			"head_queue":                 "head",
			"compute_queue":              "compute",
			"head_job_options":           "-x something",
			"propagate_head_job_options": false,
			"per_job_mem_limit":          true,
			"per_task_reserve":           true,
			"port":                       8000,
			"max_queue_size":             100,
			"pre_run_script":             "set -x",
			"post_run_script":            "echo done",
			"unit_for_limits":            "GB",
		},
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}
