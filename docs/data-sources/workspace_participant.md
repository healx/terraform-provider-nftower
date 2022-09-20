---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "nftower_workspace_participant Data Source - terraform-provider-nftower"
subcategory: ""
description: |-
  A member who has been granted access to a workspace.
---

# nftower_workspace_participant (Data Source)

A member who has been granted access to a workspace.



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `email` (String) The email of the member.
- `workspace_id` (String) The id of the workspace.

### Read-Only

- `first_name` (String) The first name of the member.
- `id` (String) The ID of this resource.
- `last_name` (String) The last name of the member.
- `member_id` (String) The id of the member in the organization.
- `role` (String) The role of the participant.

