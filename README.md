# Nextflow Tower terraform provider

This provider manages configuration for Sequera's [Nextflow Tower](https://tower.nf/) product. It has been tested using the SaaS version.

For detailed documentation see the [registry page](https://registry.terraform.io/providers/healx/nftower).

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
-	[Go](https://golang.org/doc/install) >= 1.18

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 
```sh
$ go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

```hcl
provider "nftower" {
  api_key      = "..." // can also be set using NFTOWER_API_KEY
  organization = "my-org"
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`. You will need to specify an api key using `NFTOWER_API_KEY` and an organization using `NFTOWER_ORGANIZATION`

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

## Making a release

If you wish to make a release, you must tag a commit with the version you wish to release and then push the tag to Github. A Github action will trigger to create the release and then the terraform registry will detect it and update.

```
git tag v1.2.3
git push origin v1.2.3
```