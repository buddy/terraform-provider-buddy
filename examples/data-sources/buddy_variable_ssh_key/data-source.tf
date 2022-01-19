data "buddy_variable" "by_id" {
  variable_id = 123456
}

data "buddy_variable" "by_key" {
  key = "MY_KEY"
}

data "buddy_variable" "by_key_in_project" {
  key          = "MY_KEY"
  project_name = "myproject"
}



