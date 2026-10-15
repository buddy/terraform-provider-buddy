resource "buddy_sandbox" "sb" {
  domain              = "mydomain"
  project_name        = "test"
  name                = "sb"
  install_commands    = "apt-get update && apt-get install -y curl"
  wait_for_running    = true
  wait_for_configured = true
  endpoints = {
    "local" = {
      endpoint = "localhost:3333"
      type     = "HTTP"
    }
  }
}
resource "buddy_sandbox" "fetch" {
  domain              = "mydomain"
  project_name        = "test"
  name                = "fetch"
  wait_for_running    = true
  wait_for_configured = true

  app {
    identifier = "web"
    command    = "npm run dev"
    path       = "/app"
  }

  fetch {
    type          = "PROJECT_REPO"
    ref           = "main"
    path          = "/app"
    build_command = "npm ci"
  }

  fetch {
    type       = "PUBLIC_REPO"
    repository = "https://github.com/octocat/Hello-World"
    path       = "/vendor/hello-world"
  }

  fetch {
    type        = "ARTIFACT"
    artifact    = "my-artifact:latest"
    source_path = "dist.zip"
    path        = "/data"
  }
}
