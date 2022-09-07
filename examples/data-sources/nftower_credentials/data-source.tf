data "nftower_workspace" "foo" {
  name = "foo"
}

data "nftower_credentials" "foo" {
  name         = "foo"
  workspace_id = data.nftower_workspace.foo.id
}