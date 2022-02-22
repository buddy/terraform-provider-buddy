resource "buddy_pipeline" "click" {
  domain              = "mydomain"
  project_name        = "myproject"
  name                = "click"
  on                  = "CLICK"
  refs                = ["main"]
  always_from_scratch = true
}

resource "buddy_pipeline" "event" {
  domain         = "mydomain"
  project_name   = "myproject"
  name           = "event"
  on             = "EVENT"
  priority       = "HIGH"
  fetch_all_refs = true

  event {
    type = "PUSH"
    refs = ["master"]
  }
}

resource "buddy_pipeline" "schedule" {
  domain                      = "mydomain"
  project_name                = "myproject"
  name                        = "schedule"
  on                          = "SCHEDULE"
  priority                    = "LOW"
  fail_on_prepare_env_warning = true
  start_date                  = "2016-11-18T12:38:16.000Z"
  delay                       = 10
  paused                      = true
}

resource "buddy_pipeline" "schedule_cron" {
  domain       = "mydomain"
  project_name = "myproject"
  name         = "schedule_cron"
  on           = "SCHEDULE"
  cron         = "15 14 1 * *"
}

resource "buddy_pipeline" "conditions" {
  domain       = "mydomain"
  project_name = "myproject"
  name         = "conditions"
  on           = "CLICK"
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
    zone_id   = "America/Monterrey"
  }
}