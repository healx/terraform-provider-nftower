resource "nftower_workspace" "example" {
  name      = "foo"
  full_name = "foo bar baz"
}

resource "nftower_dataset" "example" {
  name         = "my-dataset"
  workspace_id = nftower_workspace.example.id
}
