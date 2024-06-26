---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "nftower_pipeline_secrets Resource - terraform-provider-nftower"
subcategory: ""
description: |-
  A pipeline-secret for use by Tower.
---

# nftower_pipeline_secrets (Resource)

A pipeline-secret for use by Tower.

## Example Usage

```terraform
resource "nftower_workspace" "example" {
  name      = "foo"
  full_name = "foo bar baz"
}

resource "nftower_pipeline_secrets" "secret" {
  name         = "PIPELINE_SECRET"
  workspace_id = nftower_workspace.example.id
  value        = "some secret value"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the pipeline-secret.
- `value` (String) The value of the pipeline-secret.
- `workspace_id` (String) The id of the workspace in which to create the pipeline-secret.

### Read-Only

- `date_created` (String) The datetime the workspace was created.
- `id` (String) The ID of this resource.
- `last_updated` (String) The last updated datetime of the workspace.
