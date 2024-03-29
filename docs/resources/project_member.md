---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "buddy_project_member Resource - terraform-provider-buddy"
subcategory: ""
description: |-
  Manage a member's permission (role) in a project
  Workspace administrator rights are required
  Token scope required: WORKSPACE
---

# buddy_project_member (Resource)

Manage a member's permission (role) in a project

Workspace administrator rights are required

Token scope required: `WORKSPACE`

## Example Usage

```terraform
resource "buddy_project_member" "john_in_test" {
  domain        = "mydomain"
  project_name  = "test"
  member_id     = 1234
  permission_id = 5678
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `domain` (String) The workspace's URL handle
- `member_id` (Number) The member's ID
- `permission_id` (Number) The permission's ID
- `project_name` (String) The project's name

### Read-Only

- `admin` (Boolean) Is the member a workspace administrator
- `avatar_url` (String) The member's avatar URL
- `email` (String) The member's email
- `html_url` (String) The member's URL
- `id` (String) The Terraform resource identifier for this item
- `name` (String) The member's name
- `permission` (Attributes Set) The member's permission in the project (see [below for nested schema](#nestedatt--permission))
- `workspace_owner` (Boolean) Is the member the workspace owner

<a id="nestedatt--permission"></a>
### Nested Schema for `permission`

Read-Only:

- `html_url` (String)
- `name` (String)
- `permission_id` (Number)
- `pipeline_access_level` (String)
- `project_team_access_level` (String)
- `repository_access_level` (String)
- `sandbox_access_level` (String)
- `type` (String)

## Import

Import is supported using the following syntax:

```shell
# import using domain(mydomain), project name(test), member_id(1234)
terraform import buddy_project_member.john_in_test mydomain:test:1234
```
