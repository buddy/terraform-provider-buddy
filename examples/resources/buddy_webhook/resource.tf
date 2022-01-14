resource "buddy_webhook" "ws" {
  domain     = "mydomain"
  events     = ["PUSH"]
  projects   = ["myproject"]
  target_url = "https://mydomain.com"
}