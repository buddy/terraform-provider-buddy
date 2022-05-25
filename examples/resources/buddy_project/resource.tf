resource "buddy_project" "test" {
  domain       = "mydomain"
  display_name = "test"
}

resource "buddy_project" "github" {
  domain              = "mydomain"
  display_name        = "github"
  integration_id      = "githubIntegrationId"
  external_project_id = "account/project"
}

resource "buddy_project" "bitbucket" {
  domain              = "mydomain"
  display_name        = "bitbucket"
  integration_id      = "bitbucketIntegrationId"
  external_project_id = "account/project"
}

resource "buddy_project" "gitlab" {
  domain              = "mydomain"
  display_name        = "gitlab"
  integration_id      = "gitlabIntegrationId"
  external_project_id = "account/project"
  git_lab_project_id  = "12345678"
}

resource "buddy_project" "custom_http" {
  domain           = "mydomain"
  display_name     = "custom"
  custom_repo_url  = "https://mygit.repo"
  custom_repo_user = "user"
  custom_repo_pass = "pass"
}

resource "buddy_project" "custom_ssh" {
  domain                 = "mydomain"
  display_name           = "custom"
  custom_repo_url        = "ssh://mygit.repo"
  custom_repo_ssh_key_id = 12345
}