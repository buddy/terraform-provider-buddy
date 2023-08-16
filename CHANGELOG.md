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