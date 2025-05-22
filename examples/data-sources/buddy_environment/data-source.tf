data "buddy_environment" "dev" {
  domain         = "mydomain"
  environment_id = 1234
}

data "buddy_environment" "stage" {
  domain = "mydomain"
  name   = "stage"
}