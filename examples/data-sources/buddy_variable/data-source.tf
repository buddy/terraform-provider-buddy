data "buddy_variable_ssh_key" "by_id" {
  variable_id = 123456
}

data "buddy_variable_ssh_key" "by_key" {
  key = "MY_KEY"
}

data "buddy_variable_ssh_key" "by_key_in_project" {
  key          = "MY_KEY"
  project_name = "myproject"
}



