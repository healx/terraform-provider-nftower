data "nftower_workspace" "foo" {
  name = "foo"
}

data "nftower_pipeline" "hello" {
  name         = "hello"
  workspace_id = data.nftower_workspace.foo.id
}