---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "buddy_permission Resource - terraform-provider-buddy"
subcategory: ""
description: |-
  Create and manage a workspace permission (role)
  Workspace administrator rights are required
  Token scope required: WORKSPACE
---

# buddy_permission (Resource)

Create and manage a workspace permission (role)

Workspace administrator rights are required

Token scope required: `WORKSPACE`

## Example Usage

```terraform
resource "buddy_permission" "testers" {
  domain                  = "mydomain"
  name                    = "testers"
  pipeline_access_level   = "RUN_ONLY"
  repository_access_level = "READ_ONLY"
  sandbox_access_level    = "DENIED"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `domain` (String) The workspace's URL handle
- `name` (String) The permission's name
- `pipeline_access_level` (String) The permission's access level to pipelines. Allowed: `DENIED`, `READ_ONLY`, `RUN_ONLY`, `READ_WRITE`
- `repository_access_level` (String) The permission's access level to repository. Allowed: `READ_ONLY`, `READ_WRITE`, `MANAGE`
- `sandbox_access_level` (String) The permission's access level to sandboxes. Allowed: `DENIED`, `READ_ONLY`, `READ_WRITE`

### Optional

- `description` (String) The permission's description
- `project_team_access_level` (String) The permission's access level to team. Allowed: `READ_ONLY`, `MANAGE`

### Read-Only

- `html_url` (String) The permission's URL
- `id` (String) The Terraform resource identifier for this item
- `permission_id` (Number) The permission's ID
- `type` (String) The permission's type

## Import

Import is supported using the following syntax:

```shell
# import using domain(mydomain), permission_id(1234)
terraform import buddy_permission.testers mydomain:1234
```
