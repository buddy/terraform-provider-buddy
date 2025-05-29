# Get all targets in workspace
data "buddy_targets" "all" {
    domain = "myworkspace"
}

# Get targets in a specific project
data "buddy_targets" "project" {
    domain       = "myworkspace"
    project_name = "my-project"
}

# Get targets in a specific pipeline
data "buddy_targets" "pipeline" {
    domain      = "myworkspace"
    pipeline_id = 12345
}

# Get targets in a specific environment
data "buddy_targets" "environment" {
    domain         = "myworkspace"
    environment_id = "env123"
}

# Example usage - list all target names
output "all_target_names" {
    value = [for t in data.buddy_targets.all.targets : t.name]
}

# Example usage - filter SSH targets
output "ssh_targets" {
    value = [for t in data.buddy_targets.all.targets : t if t.type == "SSH"]
}

# Example usage - get disabled targets
output "disabled_targets" {
    value = [for t in data.buddy_targets.all.targets : t.name if t.disabled]
}