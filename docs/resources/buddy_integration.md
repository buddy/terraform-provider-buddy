---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "buddy_integration Resource - buddy-terraform"
subcategory: ""
description: |-
  Create and manage an integration
  Token scopes required: INTEGRATION_ADD, INTEGRATION_MANAGE, INTEGRATION_INFO
---

# buddy_integration (Resource)

Create and manage an integration

Token scopes required: `INTEGRATION_ADD`, `INTEGRATION_MANAGE`, `INTEGRATION_INFO`

## Example Usage

```terraform
resource "buddy_integration" "aws" {
  domain     = "mydomain"
  name       = "ec2 access"
  type       = "AMAZON"
  scope      = "ADMIN"
  access_key = "abcdefghijkl"
  secret_key = "abcdefghijkl"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **domain** (String) The workspace's URL handle
- **name** (String) The integration's name
- **scope** (String) The integration's scope. Allowed:

`PRIVATE` - only creator of the integration can use it

`WORKSPACE` - all workspace members can use the integration

`ADMIN` - only workspace administrators can use the integration

`GROUP` - only group members can use the integration

`PROJECT` - only project members can use the integration

`ADMIN_IN_PROJECT` - only workspace administrators in specified project can use the integration

`GROUP_IN_PROJECT` - only group members in specified project can use the integration
- **type** (String) The integration's type. Allowed: `DIGITAL_OCEAN`, `AMAZON`, `SHOPIFY`, `PUSHOVER`, `RACKSPACE`, `CLOUDFLARE`, `NEW_RELIC`, `SENTRY`, `ROLLBAR`, `DATADOG`, `DO_SPACES`, `HONEYBADGER`, `VULTR`, `SENTRY_ENTERPRISE`, `LOGGLY`, `FIREBASE`, `UPCLOUD`, `GHOST_INSPECTOR`, `AZURE_CLOUD`, `DOCKER_HUB`, `GOOGLE_SERVICE_ACCOUNT`

### Optional

- **access_key** (String, Sensitive) The integration's access key. Provide for: `DO_SPACES`, `AMAZON`, `PUSHOVER`
- **api_key** (String, Sensitive) The integration's API key. Provide for: `CLOUDFLARE`, `GOOGLE_SERVICE_ACCOUNT`
- **app_id** (String) The integration's application's ID. Provide for: `AZURE_CLOUD`
- **email** (String) The integration's email. Provide for: `CLOUDFLARE`
- **group_id** (Number) The group's ID. Provide along with scopes: `GROUP`, `GROUP_IN_PROJECT`
- **password** (String, Sensitive) The integration's password. Provide for: `AZURE_CLOUD`, `UPCLOUD`, `DOCKER_HUB`
- **project_name** (String) The project's name. Provide along with scopes: `PROJECT`, `ADMIN_IN_PROJECT`, `GROUP_IN_PROJECT`
- **role_assumption** (Block List) The integration's AWS role to assume. Provide for: `AMAZON` (see [below for nested schema](#nestedblock--role_assumption))
- **secret_key** (String, Sensitive) The integration's secret key. Provide for: `DO_SPACES`, `AMAZON`
- **shop** (String) The integration's shop. Provide for: `SHOPIFY`
- **tenant_id** (String) The integration's tenant's ID. Provide for: `AZURE_CLOUD`
- **token** (String, Sensitive) The integration's token. Provide for: `DIGITAL_OCEAN`, `SHOPIFY`, `RACKSPACE`, `CLOUDFLARE`, `NEW_RELIC`, `SENTRY`, `ROLLBAR`, `DATADOG`, `HONEYBADGER`, `VULTR`, `SENTRY_ENTERPRISE`, `LOGGLY`, `FIREBASE`, `GHOST_INSPECTOR`, `PUSHOVER`
- **username** (String) The integration's username. Provide for: `UPCLOUD`, `RACKSPACE`, `DOCKER_HUB`

### Read-Only

- **html_url** (String) The integration's URL
- **id** (String) The Terraform resource identifier for this item
- **integration_id** (String) The integration's ID

<a id="nestedblock--role_assumption"></a>
### Nested Schema for `role_assumption`

Required:

- **arn** (String)

Optional:

- **duration** (Number)
- **external_id** (String)

## Import

Import is supported using the following syntax:

```shell
# import using domain(mydomain), integration_id(abc123)
terraform import buddy_integration.aws mydomain:abc123
```