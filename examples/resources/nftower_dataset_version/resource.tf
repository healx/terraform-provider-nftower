resource "nftower_workspace" "example" {
  name      = "foo"
  full_name = "foo bar baz"
}

resource "nftower_dataset" "example" {
  name         = "my-dataset"
  workspace_id = nftower_workspace.example.id
}

resource "nftower_dataset_version" "example_csv" {
  dataset_id   = nftower_dataset.example.id
  workspace_id = nftower_workspace.example.id

  file_name  = "foo.csv"
  contents   = <<EOF
one,two,three,four
1,2,3,4
EOF
  has_header = true
}
