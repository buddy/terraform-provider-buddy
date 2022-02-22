resource "buddy_variable" "mysecret" {
  domain    = "mydomain"
  key       = "ACCESS_KEY"
  value     = "abcdefgh"
  encrypted = true
}

resource "buddy_variable" "settable_in_project" {
  domain       = "mydomain"
  project_name = "myproject"
  key          = "KEY"
  value        = "init"
  settable     = true
}

resource "buddy_variable" "in_pipeline" {
  domain       = "mydomain"
  project_name = "myproject"
  pipeline_id  = 123456
  key          = "KEY"
  value        = "VAL"
  description  = "variable visibile only in this pipeline"
}