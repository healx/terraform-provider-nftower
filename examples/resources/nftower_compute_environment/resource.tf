resource "nftower_workspace" "example" {
  name      = "foo"
  full_name = "foo bar baz"
}

resource "nftower_credentials" "aws" {
  name         = "aws-creds"
  workspace_id = nftower_workspace.example.id

  aws {
    access_key = "ABDFRGTEDRFS"
    secret_key = "sdkjdlgkdjflgkdglkdnflsrkgdlvkslgkdn"
  }
}


resource "nftower_compute_environment" "example-awsbatch" {
  name           = "example-awsbatch"
  workspace_id   = nftower_workspace.example.id
  credentials_id = nftower_workspace.aws.id

  aws_batch {
    region        = "eu-west-1"
    compute_queue = "compute"
    head_queue    = "head"
    work_dir      = "s3://my-nf-workdir"
  }
}
