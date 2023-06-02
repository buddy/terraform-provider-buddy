resource "buddy_integration" "aws" {
  domain     = "mydomain"
  name       = "ec2 access"
  type       = "AMAZON"
  scope      = "ADMIN"
  access_key = "key"
  secret_key = "secret"

  role_assumption {
    arn = "arn1"
  }

  role_assumption {
    arn         = "arn2"
    external_id = "3"
    duration    = 100
  }
}

resource "buddy_integration" "do_private" {
  domain = "mydomain"
  name   = "digital ocean"
  type   = "DIGITAL_OCEAN"
  scope  = "PRIVATE"
  token  = "abc"
}

resource "buddy_integration" "shopify_workspace" {
  domain = "mydomain"
  name   = "shopify"
  type   = "SHOPIFY"
  scope  = "WORKSPACE"
  shop   = "myshop"
  token  = "abc"
}

resource "buddy_integration" "pushover_group" {
  domain     = "mydomain"
  name       = "pushover"
  type       = "PUSHOVER"
  scope      = "GROUP"
  group_id   = 123
  token      = "abc"
  access_key = "key"
}

resource "buddy_integration" "rackspace_project" {
  domain       = "mydomain"
  name         = "rackspace"
  type         = "RACKSPACE"
  scope        = "PROJECT"
  project_name = "myproject"
  username     = "abc"
  token        = "abc"
}

resource "buddy_integration" "cloudflare_admin_project" {
  domain       = "mydomain"
  name         = "cloudflare"
  type         = "CLOUDFLARE"
  scope        = "ADMIN_IN_PROJECT"
  project_name = "myproject"
  token        = "abc"
  api_key      = "key"
  email        = "email@email.com"
}

resource "buddy_integration" "new_relic_group_project" {
  domain       = "mydomain"
  name         = "new_relic"
  type         = "NEW_RELIC"
  scope        = "GROUP_IN_PROJECT"
  group_id     = 123
  project_name = "myproject"
  token        = "abc"
}

resource "buddy_integration" "sentry" {
  domain = "mydomain"
  name   = "sentry"
  type   = "SENTRY"
  scope  = "ADMIN"
  token  = "abc"
}

resource "buddy_integration" "rollbar" {
  domain = "mydomain"
  name   = "rollbar"
  type   = "ROLLBAR"
  scope  = "ADMIN"
  token  = "abc"
}

resource "buddy_integration" "datadog" {
  domain = "mydomain"
  name   = "datadog"
  type   = "DATADOG"
  scope  = "ADMIN"
  token  = "abc"
}

resource "buddy_integration" "do_spaces" {
  domain     = "mydomain"
  name       = "do_spaces"
  type       = "DO_SPACES"
  scope      = "ADMIN"
  access_key = "key"
  secret_key = "secret"
}

resource "buddy_integration" "honeybadger" {
  domain = "mydomain"
  name   = "honeybadger"
  type   = "HONEYBADGER"
  scope  = "ADMIN"
  token  = "abc"
}

resource "buddy_integration" "vultr" {
  domain = "mydomain"
  name   = "vultr"
  type   = "VULTR"
  scope  = "ADMIN"
  token  = "abc"
}

resource "buddy_integration" "sentry_enterprise" {
  domain = "mydomain"
  name   = "sentry_enterprise"
  type   = "SENTRY_ENTERPRISE"
  scope  = "ADMIN"
  token  = "abc"
}

resource "buddy_integration" "loggly" {
  domain = "mydomain"
  name   = "loggly"
  type   = "LOGGLY"
  scope  = "ADMIN"
  token  = "abc"
}

resource "buddy_integration" "firebase" {
  domain = "mydomain"
  name   = "firebase"
  type   = "FIREBASE"
  scope  = "ADMIN"
  token  = "abc"
}

resource "buddy_integration" "upcloud" {
  domain   = "mydomain"
  name     = "upcloud"
  type     = "UPCLOUD"
  scope    = "ADMIN"
  username = "abc"
  password = "pass"
}

resource "buddy_integration" "ghost_inspector" {
  domain = "mydomain"
  name   = "ghost_inspector"
  type   = "GHOST_INSPECTOR"
  scope  = "ADMIN"
  token  = "abc"
}

resource "buddy_integration" "azure_cloud" {
  domain    = "mydomain"
  name      = "azure_cloud"
  type      = "AZURE_CLOUD"
  scope     = "ADMIN"
  app_id    = "id"
  tenant_id = "tenant"
  password  = "pass"
}

resource "buddy_integration" "docker_hub" {
  domain   = "mydomain"
  name     = "docker_hub"
  type     = "DOCKER_HUB"
  scope    = "ADMIN"
  username = "abc"
  password = "pass"
}

resource "buddy_integration" "github" {
  domain = "mydomain"
  name   = "github"
  type   = "GIT_HUB"
  scope  = "ADMIN"
  token  = "github token"
}

resource "buddy_integration" "gitlab" {
  domain = "mydomain"
  name   = "gitlab"
  type   = "GIT_LAB"
  scope  = "ADMIN"
  token  = "gitlab token"
}

resource "buddy_integration" "google_service_account" {
  domain  = "mydomain"
  name    = "google_service_account"
  type    = "GOOGLE_SERVICE_ACCOUNT"
  scope   = "ADMIN"
  api_key = "key"
}