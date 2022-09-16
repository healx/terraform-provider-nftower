data "nftower_workspace" "example" {
  name = "my-workspace"
}

data "nftower_workspace_participant" "example" {
  workspace_id = data.nftower_workspace.example.id
  email        = "someone@mydomain.com"
}
