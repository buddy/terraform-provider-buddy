package test

import (
	"fmt"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccTarget_workspace_ssh(t *testing.T) {
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	host := "example.com"
	newHost := "newexample.com"
	username := "testuser"
	newUsername := "newtestuser"
	port := "22"
	newPort := "2222"
	newPath := "/var/www/public"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create with minimal config
			{
				Config: testAccTargetConfigWorkspaceSsh(domain, name, host, port, username),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetCheckExist("buddy_target.test", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "domain", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "name", name),
					resource.TestCheckResourceAttr("buddy_target.test", "type", "SSH"),
					resource.TestCheckResourceAttr("buddy_target.test", "host", host),
					resource.TestCheckResourceAttr("buddy_target.test", "port", port),
					resource.TestCheckResourceAttr("buddy_target.test", "auth_username", username),
					resource.TestCheckResourceAttr("buddy_target.test", "auth_method", "PASS"),
					resource.TestCheckResourceAttr("buddy_target.test", "scope", "PRIVATE"),
				),
			},
			// Update with full config
			{
				Config: testAccTargetConfigWorkspaceSshFull(domain, newName, newHost, newPort, newUsername, newPath),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetCheckExist("buddy_target.test", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "domain", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "name", newName),
					resource.TestCheckResourceAttr("buddy_target.test", "type", "SSH"),
					resource.TestCheckResourceAttr("buddy_target.test", "host", newHost),
					resource.TestCheckResourceAttr("buddy_target.test", "port", newPort),
					resource.TestCheckResourceAttr("buddy_target.test", "auth_username", newUsername),
					resource.TestCheckResourceAttr("buddy_target.test", "path", newPath),
					resource.TestCheckResourceAttr("buddy_target.test", "auth_method", "KEY"),
					resource.TestCheckResourceAttr("buddy_target.test", "disabled", "false"),
					resource.TestCheckResourceAttr("buddy_target.test", "tags.#", "2"),
				),
			},
			// Import
			{
				ResourceName:            "buddy_target.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_password", "auth_passphrase", "auth_key"},
			},
		},
	})
}

func TestAccTarget_workspace_ftp(t *testing.T) {
	domain := util.UniqueString()
	name := util.RandString(10)
	host := "ftp.example.com"
	username := "ftpuser"
	port := "21"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create FTP target
			{
				Config: testAccTargetConfigWorkspaceFtp(domain, name, host, port, username),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetCheckExist("buddy_target.test", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "domain", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "name", name),
					resource.TestCheckResourceAttr("buddy_target.test", "type", "FTP"),
					resource.TestCheckResourceAttr("buddy_target.test", "host", host),
					resource.TestCheckResourceAttr("buddy_target.test", "port", port),
					resource.TestCheckResourceAttr("buddy_target.test", "auth_username", username),
				),
			},
			// Import
			{
				ResourceName:            "buddy_target.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_password"},
			},
		},
	})
}

func TestAccTarget_workspace_s3(t *testing.T) {
	domain := util.UniqueString()
	name := util.RandString(10)
	repository := "my-bucket"
	path := "/uploads"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create S3 target
			{
				Config: testAccTargetConfigWorkspaceS3(domain, name, repository, path),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetCheckExist("buddy_target.test", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "domain", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "name", name),
					resource.TestCheckResourceAttr("buddy_target.test", "type", "AMAZON_S3"),
					resource.TestCheckResourceAttr("buddy_target.test", "repository", repository),
					resource.TestCheckResourceAttr("buddy_target.test", "path", path),
					resource.TestCheckResourceAttrSet("buddy_target.test", "integration"),
				),
			},
			// Import
			{
				ResourceName:      "buddy_target.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTarget_workspace_docker_registry(t *testing.T) {
	domain := util.UniqueString()
	name := util.RandString(10)
	host := "docker.io"
	repository := "myorg/myapp"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create Docker Registry target
			{
				Config: testAccTargetConfigWorkspaceDockerRegistry(domain, name, host, repository),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetCheckExist("buddy_target.test", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "domain", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "name", name),
					resource.TestCheckResourceAttr("buddy_target.test", "type", "DOCKER_REGISTRY"),
					resource.TestCheckResourceAttr("buddy_target.test", "host", host),
					resource.TestCheckResourceAttr("buddy_target.test", "repository", repository),
					resource.TestCheckResourceAttr("buddy_target.test", "secure", "true"),
				),
			},
			// Import
			{
				ResourceName:            "buddy_target.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auth_password"},
			},
		},
	})
}

func testAccTargetConfigWorkspaceSsh(domain, name, host, port, username string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
	domain = "%s"
}

resource "buddy_target" "test" {
	domain = buddy_workspace.test.domain
	name = "%s"
	type = "SSH"
	host = "%s"
	port = "%s"
	auth_method = "PASS"
	auth_username = "%s"
	auth_password = "secret123"
}
`, domain, name, host, port, username)
}

func testAccTargetConfigWorkspaceSshFull(domain, name, host, port, username, path string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
	domain = "%s"
}

resource "buddy_target" "test" {
	domain = buddy_workspace.test.domain
	name = "%s"
	type = "SSH"
	host = "%s"
	port = "%s"
	path = "%s"
	auth_method = "KEY"
	auth_username = "%s"
	auth_key = "-----BEGIN RSA PRIVATE KEY-----\ntest\n-----END RSA PRIVATE KEY-----"
	auth_passphrase = "passphrase123"
	disabled = false
	tags = ["production", "server"]
}
`, domain, name, host, port, path, username)
}

func testAccTargetConfigWorkspaceFtp(domain, name, host, port, username string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
	domain = "%s"
}

resource "buddy_target" "test" {
	domain = buddy_workspace.test.domain
	name = "%s"
	type = "FTP"
	host = "%s"
	port = "%s"
	auth_username = "%s"
	auth_password = "secret123"
}
`, domain, name, host, port, username)
}

func testAccTargetConfigWorkspaceS3(domain, name, repository, path string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
	domain = "%s"
}

data "buddy_integrations" "aws" {
	domain = buddy_workspace.test.domain
	type = "AMAZON"
}

resource "buddy_target" "test" {
	domain = buddy_workspace.test.domain
	name = "%s"
	type = "AMAZON_S3"
	repository = "%s"
	path = "%s"
	integration = data.buddy_integrations.aws.integrations[0].integration_id
}
`, domain, name, repository, path)
}

func testAccTargetConfigWorkspaceDockerRegistry(domain, name, host, repository string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
	domain = "%s"
}

resource "buddy_target" "test" {
	domain = buddy_workspace.test.domain
	name = "%s"
	type = "DOCKER_REGISTRY"
	host = "%s"
	repository = "%s"
	secure = true
	auth_username = "dockeruser"
	auth_password = "dockerpass"
}
`, domain, name, host, repository)
}

func testAccTargetCheckExist(n string, domain string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		targetId := rs.Primary.Attributes["target_id"]
		if targetId == "" {
			return fmt.Errorf("no target ID is set")
		}

		_, _, err := acc.ApiClient.TargetService.Get(domain, targetId)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccTargetCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_target" {
			continue
		}

		domain := rs.Primary.Attributes["domain"]
		targetId := rs.Primary.Attributes["target_id"]

		target, resp, err := acc.ApiClient.TargetService.Get(domain, targetId)
		if err == nil && target != nil {
			return fmt.Errorf("target still exists")
		}

		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}

	return nil
}
