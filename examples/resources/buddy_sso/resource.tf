resource "buddy_sso" "saml" {
  domain      = "mydomain"
  sso_url     = "https://login.microsoftonline.com/test/saml2"
  issuer      = "https://sts.windows.net/test"
  certificate = "..."
  signature   = "sha256"
  digest      = "sha256"
}

resource "buddy_sso" "oidc" {
  domain        = "mydomain"
  issuer        = "https://sts.windows.net/test"
  client_id     = "12345"
  client_secret = "12345"
}