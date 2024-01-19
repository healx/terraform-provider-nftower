resource "nftower_workspace" "example" {
  name      = "foo"
  full_name = "foo bar baz"
}

resource "nftower_pipeline_secrets" "secret" {
  name         = "PIPELINE_SECRET"
  workspace_id = nftower_workspace.example.id
  value        = "some secret value"
}
