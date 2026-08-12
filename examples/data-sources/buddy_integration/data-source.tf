data "buddy_integration" "amazon" {
  domain         = "mydomain"
  integration_id = "abcd1234"
}

data "buddy_integration" "azure" {
  domain = "mydomain"
  name   = "azure"
}

data "buddy_integration" "shopify" {
  domain       = "mydomain"
  name         = "shopify"
  project_name = "myproject"
}
