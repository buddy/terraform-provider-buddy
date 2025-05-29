resource "buddy_target" "example" {
  domain     = "mydomain"
  name       = "prod"
  identifier = "prod"
  type       = "FTP"
  host       = "1.1.1.1"
  port       = "21"
  secure     = true
}
