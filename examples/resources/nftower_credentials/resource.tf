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

resource "nftower_credentials" "github" {
  name         = "github-creds"
  workspace_id = nftower_workspace.example.id

  github {
    username     = "my-user"
    access_token = "sdkjdlgkdjflgkdglkdnflsrkgdlvkslgkdn" // a personal access token (PAT)
  }
}
