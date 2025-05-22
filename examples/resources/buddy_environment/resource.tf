resource "buddy_environment" "dev" {
  domain                = "mydomain"
  project_name          = "myproject"
  name                  = "dev"
  identifier            = "dev"
  type                  = "dev"
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