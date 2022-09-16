resource "nftower_organization_member" "example" {
  email = "myuser@domain.com"
}

resource "nftower_organization_member" "owner" {
  email = "myadminuser@domain.com"
  role  = "owner"
}
