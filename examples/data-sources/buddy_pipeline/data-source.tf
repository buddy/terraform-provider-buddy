data "buddy_pipeline" "by_name" {
  domain       = "mydomain"
  project_name = "myproject"
  name         = "mypipeline"
}

data "buddy_pipeline" "by_id" {
  domain       = "mydomain"
  project_name = "myproject"
  pipeline_id  = 123456
}



