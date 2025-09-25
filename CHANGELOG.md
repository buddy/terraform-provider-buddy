# Release v1.35.0 (2025-09-30)
* Add new pipeline event - WEBHOOK

# Release v1.34.0 (2025-09-23)
* Add sandbox resource & data

# Release v1.33.0 (2025-08-13)
* Bump deps & go version to 1.23

# Release v1.32.1 (2025-08-13)
* Bump api sdk version to v1.30.1

# Release v1.32.0 (2025-07-23)
* Adds `remote_ref` to pipeline, `remote_branch` deprecated

# Release v1.31.0 (2025-07-22)
* Adds geolocation to domain record

# Release v1.30.0 (2025-06-06)
* [Breaking] - 'type' removed from environment
* [Breaking] - 'allowed_pipelines' removed from environment

# Release v1.29.1 (2025-06-03)
* Adds variable scope validator

# Release v1.29.0 (2025-06-02)
* Adds targets

# Release v1.28.1 (2025-05-29)
* Adds identifier to pipeline

# Release v1.28.0 (2025-05-21)
* Adds environments

# Release v1.27.0 (2025-04-09)
* Adds manage_variables_by_yaml, manage_permissions_by_yaml to pipeline
* Adds domain records management

# Release v1.26.0 (2025-02-04)
* [Breaking] - 'on' in pipeline is removed

# Release v1.25.0 (2025-01-28)
* [Breaking] - cron, start_date and delay in pipeline are moved to event
* [Breaking] - zone_id -> timezone in pipeline trigger condition

# Release v1.24.1 (2024-12-10)
* Adds cpu to pipeline

# Release v1.24.0 (2024-11-08)
* Fix variable import

# Release v1.23.0 (2024-10-15)
* Adds new props to pipeline
* concurrent_pipeline_runs
* description_required
* git_changeset_base
* filesystem_changeset_base

# Release v1.22.0 (2024-09-25)
* Adds new events to pipeline

# Release v1.21.0 (2024-07-18)
* Adds `pause_on_repeated_failures` to pipeline

# Release v1.20.1 (2024-07-09)
* goreleaser fix

# Release v1.20.0 (2024-07-09)
* Changes in resource integration
* [BREAKING] Removes scopes other than `WORKSPACE`, `PROJECT`
* Adds `permissions`
* Adds `allowed_pipelines`

# Release v1.19.0 (2024-04-19)
* Adds identifier to integration

# Release v1.18.0 (2024-04-12)
* Go version 1.21
* Bump deps

# Release v1.17.0 (2024-03-13)
* Adds GIT configuration to pipeline

# Release v1.16.0 (2023-09-13)
* Adds new trigger conditions to pipeline

# Release v1.15.1 (2023-08-16)
* Adds timeout in api client
 
# Release v1.15.0 (2023-08-01)
* Adds OIDC support in integrations & login

# Release v1.14.0 (2023-07-17)
* Adds `STACK_HAWK` integration type

# Release v1.13.1 (2023-06-16)
* Fixes without_repository in buddy_project

# Release v1.13.0 (2023-06-02)
* Bumps min Go version to 1.19
* Migrates from sdk v2 to terraform plugin framework (no breaking changes)

# Release v1.12.0 (2023-03-28)
* Adds `buddy_sso` resource to manage SSO in workspace

# Release v1.11.1 (2023-01-23)
* Fix `auto_assign_permission_set_id` when `auto_assign_to_new_projects` == `false`

# Release v1.11.0 (2023-01-17)
* Adds `partner_token` to `buddy_integration`

# Release v1.10.0 (2023-01-10)
* Adds `permissions` to `buddy_pipeline`

# Release v1.9.0 (2022-12-13)
* Remove `display_name` from `buddy_variable_ssh_key`
* Adds `without_repository` to `buddy_project`

# Release v1.8.0 (2022-11-29)
* Adds `fetch_submodules, fetch_submodules_env_key, access, allow_pull_requests` to `buddy_project`
* Adds new scope `PRIVATE_IN_PROJECT` to `buddy_integration`

# Release v1.7.0 (2022-07-20)
* Adds `update_default_branch_from_external` to `buddy_project`

# Release v1.6.0 (2022-06-29)
* Adds `status` (`MANAGER`, `MEMBER`) to `buddy_group_member`

# Release v1.5.0 (2022-05-31)
* Adds `project_team_access_level` to permission

# Release v1.4.6 (2022-05-25)
* Adds `custom_repo_ssh_key_id` when creating custom project

# Release v1.4.5 (2022-04-19)
* Adds `GitHub` & `GitLab` token integration

# Release v1.4.4 (2022-04-19)
* Adds possibility to assign `buddy_group` and `buddy_member` to new projects by default using provided `buddy_permission`
