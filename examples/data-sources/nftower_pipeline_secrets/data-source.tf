data "nftower_workspace" "foo" {
  name = "foo"
}

data "nftower_pipeline_secrets" "foo" {
  name         = "foo"
  workspace_id = data.nftower_workspace.foo.id
}