data "buddy_webhook" "by_id" {
  webhook_id = 123456
}

data "buddy_webhook" "by_target_url" {
  target_url = "https://mydomain.com"
}