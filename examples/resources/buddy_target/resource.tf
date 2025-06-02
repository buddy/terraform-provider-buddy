resource "buddy_target" "ftp" {
  domain     = "myworkspace"
  name       = "FTP Server"
  identifier = "ftp-server"
  type       = "FTP"
  host       = "ftp.example.com"
  port       = "21"
  secure     = false
  auth {
    username = "ftpuser"
    password = "secret"
  }
}

resource "buddy_target" "ftps_in_project" {
  domain       = "myworkspace"
  project_name = "myproject"
  name         = "FTPS Server"
  identifier   = "ftps-server"
  type         = "FTP"
  host         = "ftp.example.com"
  port         = "21"
  secure       = true
  auth {
    username = "ftpuser"
    password = "secret"
  }
}

resource "buddy_target" "ssh_password_in_environment" {
  domain         = "myworkspace"
  project_name   = "myproject"
  environment_id = "myenv"
  name           = "SSH Server"
  identifier     = "ssh-server"
  type           = "SSH"
  host           = "ssh.example.com"
  port           = "22"
  path           = "/var/www"
  auth {
    method   = "PASSWORD"
    username = "sshuser"
    password = "secret"
  }
}

resource "buddy_target" "ssh_key" {
  domain     = "myworkspace"
  name       = "SSH Server"
  identifier = "ssh-server"
  type       = "SSH"
  host       = "ssh.example.com"
  port       = "22"
  path       = "/var/www"
  auth {
    method     = "SSH_KEY"
    username   = "sshuser"
    key        = "-----BEGIN RSA PRIVATE KEY-----\n..."
    passphrase = "keypassphrase"
  }
}

resource "buddy_target" "ssh_asset" {
  domain     = "myworkspace"
  name       = "SSH Server with Asset"
  identifier = "ssh-server-asset"
  type       = "SSH"
  host       = "ssh.example.com"
  port       = "22"
  path       = "/var/www"
  auth {
    method   = "ASSETS_KEY"
    username = "sshuser"
    asset    = "id_workspace"
  }
}

resource "buddy_target" "ssh_proxy_in_pipeline" {
  domain       = "myworkspace"
  project_name = "myproject"
  pipeline_id  = 12345
  name         = "SSH via Proxy"
  identifier   = "ssh-proxy"
  type         = "SSH"
  host         = "internal.example.com"
  port         = "22"
  path         = "/var/www"
  auth {
    method = "PROXY_CREDENTIALS"
  }
  proxy {
    name = "Jump Host"
    host = "proxy.example.com"
    port = "22"
    auth {
      method   = "PASSWORD"
      username = "proxyuser"
      password = "proxypass"
    }
  }
}

resource "buddy_target" "ssh_proxy_key" {
  domain       = "myworkspace"
  project_name = "myproject"
  name         = "SSH via Proxy"
  identifier   = "ssh-proxy"
  type         = "SSH"
  host         = "internal.example.com"
  port         = "22"
  path         = "/var/www"
  auth {
    method   = "PROXY_KEY"
    username = "myuser"
  }
  proxy {
    name = "Jump Host"
    host = "proxy.example.com"
    port = "22"
    auth {
      method     = "SSH_KEY"
      username   = "proxyuser"
      key        = "-----BEGIN RSA PRIVATE KEY-----\n..."
      passphrase = "keypassphrase"
    }
  }
}

resource "buddy_target" "git_http" {
  domain     = "myworkspace"
  name       = "Git Repository"
  identifier = "git-repo"
  type       = "GIT"
  repository = "https://github.com/example/repo.git"
  tags       = ["production", "backend"]
  auth {
    method   = "HTTP"
    username = "gituser"
    password = "token"
  }
}

resource "buddy_target" "git_ssh" {
  domain     = "myworkspace"
  name       = "Git SSH Repository"
  identifier = "git-ssh-repo"
  type       = "GIT"
  repository = "git@github.com:example/repo.git"
  auth {
    method = "SSH_KEY"
    key    = "-----BEGIN RSA PRIVATE KEY-----\n..."
  }
}

resource "buddy_target" "digitalocean" {
  domain      = "myworkspace"
  name        = "DigitalOcean Server"
  identifier  = "do-server"
  type        = "DIGITAL_OCEAN"
  host        = "droplet.example.com"
  port        = "22"
  integration = buddy_integration.do.identifier
  auth {
    method   = "PASSWORD"
    username = "root"
    password = "secret"
  }
}

resource "buddy_target" "restricted" {
  domain     = "myworkspace"
  name       = "Restricted Server"
  identifier = "restricted-server"
  type       = "SSH"
  host       = "secure.example.com"
  port       = "22"
  auth {
    method   = "SSH_KEY"
    username = "deploy"
    key      = "-----BEGIN RSA PRIVATE KEY-----\n..."
  }
  permissions {
    others = "USE_ONLY"
    users {
      id           = buddy_member.developer.member_id
      access_level = "USE_ONLY"
    }
    groups {
      id           = buddy_group.devops.group_id
      access_level = "MANAGE"
    }
  }
}