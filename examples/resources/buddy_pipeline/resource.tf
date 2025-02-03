resource "buddy_pipeline" "click" {
  domain                   = "mydomain"
  project_name             = "myproject"
  name                     = "click"
  refs                     = ["main"]
  always_from_scratch      = true
  concurrent_pipeline_runs = true
  git_config_ref           = "NONE"
  permissions {
    others = "DENIED"
    user {
      id           = 1
      access_level = "READ_WRITE"
    }
  }
}

resource "buddy_pipeline" "event_push" {
  domain               = "mydomain"
  project_name         = "myproject"
  name                 = "event_push"
  priority             = "HIGH"
  fetch_all_refs       = true
  description_required = true

  event {
    type = "PUSH"
    refs = ["refs/heads/master"]
  }
}

resource "buddy_pipeline" "event_create_ref" {
  domain       = "mydomain"
  project_name = "myproject"
  name         = "event_create_ref"

  event {
    type = "CREATE_REF"
    refs = ["refs/heads/*"]
  }
}

resource "buddy_pipeline" "event_delete_ref" {
  domain       = "mydomain"
  project_name = "myproject"
  name         = "event_delete_ref"

  event {
    type = "DELETE_REF"
    refs = ["refs/heads/*"]
  }
}

resource "buddy_pipeline" "schedule" {
  domain       = "mydomain"
  project_name = "myproject"
  name         = "schedule"
  event {
    type       = "SCHEDULE"
    start_date = "2016-11-18T12:38:16.000Z"
    delay      = 10
  }
  priority                    = "LOW"
  fail_on_prepare_env_warning = true
  paused                      = true
  git_config_ref              = "FIXED"
  git_config = {
    project = "project_name"
    branch  = "branch_name"
    path    = "path/to/definition.yml"
  }
}

resource "buddy_pipeline" "schedule_cron" {
  domain       = "mydomain"
  project_name = "myproject"
  name         = "schedule_cron"
  event {
    cron = "15 14 1 * *"
  }
  git_config_ref = "DYNAMIC"
}

resource "buddy_pipeline" "remote" {
  domain              = "mydomain"
  project_name        = "myproject"
  name                = "remote_pipeline"
  definition_source   = "REMOTE"
  remote_project_name = "remote_project"
  remote_branch       = "remote_branch"
  remote_path         = "remote.yml"

  remote_parameter {
    key   = "myparam"
    value = "val"
  }
}

resource "buddy_pipeline" "conditions" {
  domain       = "mydomain"
  project_name = "myproject"
  name         = "conditions"
  refs         = ["main"]

  trigger_condition {
    condition = "ON_CHANGE"
  }
  trigger_condition {
    condition = "ON_CHANGE_AT_PATH"
    paths     = ["/abc"]
  }
  trigger_condition {
    condition      = "VAR_IS"
    variable_key   = "KEY"
    variable_value = "VAL"
  }
  trigger_condition {
    condition      = "VAR_IS_NOT"
    variable_key   = "KEY"
    variable_value = "VAL"
  }
  trigger_condition {
    condition      = "VAR_CONTAINS"
    variable_key   = "KEY"
    variable_value = "VAL"
  }
  trigger_condition {
    condition      = "VAR_NOT_CONTAINS"
    variable_key   = "KEY"
    variable_value = "VAL"
  }
  trigger_condition {
    condition = "DATETIME"
    hours     = [10]
    days      = [1, 20]
    timezone  = "America/Monterrey"
  }
  trigger_condition {
    condition     = "TRIGGERING_USER_IS_NOT_IN_GROUP"
    trigger_group = "devs"
  }
  trigger_condition {
    condition     = "TRIGGERING_USER_IS_IN_GROUP"
    trigger_group = "admins"
  }
  trigger_condition {
    condition    = "TRIGGERING_USER_IS_NOT"
    trigger_user = "test1@test.com"
  }
  trigger_condition {
    condition    = "TRIGGERING_USER_IS"
    trigger_user = "test2@test.com"
  }
}