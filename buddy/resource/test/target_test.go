package test

import (
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"
)

var targetIgnoreImportVerify = []string{
	"auth.0.password",
	"auth.0.key",
	"auth.0.passphrase",
	"proxy.0.auth.0.password",
	"proxy.0.auth.0.key",
	"proxy.0.auth.0.passphrase",
}

func TestAccTarget_ftps(t *testing.T) {
	var target buddy.Target
	domain := util.UniqueString()
	name := util.RandString(10)
	identifier := util.UniqueString()
	host := "1.1.1.1"
	port := "33"
	username := util.RandString(10)
	password := util.RandString(10)

	newName := util.RandString(10)
	newIdentifier := util.UniqueString()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create with required fields
			{
				Config: testAccTargetFtpsConfig(domain, name, identifier, host, port, username, password, true, true),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, domain, name, identifier, buddy.TargetTypeFtp, host, port, true, true),
					resource.TestCheckResourceAttr("buddy_target.test", "auth.0.method", buddy.TargetAuthMethodPassword),
					resource.TestCheckResourceAttr("buddy_target.test", "auth.0.username", username),
					resource.TestCheckResourceAttrSet("buddy_target.test", "auth.0.password"),
				),
			},
			// Update fields
			{
				Config: testAccTargetFtpsConfig(domain, newName, newIdentifier, host, port, username, password, false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, domain, newName, newIdentifier, buddy.TargetTypeFtp, host, port, false, false),
				),
			},
			// Import
			{
				ResourceName: "buddy_target.test",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["buddy_target.test"]
					if !ok {
						return "", fmt.Errorf("Not found: %s", "buddy_target.test")
					}
					return fmt.Sprintf("%s:%s", rs.Primary.Attributes["domain"], rs.Primary.Attributes["target_id"]), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: targetIgnoreImportVerify,
			},
		},
	})
}

func TestAccTarget_sshPassword(t *testing.T) {
	var target buddy.Target
	domain := util.UniqueString()
	name := util.RandString(10)
	identifier := util.UniqueString()
	host := "1.1.1.1"
	port := "44"
	path := util.RandString(10)
	username := util.RandString(10)
	password := util.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccTargetSshPasswordConfig(domain, name, identifier, host, port, path, username, password),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, domain, name, identifier, buddy.TargetTypeSsh, host, port, false, false),
					resource.TestCheckResourceAttr("buddy_target.test", "auth.0.method", buddy.TargetAuthMethodPassword),
					resource.TestCheckResourceAttr("buddy_target.test", "auth.0.username", username),
					resource.TestCheckResourceAttr("buddy_target.test", "path", path),
				),
			},
			// Import
			{
				ResourceName: "buddy_target.test",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["buddy_target.test"]
					if !ok {
						return "", fmt.Errorf("Not found: %s", "buddy_target.test")
					}
					return fmt.Sprintf("%s:%s", rs.Primary.Attributes["domain"], rs.Primary.Attributes["target_id"]), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: targetIgnoreImportVerify,
			},
		},
	})
}

func TestAccTarget_sshKey(t *testing.T) {
	var target buddy.Target
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name := util.RandString(10)
	identifier := util.UniqueString()
	host := "1.1.1.1"
	port := "44"
	path := util.RandString(10)
	username := util.RandString(10)
	key := util.RandString(10)
	passphrase := util.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccTargetSshKeyConfig(domain, projectName, name, identifier, host, port, path, username, key, passphrase),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, domain, name, identifier, buddy.TargetTypeSsh, host, port, false, false),
					resource.TestCheckResourceAttr("buddy_target.test", "auth.0.method", buddy.TargetAuthMethodSshKey),
					resource.TestCheckResourceAttr("buddy_target.test", "auth.0.username", username),
					resource.TestCheckResourceAttrSet("buddy_target.test", "auth.0.key"),
					resource.TestCheckResourceAttrSet("buddy_target.test", "auth.0.passphrase"),
					resource.TestCheckResourceAttr("buddy_target.test", "scope", buddy.TargetScopeProject),
				),
			},
			// Import
			{
				ResourceName: "buddy_target.test",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["buddy_target.test"]
					if !ok {
						return "", fmt.Errorf("Not found: %s", "buddy_target.test")
					}
					return fmt.Sprintf("%s:%s", rs.Primary.Attributes["domain"], rs.Primary.Attributes["target_id"]), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: targetIgnoreImportVerify,
			},
		},
	})
}

func TestAccTarget_sshAsset(t *testing.T) {
	var target buddy.Target
	domain := util.UniqueString()
	projectName := util.UniqueString()
	envName := util.UniqueString()
	envId := util.UniqueString()
	name := util.RandString(10)
	identifier := util.UniqueString()
	host := "1.1.1.1"
	port := "44"
	path := util.RandString(10)
	username := util.RandString(10)
	asset := "id_workspace"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccTargetSshAssetConfig(domain, projectName, envName, envId, name, identifier, host, port, path, username, asset),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, domain, name, identifier, buddy.TargetTypeSsh, host, port, false, false),
					resource.TestCheckResourceAttr("buddy_target.test", "auth.0.method", buddy.TargetAuthMethodAssetsKey),
					resource.TestCheckResourceAttr("buddy_target.test", "auth.0.username", username),
					resource.TestCheckResourceAttr("buddy_target.test", "auth.0.asset", asset),
					resource.TestCheckResourceAttr("buddy_target.test", "scope", buddy.TargetScopeEnvironment),
				),
			},
			// Import
			{
				ResourceName: "buddy_target.test",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["buddy_target.test"]
					if !ok {
						return "", fmt.Errorf("Not found: %s", "buddy_target.test")
					}
					return fmt.Sprintf("%s:%s", rs.Primary.Attributes["domain"], rs.Primary.Attributes["target_id"]), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: targetIgnoreImportVerify,
			},
		},
	})
}

func TestAccTarget_gitHttp(t *testing.T) {
	var target buddy.Target
	domain := util.UniqueString()
	name := util.RandString(10)
	identifier := util.UniqueString()
	repository := "https://a" + util.UniqueString() + ".com"
	username := util.RandString(10)
	password := util.RandString(10)
	tag1 := "a"
	tag2 := "b"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccTargetGitHttpConfig(domain, name, identifier, repository, username, password, tag1, tag2),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					resource.TestCheckResourceAttr("buddy_target.test", "name", name),
					resource.TestCheckResourceAttr("buddy_target.test", "identifier", identifier),
					resource.TestCheckResourceAttr("buddy_target.test", "type", buddy.TargetTypeGit),
					resource.TestCheckResourceAttr("buddy_target.test", "repository", repository),
					resource.TestCheckResourceAttr("buddy_target.test", "auth.0.method", buddy.TargetAuthMethodHttp),
					resource.TestCheckResourceAttr("buddy_target.test", "auth.0.username", username),
					resource.TestCheckResourceAttr("buddy_target.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("buddy_target.test", "tags.0", tag1),
					resource.TestCheckResourceAttr("buddy_target.test", "tags.1", tag2),
				),
			},
			// Import
			{
				ResourceName: "buddy_target.test",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["buddy_target.test"]
					if !ok {
						return "", fmt.Errorf("Not found: %s", "buddy_target.test")
					}
					return fmt.Sprintf("%s:%s", rs.Primary.Attributes["domain"], rs.Primary.Attributes["target_id"]), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: targetIgnoreImportVerify,
			},
		},
	})
}

func TestAccTarget_digitalOcean(t *testing.T) {
	var target buddy.Target
	domain := util.UniqueString()
	integrationName := util.UniqueString()
	integrationToken := util.UniqueString()
	name := util.RandString(10)
	identifier := util.UniqueString()
	host := "1.1.1.1"
	port := "44"
	username := util.RandString(10)
	password := util.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccTargetDigitalOceanConfig(domain, integrationName, integrationToken, name, identifier, host, port, username, password),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, domain, name, identifier, buddy.TargetTypeDigitalOcean, host, port, false, false),
					resource.TestCheckResourceAttrSet("buddy_target.test", "integration"),
				),
			},
			// Import
			{
				ResourceName: "buddy_target.test",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["buddy_target.test"]
					if !ok {
						return "", fmt.Errorf("Not found: %s", "buddy_target.test")
					}
					return fmt.Sprintf("%s:%s", rs.Primary.Attributes["domain"], rs.Primary.Attributes["target_id"]), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: targetIgnoreImportVerify,
			},
		},
	})
}

func TestAccTarget_permissions(t *testing.T) {
	var target buddy.Target
	domain := util.UniqueString()
	name := util.RandString(10)
	memberEmail := util.RandEmail()
	groupName := util.UniqueString()
	identifier := util.UniqueString()
	host := "1.1.1.1"
	port := "44"
	path := util.RandString(10)
	username := util.RandString(10)
	key := util.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create with permissions
			{
				Config: testAccTargetPermissionsConfig(domain, memberEmail, groupName, name, identifier, host, port, path, username, key),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, domain, name, identifier, buddy.TargetTypeSsh, host, port, false, false),
					resource.TestCheckResourceAttr("buddy_target.test", "permissions.0.others", buddy.TargetPermissionManage),
					resource.TestCheckResourceAttrSet("buddy_target.test", "permissions.0.users.0.id"),
					resource.TestCheckResourceAttr("buddy_target.test", "permissions.0.users.0.access_level", buddy.TargetPermissionUseOnly),
					resource.TestCheckResourceAttrSet("buddy_target.test", "permissions.0.groups.0.id"),
					resource.TestCheckResourceAttr("buddy_target.test", "permissions.0.groups.0.access_level", buddy.TargetPermissionManage),
				),
			},
		},
	})
}

func TestAccTarget_sshProxyCredentials(t *testing.T) {
	var target buddy.Target
	domain := util.UniqueString()
	projectName := util.UniqueString()
	pipelineName := util.UniqueString()
	name := util.RandString(10)
	identifier := util.UniqueString()
	host := "1.1.1.1"
	port := "44"
	path := util.RandString(10)
	proxyName := util.UniqueString()
	proxyHost := "2.2.2.2"
	proxyPort := "33"
	proxyUsername := util.RandString(10)
	proxyPassword := util.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccTargetSshProxyCredentialsConfig(domain, projectName, pipelineName, name, identifier, host, port, path, proxyName, proxyHost, proxyPort, proxyUsername, proxyPassword),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetGet("buddy_target.test", &target),
					testAccTargetAttributes("buddy_target.test", &target, domain, name, identifier, buddy.TargetTypeSsh, host, port, false, false),
					resource.TestCheckResourceAttr("buddy_target.test", "auth.0.method", buddy.TargetAuthMethodProxyCredentials),
					resource.TestCheckResourceAttr("buddy_target.test", "proxy.0.name", proxyName),
					resource.TestCheckResourceAttr("buddy_target.test", "proxy.0.host", proxyHost),
					resource.TestCheckResourceAttr("buddy_target.test", "proxy.0.port", proxyPort),
					resource.TestCheckResourceAttr("buddy_target.test", "proxy.0.auth.0.method", buddy.TargetAuthMethodPassword),
					resource.TestCheckResourceAttr("buddy_target.test", "proxy.0.auth.0.username", proxyUsername),
					resource.TestCheckResourceAttr("buddy_target.test", "scope", buddy.TargetScopePipeline),
				),
			},
			// Import
			{
				ResourceName: "buddy_target.test",
				ImportState:  true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources["buddy_target.test"]
					if !ok {
						return "", fmt.Errorf("Not found: %s", "buddy_target.test")
					}
					return fmt.Sprintf("%s:%s", rs.Primary.Attributes["domain"], rs.Primary.Attributes["target_id"]), nil
				},
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: targetIgnoreImportVerify,
			},
		},
	})
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
			return util.ErrorResourceExists()
		}
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}

func testAccTargetGet(n string, target *buddy.Target) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		domain := rs.Primary.Attributes["domain"]
		targetId := rs.Primary.Attributes["target_id"]
		t, _, err := acc.ApiClient.TargetService.Get(domain, targetId)
		if err != nil {
			return err
		}
		*target = *t
		return nil
	}
}

func testAccTargetAttributes(n string, target *buddy.Target, domain string, name string, identifier string, typ string, host string, port string, secure bool, disabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if err := util.CheckFieldEqual("name", target.Name, name); err != nil {
			return err
		}

		if err := util.CheckFieldEqual("identifier", target.Identifier, identifier); err != nil {
			return err
		}

		if err := util.CheckFieldEqual("type", target.Type, typ); err != nil {
			return err
		}

		if err := util.CheckFieldEqual("host", target.Host, host); err != nil {
			return err
		}

		if err := util.CheckFieldEqual("port", target.Port, port); err != nil {
			return err
		}

		if err := util.CheckBoolFieldEqual("secure", target.Secure, secure); err != nil {
			return err
		}

		if err := util.CheckBoolFieldEqual("disabled", target.Disabled, disabled); err != nil {
			return err
		}

		attribs := map[string]string{
			"domain":     domain,
			"target_id":  target.Id,
			"name":       name,
			"identifier": identifier,
			"type":       typ,
			"host":       host,
			"port":       port,
			"secure":     strconv.FormatBool(secure),
			"disabled":   strconv.FormatBool(disabled),
		}

		for key, value := range attribs {
			if err := util.CheckFieldEqual(key, rs.Primary.Attributes[key], value); err != nil {
				return err
			}
		}

		return nil
	}
}

func testAccTargetFtpsConfig(domain string, name string, identifier string, host string, port string, username string, password string, secure bool, disabled bool) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_target" "test" {
    domain     = buddy_workspace.test.domain
    name       = "%s"
    identifier = "%s"
    type       = "%s"
    host       = "%s"
    port       = "%s"
    secure     = %t
    disabled   = %t
    auth {
        method   = "%s"
        username = "%s"
        password = "%s"
    }
}`, domain, name, identifier, buddy.TargetTypeFtp, host, port, secure, disabled, buddy.TargetAuthMethodPassword, username, password)
}

func testAccTargetSshPasswordConfig(domain string, name string, identifier string, host string, port string, path string, username string, password string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_target" "test" {
    domain     = buddy_workspace.test.domain
    name       = "%s"
    identifier = "%s"
    type       = "%s"
    host       = "%s"
    port       = "%s"
    path       = "%s"
    auth {
        method   = "%s"
        username = "%s"
        password = "%s"
    }
}`, domain, name, identifier, buddy.TargetTypeSsh, host, port, path, buddy.TargetAuthMethodPassword, username, password)
}

func testAccTargetSshKeyConfig(domain string, projectName string, name string, identifier string, host string, port string, path string, username string, key string, passphrase string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_project" "test" {
    domain       = buddy_workspace.test.domain
    display_name = "%s"
}

resource "buddy_target" "test" {
    domain       = buddy_workspace.test.domain
    project_name = buddy_project.test.name
    name         = "%s"
    identifier   = "%s"
    type         = "%s"
    host         = "%s"
    port         = "%s"
    path         = "%s"
    auth {
        method     = "%s"
        username   = "%s"
        key        = "%s"
        passphrase = "%s"
    }
}`, domain, projectName, name, identifier, buddy.TargetTypeSsh, host, port, path, buddy.TargetAuthMethodSshKey, username, key, passphrase)
}

func testAccTargetSshAssetConfig(domain string, projectName string, envName string, envId string, name string, identifier string, host string, port string, path string, username string, asset string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_project" "test" {
    domain       = buddy_workspace.test.domain
    display_name = "%s"
}

resource "buddy_environment" "test" {
    domain       = buddy_workspace.test.domain
    project_name = buddy_project.test.name
    name         = "%s"
    identifier   = "%s"
    type         = "%s"
}

resource "buddy_target" "test" {
    domain          = buddy_workspace.test.domain
    environment_id  = buddy_environment.test.environment_id
    name            = "%s"
    identifier      = "%s"
    type            = "%s"
    host            = "%s"
    port            = "%s"
    path            = "%s"
    auth {
        method   = "%s"
        username = "%s"
        asset    = "%s"
    }
}`, domain, projectName, envName, envId, buddy.EnvironmentTypeDev, name, identifier, buddy.TargetTypeSsh, host, port, path, buddy.TargetAuthMethodAssetsKey, username, asset)
}

func testAccTargetGitHttpConfig(domain string, name string, identifier string, repository string, username string, password string, tag1 string, tag2 string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_target" "test" {
    domain     = buddy_workspace.test.domain
    name       = "%s"
    identifier = "%s"
    type       = "%s"
    repository = "%s"
    tags       = ["%s", "%s"]
    auth {
        method   = "%s"
        username = "%s"
        password = "%s"
    }
}`, domain, name, identifier, buddy.TargetTypeGit, repository, tag1, tag2, buddy.TargetAuthMethodHttp, username, password)
}

func testAccTargetDigitalOceanConfig(domain string, integrationName string, integrationToken string, name string, identifier string, host string, port string, username string, password string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_integration" "test" {
    domain = buddy_workspace.test.domain
    name   = "%s"
    type   = "%s"
    scope  = "%s"
    token  = "%s"
}

resource "buddy_target" "test" {
    domain      = buddy_workspace.test.domain
    name        = "%s"
    identifier  = "%s"
    type        = "%s"
    host        = "%s"
    port        = "%s"
    integration = buddy_integration.test.identifier
    auth {
        method   = "%s"
        username = "%s"
        password = "%s"
    }
}`, domain, integrationName, buddy.IntegrationTypeDigitalOcean, buddy.IntegrationScopeWorkspace, integrationToken, name, identifier, buddy.TargetTypeDigitalOcean, host, port, buddy.TargetAuthMethodPassword, username, password)
}

func testAccTargetPermissionsConfig(domain string, memberEmail string, groupName string, name string, identifier string, host string, port string, path string, username string, key string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_member" "test" {
    domain = buddy_workspace.test.domain
    email  = "%s"
}

resource "buddy_group" "test" {
    domain = buddy_workspace.test.domain
    name   = "%s"
}

resource "buddy_target" "test" {
    domain     = buddy_workspace.test.domain
    name       = "%s"
    identifier = "%s"
    type       = "%s"
    host       = "%s"
    port       = "%s"
    path       = "%s"
    auth {
        method   = "%s"
        username = "%s"
        key      = "%s"
    }
    permissions {
        others = "%s"
        users {
            id           = buddy_member.test.member_id
            access_level = "%s"
        }
        groups {
            id           = buddy_group.test.group_id
            access_level = "%s"
        }
    }
}`, domain, memberEmail, groupName, name, identifier, buddy.TargetTypeSsh, host, port, path, buddy.TargetAuthMethodSshKey, username, key, buddy.TargetPermissionManage, buddy.TargetPermissionUseOnly, buddy.TargetPermissionManage)
}

func testAccTargetSshProxyCredentialsConfig(domain string, projectName string, pipelineName string, name string, identifier string, host string, port string, path string, proxyName string, proxyHost string, proxyPort string, proxyUsername string, proxyPassword string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_project" "test" {
    domain       = buddy_workspace.test.domain
    display_name = "%s"
}

resource "buddy_pipeline" "test" {
    domain       = buddy_workspace.test.domain
    project_name = buddy_project.test.name
    name         = "%s"
}

resource "buddy_target" "test" {
    domain      = buddy_workspace.test.domain
    pipeline_id = buddy_pipeline.test.pipeline_id
    name        = "%s"
    identifier  = "%s"
    type        = "%s"
    host        = "%s"
    port        = "%s"
    path        = "%s"
    auth {
        method = "%s"
    }
    proxy {
        name = "%s"
        host = "%s"
        port = "%s"
        auth {
            method   = "%s"
            username = "%s"
            password = "%s"
        }
    }
}`, domain, projectName, pipelineName, name, identifier, buddy.TargetTypeSsh, host, port, path, buddy.TargetAuthMethodProxyCredentials, proxyName, proxyHost, proxyPort, buddy.TargetAuthMethodPassword, proxyUsername, proxyPassword)
}
