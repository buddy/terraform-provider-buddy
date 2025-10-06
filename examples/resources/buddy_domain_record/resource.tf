resource "buddy_domain_record" "a" {
  workspace_domain = "myworkspace"
  domain_id        = "abcdef1234"
  domain           = "a.test.com"
  type             = "A"
  ttl              = 60
  value            = ["1.1.1.1", "2.2.2.2"]
}

resource "buddy_domain_record" "geo" {
  workspace_domain = "myworkspace"
  domain_id        = "abcdef1234"
  domain           = "geo.test.com"
  type             = "TXT"
  ttl              = 60
  value            = ["Fallback"]
  continent = {
    "Asia"   = ["Asia"]
    "Europe" = ["Europe"]
  }
  country = {
    "GB" = ["United Kingdom"]
    "TR" = ["Turkey"]
  }
}