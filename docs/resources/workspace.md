---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "buddy_workspace Resource - terraform-provider-buddy"
subcategory: ""
description: |-
  Create and manage a workspace
  Invite-only token is required. Contact support@buddy.works for more details
  Token scope required: WORKSPACE
---

# buddy_workspace (Resource)

Create and manage a workspace

Invite-only token is required. Contact support@buddy.works for more details

Token scope required: `WORKSPACE`

## Example Usage

```terraform
resource "buddy_workspace" "ws" {
  domain          = "mydomain"
  name            = "Myname"
  encryption_salt = "mysalt"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `domain` (String) The workspace's URL handle

### Optional

- `encryption_salt` (String) The workspace's salt to encrypt secrets in YAML & API
- `name` (String) The workspace's name

### Read-Only

- `create_date` (String) The workspace's create date
- `frozen` (Boolean) Is the workspace frozen
- `html_url` (String) The workspace's URL
- `id` (String) The Terraform resource identifier for this item
- `owner_id` (Number) The workspace's owner ID
- `workspace_id` (Number) The workspace's ID

## Import

Import is supported using the following syntax:

```shell
# import using domain(mydomain)
terraform import buddy_workspace.ws mydomain
```
