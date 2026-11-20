resource "buddy_ruleset" "protect_main" {
  domain       = "mydomain"
  project_name = "myproject"
  identifier   = "protect-main"
  name         = "Protect main"
  note         = "Main changes only through pull requests with a green build"
  agent_note   = "Open a pull request instead of pushing to main"
  refs         = [":default"]

  bypass_admins = true
  bypass_others = "DENIED"

  bypass_user {
    id     = 123
    access = "ALLOWED"
  }

  bypass_group {
    id     = 456
    access = "ALLOWED"
  }

  required_pipelines           = ["build"]
  required_unit_tests_suites   = ["backend-junit"]
  required_visual_tests_suites = ["storybook"]
  require_up_to_date           = true
  require_pull_request         = true
  block_force_push             = true
  restrict_deletion            = true
}

resource "buddy_ruleset" "release_tags" {
  domain            = "mydomain"
  project_name      = "myproject"
  identifier        = "release-tags"
  refs              = ["refs/tags/v*"]
  bypass_admins     = true
  restrict_creation = true
  restrict_update   = true
  restrict_deletion = true
}
