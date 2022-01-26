data "buddy_integration" "amazon" {
  domain         = "mydomain"
  integration_id = "abcd1234"
}

data "buddy_integration" "azure" {
  domain = "mydomain"
  name   = "azure"
}