data "buddy_integrations" "all" {
  domain = "mydomain"
}

data "buddy_integrations" "amazon" {
  domain = "mydomain"
  type   = "AMAZON"
}




