# FTP/FTPS target
resource "buddy_target" "ftp" {
    domain     = "myworkspace"
    name       = "FTP Server"
    identifier = "ftp-server"
    type       = "FTP"
    host       = "ftp.example.com"
    port       = "21"
    secure     = true
    disabled   = false
    auth {
        method   = "PASSWORD"
        username = "ftpuser"
        password = "secret"
    }
}

# SSH target with password
resource "buddy_target" "ssh_password" {
    domain     = "myworkspace"
    name       = "SSH Server"
    identifier = "ssh-server"
    type       = "SSH"
    host       = "ssh.example.com"
    port       = "22"
    path       = "/var/www"
    auth {
        method   = "PASSWORD"
        username = "sshuser"
        password = "secret"
    }
}

# SSH target with SSH key (project-scoped)
resource "buddy_target" "ssh_key" {
    domain       = "myworkspace"
    project_name = "my-project"
    name         = "SSH Server with Key"
    identifier   = "ssh-server-key"
    type         = "SSH"
    host         = "ssh.example.com"
    port         = "22"
    path         = "/var/www"
    auth {
        method     = "SSH_KEY"
        username   = "sshuser"
        key        = "-----BEGIN RSA PRIVATE KEY-----\n..."
        passphrase = "keypassphrase"
    }
}

# SSH target with asset key (environment-scoped)
resource "buddy_target" "ssh_asset" {
    domain         = "myworkspace"
    environment_id = "env123"
    name           = "SSH Server with Asset"
    identifier     = "ssh-server-asset"
    type           = "SSH"
    host           = "ssh.example.com"
    port           = "22"
    path           = "/var/www"
    auth {
        method   = "ASSETS_KEY"
        username = "sshuser"
        asset    = "id_workspace"
    }
}

# Git target with HTTP auth
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

# Git target with SSH key
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

# DigitalOcean target
resource "buddy_integration" "do" {
    domain = "myworkspace"
    name   = "DigitalOcean"
    type   = "DIGITAL_OCEAN"
    scope  = "WORKSPACE"
    token  = "digitalocean-api-token"
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

# SSH target with proxy (pipeline-scoped)
resource "buddy_target" "ssh_proxy" {
    domain      = "myworkspace"
    pipeline_id = 12345
    name        = "SSH via Proxy"
    identifier  = "ssh-proxy"
    type        = "SSH"
    host        = "internal.example.com"
    port        = "22"
    path        = "/var/www"
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

# Target with permissions
resource "buddy_member" "developer" {
    domain = "myworkspace"
    email  = "developer@example.com"
}

resource "buddy_group" "devops" {
    domain = "myworkspace"
    name   = "DevOps Team"
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