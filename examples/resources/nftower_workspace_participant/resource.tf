resource "nftower_organization_member" "example" {
  email = "myuser@domain.com"
}

resource "nftower_workspace" "example" {
  name        = "foo"
  full_name   = "foo bar baz"
  description = "A foo workspace"
  visibility  = "PRIVATE"
}

resource "nftower_workspace_participant" "example" {
  workspace_id = nftower_workspace.example.id
  member_id    = nftower_organization_member.example.id
  role         = "maintain"
}

resource "nftower_workspace_participant" "example_by_email" {
  workspace_id = nftower_workspace.example.id
  email        = "myuser@domain.com"
  role         = "maintain"

  depends_on = [nftower_organization_member.example]
}
