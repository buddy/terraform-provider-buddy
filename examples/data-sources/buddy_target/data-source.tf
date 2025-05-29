# Get target by ID
data "buddy_target" "by_id" {
    domain    = "myworkspace"
    target_id = "abc123"
}

# Get target by identifier
data "buddy_target" "by_identifier" {
    domain     = "myworkspace"
    identifier = "my-target"
}

# Example usage
output "target_name" {
    value = data.buddy_target.by_identifier.name
}

output "target_type" {
    value = data.buddy_target.by_identifier.type
}

output "target_host" {
    value = data.buddy_target.by_identifier.host
}

output "target_url" {
    value = data.buddy_target.by_identifier.html_url
}