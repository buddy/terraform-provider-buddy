---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "buddy_pipeline Resource - buddy-terraform"
subcategory: ""
description: |-
  Create and manage a pipeline
  Token scopes required: WORKSPACE, EXECUTION_MANAGE, EXECUTION_INFO
---

# buddy_pipeline (Resource)

Create and manage a pipeline

Token scopes required: `WORKSPACE`, `EXECUTION_MANAGE`, `EXECUTION_INFO`

## Example Usage

```terraform
resource "buddy_pipeline" "test" {
  domain       = "mydomain"
  project_name = "myproject"
  name         = "test"
  on           = "CLICK"
  refs         = ["main"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **domain** (String) The workspace's URL handle
- **name** (String) The pipeline's name
- **on** (String) The pipeline's trigger mode. Allowed: `CLICK`, `EVENT`, `SCHEDULE`
- **project_name** (String) The project's name

### Optional

- **always_from_scratch** (Boolean) Defines whether or not to upload everything from scratch on every run
- **auto_clear_cache** (Boolean) Defines whether or not to automatically clear cache before running the pipeline
- **cron** (String) The pipeline's CRON expression. Required if the pipeline is set to `on: SCHEDULE` and neither `start_date` nor `delay` is specified
- **delay** (Number) The pipeline's runs interval (in minutes). Required if the pipeline is set to `on: SCHEDULE` and no `cron` is specified
- **do_not_create_commit_status** (Boolean) Defines whether or not to omit sending commit statuses to GitHub or GitLab upon execution
- **event** (Block List) The pipeline's list of events. Set it if `on: EVENT` (see [below for nested schema](#nestedblock--event))
- **execution_message_template** (String) The pipeline's run title. Default: `$BUDDY_EXECUTION_REVISION_SUBJECT`
- **fail_on_prepare_env_warning** (Boolean) Defines either or not run should fail if any warning occurs in prepare environment. Default: `true`
- **fetch_all_refs** (Boolean) Defines either or not fetch all refs from repository. Default: `true`
- **ignore_fail_on_project_status** (Boolean) If set to true the status of a given pipeline will be ignored on the projects' dashboard
- **no_skip_to_most_recent** (Boolean) Defines whether or not to skip run to the most recent run
- **paused** (Boolean) Is the pipeline's run paused. Restricted to `on: SCHEDULE`
- **priority** (String) The pipeline's priority. Allowed: `LOW`, `NORMAL`, `HIGH`
- **refs** (Set of String) The pipeline's list of refs. Set it if `on: CLICK`
- **start_date** (String) The pipeline's start date. Required if the pipeline is set to `on: SCHEDULE` and no `cron` is specified. Format: `2016-11-18T12:38:16.000Z`
- **target_site_url** (String) The pipeline's website target URL
- **trigger_condition** (Block List) The pipeline's list of trigger conditions (see [below for nested schema](#nestedblock--trigger_condition))
- **worker** (String) The pipeline's worker name. Only for `Buddy Enterprise`

### Read-Only

- **create_date** (String) The pipeline's date of creation
- **creator** (List of Object) The pipeline's creator (see [below for nested schema](#nestedatt--creator))
- **html_url** (String) The pipeline's URL
- **id** (String) The Terraform resource identifier for this item
- **last_execution_revision** (String) The pipeline's last run revision
- **last_execution_status** (String) The pipeline's last run status
- **pipeline_id** (Number) The pipeline's ID
- **project** (List of Object) The pipeline's project (see [below for nested schema](#nestedatt--project))

<a id="nestedblock--event"></a>
### Nested Schema for `event`

Required:

- **refs** (Set of String)
- **type** (String)


<a id="nestedblock--trigger_condition"></a>
### Nested Schema for `trigger_condition`

Required:

- **condition** (String)

Optional:

- **days** (Set of Number)
- **hours** (Set of Number)
- **paths** (Set of String)
- **pipeline_name** (String)
- **project_name** (String)
- **variable_key** (String)
- **variable_value** (String, Sensitive)
- **zone_id** (String)


<a id="nestedatt--creator"></a>
### Nested Schema for `creator`

Read-Only:

- **admin** (Boolean)
- **avatar_url** (String)
- **email** (String)
- **html_url** (String)
- **member_id** (Number)
- **name** (String)
- **workspace_owner** (Boolean)


<a id="nestedatt--project"></a>
### Nested Schema for `project`

Read-Only:

- **display_name** (String)
- **html_url** (String)
- **name** (String)
- **status** (String)

## Import

Import is supported using the following syntax:

```shell
# import using domain(mydomain), project name (myproject) and pipeline id (123456)
terraform import buddy_pipeline.test mydomain:myproject:123456
```