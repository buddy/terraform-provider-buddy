---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "buddy_profile_public_key Resource - terraform-provider-buddy"
subcategory: ""
description: |-
  Create and manage a user's public key
  Token scope required: USER_KEY
---

# buddy_profile_public_key (Resource)

Create and manage a user's public key

Token scope required: `USER_KEY`

## Example Usage

```terraform
resource "buddy_profile_public_key" "my_key" {
  content = "ssh-rsa ..."
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `content` (String) The public key's content (starts with ssh-rsa)

### Optional

- `title` (String) The public key's title

### Read-Only

- `html_url` (String) The public key's URL
- `id` (String) The Terraform resource identifier for this item

## Import

Import is supported using the following syntax:

```shell
# import using public key id
terraform import buddy_profile_public_key.my_key 123456
```
