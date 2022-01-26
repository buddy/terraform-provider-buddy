resource "buddy_pipeline" "test" {
  domain       = "mydomain"
  project_name = "myproject"
  name         = "test"
  on           = "CLICK"
  refs         = ["main"]
}