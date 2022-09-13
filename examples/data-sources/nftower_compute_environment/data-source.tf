data "nftower_workspace" "foo" {
  name = "foo"
}

data "nftower_compute_environment" "foo" {
  name         = "foo"
  workspace_id = data.nftower_workspace.foo.id
}