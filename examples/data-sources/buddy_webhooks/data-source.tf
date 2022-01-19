data "buddy_webhooks" "all" {
  domain = "mydomain"
}

data "buddy_webhooks" "filter" {
  domain           = "mydomain"
  target_url_regex = "/my$"
}