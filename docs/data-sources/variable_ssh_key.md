---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "buddy_variable_ssh_key Data Source - terraform-provider-buddy"
subcategory: ""
description: |-
  Get variables of SSH key type by key or variable ID
  Token scope required: WORKSPACE, VARIABLE_INFO
---

# buddy_variable_ssh_key (Data Source)

Get variables of SSH key type by key or variable ID

Token scope required: `WORKSPACE`, `VARIABLE_INFO`

## Example Usage

```terraform
data "buddy_variable" "by_id" {
  variable_id = 123456
}

data "buddy_variable" "by_key" {
  key = "MY_KEY"
}

data "buddy_variable" "by_key_in_project" {
  key          = "MY_KEY"
  project_name = "myproject"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `domain` (String) The workspace's URL handle

### Optional

- `action_id` (Number) The variable's action ID
- `key` (String) The variable's name
- `pipeline_id` (Number) The variable's pipeline ID
- `project_name` (String) The variable's project name
- `variable_id` (Number) The variable's ID

### Read-Only

- `checksum` (String) The variable's checksum
- `description` (String) The variable's description
- `encrypted` (Boolean) Is the variable's value encrypted, always true for buddy_variable_ssh_key
- `file_chmod` (String) The variable's file permission in an action's container
- `file_path` (String) The variable's path in the action's container
- `file_place` (String) Should the variable's be copied to an action's container in **file_path** (`CONTAINER`, `NONE`)
- `id` (String) The Terraform resource identifier for this item
- `key_fingerprint` (String) The variable's fingerprint
- `public_value` (String) The variable's public key
- `settable` (Boolean) Is the variable's value changeable, always false for buddy_variable_ssh_key
- `value` (String, Sensitive) The variable's value, always encrypted for buddy_variable_ssh_key

