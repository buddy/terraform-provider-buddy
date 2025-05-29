# Create SSH target in workspace
resource "buddy_target" "ssh_target" {
  domain       = "myworkspace"
  name         = "Production Server"
  type         = "SSH"
  hostname     = "example.com"
  port         = 22
  username     = "deploy"
  password     = var.ssh_password
  file_path    = "/var/www/html"
  description  = "Production web server"
  tags         = ["production", "web"]
}

# Create SSH target with key authentication in project
resource "buddy_target" "ssh_key_target" {
  domain        = "myworkspace"
  project_name  = "my-project"
  name          = "Staging Server"
  type          = "SSH"
  hostname      = "staging.example.com"
  port          = 2222
  username      = "deploy"
  key_id        = buddy_variable_ssh_key.deploy_key.variable_id
  auth_mode     = "KEY"
  file_path     = "/var/www/staging"
  all_pipelines_allowed = false
  allowed_pipelines     = [buddy_pipeline.deploy.pipeline_id]
}

# Create FTP target
resource "buddy_target" "ftp_target" {
  domain   = "myworkspace"
  name     = "FTP Server"
  type     = "FTP"
  hostname = "ftp.example.com"
  port     = 21
  username = "ftpuser"
  password = var.ftp_password
}

# Create S3 target
resource "buddy_target" "s3_target" {
  domain = "myworkspace"
  name   = "S3 Bucket"
  type   = "AMAZON_S3"
}

# Create Docker Registry target
resource "buddy_target" "docker_target" {
  domain   = "myworkspace"
  name     = "Docker Registry"
  type     = "DOCKER_REGISTRY"
  hostname = "registry.example.com"
  username = "dockeruser"
  password = var.docker_password
}