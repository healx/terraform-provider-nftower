resource "nftower_token" "example" {
  name = "example"
}

output {
  value = nftower_token.example.token
}
