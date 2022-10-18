data "nftower_workspace" "example" {
  name = "foo"
}

data "nftower_compute_environment" "example-awsbatch" {
  name         = "example-awsbatch"
  workspace_id = nftower_workspace.example.id
}

resource "nftower_pipeline" "example-tower" {
  name                   = "example-tower"
  workspace_id           = data.nftower_workspace.example.id
  compute_environment_id = data.nftower_compute_environment.example-awsbatch.id
  pipeline               = "https://github.com/nextflow-io/hello"
  work_dir               = data.nftower_compute_environment.example-awsbatch.aws_batch.0.work_dir
}