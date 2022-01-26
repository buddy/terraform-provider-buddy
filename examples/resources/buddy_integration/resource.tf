resource "buddy_integration" "aws" {
  domain     = "mydomain"
  name       = "ec2 access"
  type       = "AMAZON"
  scope      = "ADMIN"
  access_key = "abcdefghijkl"
  secret_key = "abcdefghijkl"
}