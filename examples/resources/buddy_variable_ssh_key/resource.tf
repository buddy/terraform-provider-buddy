resource "buddy_variable_ssh_key" "mykey" {
  domain     = "mydomain"
  key        = "MY_KEY"
  value      = <<EOT
-----BEGIN PRIVATE KEY-----
...
-----END PRIVATE KEY-----
EOT
  file_place = "CONTAINER"
  file_path  = "~/.ssh/id_mykey"
  file_chmod = "600"
}