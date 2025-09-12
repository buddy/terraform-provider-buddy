resource "buddy_sandbox_status" "s1" {
  domain          = "mydomain"
  sandbox_id      = "12345"
  status          = "STOPPED"
  wait_for_status = true
}

resource "buddy_sandbox_status" "s2" {
  domain     = "mydomain"
  sandbox_id = "54321"
  status     = "RUNNING"
}