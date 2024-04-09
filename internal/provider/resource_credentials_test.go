package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/healx/terraform-provider-nftower/internal/template"
)

func TestAccResourceCredentialsAWS(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_credentials",
				Config:       template.ParseRandName(testAccResourceCredentialsAWS),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "name", "tf-acceptance-credentials-aws"),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "description", "tf acceptance testing aws credentials"),
					resource.TestMatchResourceAttr(
						"nftower_credentials.foo", "date_created", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestMatchResourceAttr(
						"nftower_credentials.foo", "last_updated", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "aws.0.access_key", "foo"),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "aws.0.secret_key", "bar"),
				),
			},
		},
	})
}

const testAccResourceCredentialsAWS = `
resource "nftower_workspace" "foo" {
  name        = "tf-acceptance-{{.randName}}"
  full_name   = "tf acceptance testing credentials"
	
  description = "Created by the nftower terraform provider acceptance tests. Will be deleted shortly"
  visibility  = "PRIVATE"
}

resource "nftower_credentials" "foo" {
  name        = "tf-acceptance-credentials-aws"
  description = "tf acceptance testing aws credentials"
  workspace_id = nftower_workspace.foo.id

  aws {
	access_key      = "foo"
	secret_key      = "bar"
	assume_role_arn = "baz"
  }
}
`

// TestAccResourceCredentialsContainerRegistry cannot be run withouth real credentials for accessing the container registry.
// Seqera Platform test registry connection during the creation of the credentials.
// For AWS ECR would be username: AWS_ACCESS_KEY_ID), password: AWS_SECRET_ACCESS_KEY, registry_server: AWS_ACCOUNT_ID.dkr.ecr.AWS_REGION.amazonaws.com
func TestAccResourceCredentialsContainerRegistry(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_credentials",
				Config:       template.ParseRandName(testAccResourceCredentialsContainerRegistry),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "name", "tf-acceptance-credentials-container-registry"),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "description", "tf acceptance testing container registry credentials"),
					resource.TestMatchResourceAttr(
						"nftower_credentials.foo", "date_created", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestMatchResourceAttr(
						"nftower_credentials.foo", "last_updated", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "container_registry.0.username", "<<AWS_ACCESS_KEY_ID>>"),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "container_registry.0.password", "<<AWS_SECRET_ACCESS_KEY>>"),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "container_registry.0.registry_server", "<<AWS_ACCOUNT_ID.dkr.ecr.AWS_REGION.amazonaws.com>>"),
				),
			},
		},
	})
}

const testAccResourceCredentialsContainerRegistry = `
resource "nftower_workspace" "foo" {
  name        = "tf-acceptance-{{.randName}}"
  full_name   = "tf acceptance testing credentials"

  description = "Created by the nftower terraform provider acceptance tests. Will be deleted shortly"
  visibility  = "PRIVATE"
}

resource "nftower_credentials" "foo" {
  name        = "tf-acceptance-credentials-container-registry"
  description = "tf acceptance testing container registry credentials"
  workspace_id = nftower_workspace.foo.id

  container_registry {
	username      = "<<AWS_ACCESS_KEY_ID>>"
	password      = "<<AWS_SECRET_ACCESS_KEY"
	registry_server = "<<AWS_ACCOUNT_ID.dkr.ecr.AWS_REGION.amazonaws.com>>"
  }
}
`

func TestAccResourceCredentialsGithub(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_credentials",
				Config:       template.ParseRandName(testAccResourceCredentialsGithub),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "name", "tf-acceptance-credentials-github"),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "description", "tf acceptance testing github credentials"),
					resource.TestMatchResourceAttr(
						"nftower_credentials.foo", "date_created", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestMatchResourceAttr(
						"nftower_credentials.foo", "last_updated", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "github.0.username", "foo"),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "github.0.access_token", "bar"),
				),
			},
		},
	})
}

const testAccResourceCredentialsGithub = `
resource "nftower_workspace" "foo" {
  name        = "tf-acceptance-{{.randName}}"
  full_name   = "tf acceptance testing credentials"
	
  description = "Created by the nftower terraform provider acceptance tests. Will be deleted shortly"
  visibility  = "PRIVATE"
}

resource "nftower_credentials" "foo" {
  name        = "tf-acceptance-credentials-github"
  description = "tf acceptance testing github credentials"
  workspace_id = nftower_workspace.foo.id

  github {
	username     = "foo"
	access_token = "bar"
  }
}
`

func TestAccResourceCredentialsGitlab(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_credentials",
				Config:       template.ParseRandName(testAccResourceCredentialsGitlab),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "name", "tf-acceptance-credentials-gitlab"),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "description", "tf acceptance testing gitlab credentials"),
					resource.TestMatchResourceAttr(
						"nftower_credentials.foo", "date_created", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestMatchResourceAttr(
						"nftower_credentials.foo", "last_updated", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "gitlab.0.username", "foo"),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "gitlab.0.password", "bar"),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "gitlab.0.token", "baz"),
				),
			},
		},
	})
}

const testAccResourceCredentialsGitlab = `
resource "nftower_workspace" "foo" {
  name        = "tf-acceptance-{{.randName}}"
  full_name   = "tf acceptance testing credentials"
	
  description = "Created by the nftower terraform provider acceptance tests. Will be deleted shortly"
  visibility  = "PRIVATE"
}

resource "nftower_credentials" "foo" {
  name        = "tf-acceptance-credentials-gitlab"
  description = "tf acceptance testing gitlab credentials"
  workspace_id = nftower_workspace.foo.id

  gitlab {
	username     = "foo"
	password 	 = "bar"
	token 		 = "baz"
  }
}
`

func TestAccResourceCredentialsSSH(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_credentials",
				Config:       template.ParseRandName(testAccResourceCredentialsSSH),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "name", "tf-acceptance-credentials-ssh"),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "description", "tf acceptance testing ssh credentials"),
					resource.TestMatchResourceAttr(
						"nftower_credentials.foo", "date_created", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestMatchResourceAttr(
						"nftower_credentials.foo", "last_updated", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestMatchResourceAttr(
						"nftower_credentials.foo", "ssh.0.private_key", regexp.MustCompile("BEGIN OPENSSH PRIVATE KEY")),
				),
			},
		},
	})
}

const testAccResourceCredentialsSSH = `
resource "nftower_workspace" "foo" {
  name        = "tf-acceptance-{{.randName}}"
  full_name   = "tf acceptance testing credentials"
	
  description = "Created by the nftower terraform provider acceptance tests. Will be deleted shortly"
  visibility  = "PRIVATE"
}

resource "nftower_credentials" "foo" {
  name        = "tf-acceptance-credentials-ssh"
  description = "tf acceptance testing ssh credentials"
  workspace_id = nftower_workspace.foo.id

  ssh {
	private_key	 = <<EOF
	-----BEGIN OPENSSH PRIVATE KEY-----
	b3BlbnNzaC1rZXktdjEAAAAABG5vbmUAAAAEbm9uZQAAAAAAAAABAAABlwAAAAdzc2gtcn
	NhAAAAAwEAAQAAAYEA1Ui6IQY+FdeCgGPxiK1kz1Smet1uiydviL4pSGzZkJamhIw3Zf/i
	ccUcH81Oas21fi/sSGKXyMEr1P3qPk3bs25MUpVqS8Mc/2u3grkjL+5BJin5DVpeX+Slzl
	xyrDTCV0PU+jb0vCGoo1Fuiea5/15IEswfWw+lEjJI05qvYNkX2hQA471FaZmQlbYokStQ
	OtJaDjLrbxCjNdjTz4vlRzPdpG8jvlSjT4LiP8M9gtPFRPf/MgXXx+fZUx+b/Ki8FeqTP4
	SbbDhfgtf7Kq5XeFrZOMHa2N7+iaUF3tTV/GNoywLco1sOkpx2pDy8rmoQCJaxcmuUIlji
	ZxHn3Tu8aNPtAJoWYfPtwPQpyHCcMW72vGNfNmnu4v+jfB6RDZcvThZJOKX0XSrEB+SodS
	rnBtwNJ1+yhP0wnbJ63fNvetbEq7LstGN2bVi4KWnRaDh07fyCwFL771gRKRK2g+8JhiML
	AOcqq672WLWZMZ/PRXnpU0HIPF6n7Pz2Yve8n7xnAAAFmP5TY4v+U2OLAAAAB3NzaC1yc2
	EAAAGBANVIuiEGPhXXgoBj8YitZM9Upnrdbosnb4i+KUhs2ZCWpoSMN2X/4nHFHB/NTmrN
	tX4v7Ehil8jBK9T96j5N27NuTFKVakvDHP9rt4K5Iy/uQSYp+Q1aXl/kpc5ccqw0wldD1P
	o29LwhqKNRbonmuf9eSBLMH1sPpRIySNOar2DZF9oUAOO9RWmZkJW2KJErUDrSWg4y628Q
	ozXY08+L5Ucz3aRvI75Uo0+C4j/DPYLTxUT3/zIF18fn2VMfm/yovBXqkz+Em2w4X4LX+y
	quV3ha2TjB2tje/omlBd7U1fxjaMsC3KNbDpKcdqQ8vK5qEAiWsXJrlCJY4mcR5907vGjT
	7QCaFmHz7cD0KchwnDFu9rxjXzZp7uL/o3wekQ2XL04WSTil9F0qxAfkqHUq5wbcDSdfso
	T9MJ2yet3zb3rWxKuy7LRjdm1YuClp0Wg4dO38gsBS++9YESkStoPvCYYjCwDnKquu9li1
	mTGfz0V56VNByDxep+z89mL3vJ+8ZwAAAAMBAAEAAAGBAJELsZDt5uEBu71GuqbBjLI3Fj
	SuTBQUUJSFBhw78kWTPlEb7jzOpRfL/ZFfFPorRUc4ng6oBiM/w2hI+bk/R68hzoPHGw/E
	8/58KcOb1mMtO18R4k6Da3T5UQ0i79VO1+9ysO8s2ojqtv3CTlM39rvFSWyHJrfNzuuuCL
	rnEmfhm4fyXJyERiVHiv1VcQcwlpI6JYZMeLICdYwUFg+qStV+XzgJYRx6AMn875J/W2CS
	VjDOGt3Q/Wr0sGYINBPCR1Ci1yXSaZCHtY0P9yCAh0Elh9htHuOtn6nA3nnhOUyvpVHl0S
	eQDEQI0jDycgKfsf0chHC5Abmc9rHFcRKEbfggWNjIFwl3S+Q3gCTtIAqTrtwQ2pi0giQR
	KOc4sT+g6bEmzVsgkEmuSuFAMDyp65Rb1zTB1GgCVCk594k0bRFX4AuA0/fTeSby4enD/Y
	knjWOEjK87umagvQqfp67Ur5fjQpCszrG9rWHnTTaHbCMzAeWVfplrQ8GLwn+r0GwRYQAA
	AMEA4YO9tRjgcf3L2rgGSPOLtyCBnoI+3LIPc87n8xYZkEp/HHXH3EnPIqfatZYYv8IYwD
	UJGCykixuoUpcwHrqQodNLTmm4Q9yf1slBpG08NjXA9ULYy0TZnzxufWSpxnV8JAgAPAJC
	QXkxqfnWWSNXTa+xhjroBjQoiBjLu1DCCuD26aXdGyZJKmmUIcxCskVHXXvMJw+hbb7/7/
	hrXpe7fbyDVWnMmkMGZyez7/UMLfyo3UwjiC2uzL+dzu2IQT9jAAAAwQD3igueobM0BH/K
	KQeJIxIn/K5TWV8XVIGtJzQXvExacS2iZFNP01opyeWYpzmbOpRUzuGJScwsQWWEu/UAYX
	9nf5aSxjU/XW2OHFEuiL0ha3ccBnGa+encV4CNXdzQdak2rflQr0nBGXou31cj9wRQCXUn
	tHeBwkupUA+ydCJ5UiNJYii3rRc/urAO18TKTan7bf9eIh8OZTXfvcWw3xNd21FE+XWBqW
	yimBiwhFmCW5FPZfwZoVku0xjflUzbiZEAAADBANyS83XvnEMzXLun4T37L7UNz3cFooFr
	OMJpPlPSvJfQdPfKIeHViRyIER5ONRtX8JMWzJPsoVR6M2XtghlAo0cmh7d58UQ717uUiu
	94mzwJWVyqwRU9/zLJE/HRCvJoWCJUXFMoHnrgPPGdoyylx7F3htOlu4/k3fApdSy/6ZMW
	DWPlUrsaXhWdBAoEYrv/WlIVef5XnP+sB8nC70hFMwie4YLQlX6vQ4joI/Qr1QEGhfaAc3
	krNGWJs7lOdXEqdwAAABx0aW1yaWNoYXJkc29uQEdFTC1DSjc0RDJYOU1GAQIDBAU=
	-----END OPENSSH PRIVATE KEY-----
	EOF
  }
}
`

func TestAccResourceCredentials_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				ResourceName: "nftower_credentials",
				Config:       template.ParseRandName(testAccResourceCredentials_basic),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "name", "tf-acceptance-credentials-github"),
					resource.TestMatchResourceAttr(
						"nftower_credentials.foo", "date_created", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestMatchResourceAttr(
						"nftower_credentials.foo", "last_updated", regexp.MustCompile("^[0-9-:TZ]+")),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "github.0.username", "foo"),
					resource.TestCheckResourceAttr(
						"nftower_credentials.foo", "github.0.access_token", "bar"),
				),
			},
		},
	})
}

const testAccResourceCredentials_basic = `
resource "nftower_workspace" "foo" {
  name        = "tf-acceptance-{{.randName}}"
  full_name   = "tf acceptance testing credentials"

  description = "Created by the nftower terraform provider acceptance tests. Will be deleted shortly"
  visibility  = "PRIVATE"
}

resource "nftower_credentials" "foo" {
  name        = "tf-acceptance-credentials-github"
  workspace_id = nftower_workspace.foo.id

  github {
	username     = "foo"
	access_token = "bar"
  }
}
`
