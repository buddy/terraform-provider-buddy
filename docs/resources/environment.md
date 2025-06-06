---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "buddy_environment Resource - terraform-provider-buddy"
subcategory: ""
description: |-
  Create and manage an environment
  Token scopes required: WORKSPACE, ENVIRONMENT_MANAGE, ENVIRONMENT_INFO
---

# buddy_environment (Resource)

Create and manage an environment

Token scopes required: `WORKSPACE`, `ENVIRONMENT_MANAGE`, `ENVIRONMENT_INFO`

## Example Usage

```terraform
resource "buddy_environment" "dev" {
  domain                = "mydomain"
  project_name          = "myproject"
  name                  = "dev"
  identifier            = "dev"
  public_url            = "https://dev.com"
  all_pipelines_allowed = true
  tags                  = ["frontend", "backend"]
  var {
    key   = "ENV"
    value = "DEV"
  }
  permissions {
    others = "MANAGE"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `domain` (String) The workspace's URL handle
- `identifier` (String) The environment's identifier
- `name` (String) The environment's name
- `project_name` (String) The project's name

### Optional

- `all_pipelines_allowed` (Boolean) Defines whether or not environment can be used in all pipelines
- `permissions` (Block Set) The environment's permissions (see [below for nested schema](#nestedblock--permissions))
- `public_url` (String) The environment's public URL
- `tags` (Set of String) The environment's list of tags
- `var` (Block Set) The environment's variables (see [below for nested schema](#nestedblock--var))

### Read-Only

- `environment_id` (String) The environment's ID
- `html_url` (String) The environment's URL
- `id` (String) The Terraform resource identifier for this item
- `project` (Attributes Set) The environment's project (see [below for nested schema](#nestedatt--project))

<a id="nestedblock--permissions"></a>
### Nested Schema for `permissions`

Optional:

- `group` (Block Set) (see [below for nested schema](#nestedblock--permissions--group))
- `others` (String)
- `user` (Block Set) (see [below for nested schema](#nestedblock--permissions--user))

<a id="nestedblock--permissions--group"></a>
### Nested Schema for `permissions.group`

Required:

- `access_level` (String)

Read-Only:

- `id` (Number) The ID of this resource.


<a id="nestedblock--permissions--user"></a>
### Nested Schema for `permissions.user`

Required:

- `access_level` (String)

Read-Only:

- `id` (Number) The ID of this resource.



<a id="nestedblock--var"></a>
### Nested Schema for `var`

Required:

- `key` (String)
- `value` (String, Sensitive)

Optional:

- `description` (String)
- `encrypted` (Boolean)
- `settable` (Boolean)


<a id="nestedatt--project"></a>
### Nested Schema for `project`

Read-Only:

- `display_name` (String)
- `html_url` (String)
- `name` (String)
- `status` (String)

## Import

Import is supported using the following syntax:

```shell
# import using domain(mydomain), project name (myproject) and environment id (123456)
terraform import buddy_environment.test mydomain:myproject:123456
```
