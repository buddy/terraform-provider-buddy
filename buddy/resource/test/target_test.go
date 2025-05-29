package test

import (
	"fmt"
	"strconv"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"

	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccTarget_workspace_ssh(t *testing.T) {
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	hostname := "example.com"
	newHostname := "newexample.com"
	username := "testuser"
	newUsername := "newtestuser"
	port := 22
	newPort := 2222
	description := "Test SSH target"
	newDescription := "Updated SSH target"
	filePath := "/var/www/html"
	newFilePath := "/var/www/public"
	
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create with minimal config
			{
				Config: testAccTargetConfigWorkspaceSsh(domain, name, hostname, port, username),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetCheckExist("buddy_target.test", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "domain", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "name", name),
					resource.TestCheckResourceAttr("buddy_target.test", "type", "SSH"),
					resource.TestCheckResourceAttr("buddy_target.test", "hostname", hostname),
					resource.TestCheckResourceAttr("buddy_target.test", "port", strconv.Itoa(port)),
					resource.TestCheckResourceAttr("buddy_target.test", "username", username),
					resource.TestCheckResourceAttr("buddy_target.test", "auth_mode", "PASS"),
					resource.TestCheckResourceAttr("buddy_target.test", "all_pipelines_allowed", "true"),
				),
			},
			// Update with full config
			{
				Config: testAccTargetConfigWorkspaceSshFull(domain, newName, newHostname, newPort, newUsername, newDescription, newFilePath),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetCheckExist("buddy_target.test", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "domain", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "name", newName),
					resource.TestCheckResourceAttr("buddy_target.test", "type", "SSH"),
					resource.TestCheckResourceAttr("buddy_target.test", "hostname", newHostname),
					resource.TestCheckResourceAttr("buddy_target.test", "port", strconv.Itoa(newPort)),
					resource.TestCheckResourceAttr("buddy_target.test", "username", newUsername),
					resource.TestCheckResourceAttr("buddy_target.test", "description", newDescription),
					resource.TestCheckResourceAttr("buddy_target.test", "file_path", newFilePath),
					resource.TestCheckResourceAttr("buddy_target.test", "auth_mode", "KEY"),
					resource.TestCheckResourceAttr("buddy_target.test", "all_pipelines_allowed", "false"),
					resource.TestCheckResourceAttr("buddy_target.test", "tags.#", "2"),
				),
			},
			// Import
			{
				ResourceName:            "buddy_target.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "passphrase"},
			},
		},
	})
}

func TestAccTarget_project_ftp(t *testing.T) {
	domain := util.UniqueString()
	projectName := util.RandString(10)
	name := util.RandString(10)
	hostname := "ftp.example.com"
	username := "ftpuser"
	port := 21
	
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create project
			{
				Config: testAccProjectConfig(domain, projectName),
			},
			// Create FTP target
			{
				Config: testAccProjectConfig(domain, projectName) + testAccTargetConfigProjectFtp(domain, projectName, name, hostname, port, username),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetCheckExist("buddy_target.test", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "domain", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "project_name", projectName),
					resource.TestCheckResourceAttr("buddy_target.test", "name", name),
					resource.TestCheckResourceAttr("buddy_target.test", "type", "FTP"),
					resource.TestCheckResourceAttr("buddy_target.test", "hostname", hostname),
					resource.TestCheckResourceAttr("buddy_target.test", "port", strconv.Itoa(port)),
					resource.TestCheckResourceAttr("buddy_target.test", "username", username),
				),
			},
			// Import
			{
				ResourceName:            "buddy_target.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func TestAccTarget_sftp(t *testing.T) {
	domain := util.UniqueString()
	name := util.RandString(10)
	hostname := "sftp.example.com"
	username := "sftpuser"
	port := 22
	
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create SFTP target
			{
				Config: testAccTargetConfigSftp(domain, name, hostname, port, username),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetCheckExist("buddy_target.test", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "domain", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "name", name),
					resource.TestCheckResourceAttr("buddy_target.test", "type", "SFTP"),
					resource.TestCheckResourceAttr("buddy_target.test", "hostname", hostname),
					resource.TestCheckResourceAttr("buddy_target.test", "port", strconv.Itoa(port)),
					resource.TestCheckResourceAttr("buddy_target.test", "username", username),
				),
			},
		},
	})
}

func TestAccTarget_ftps(t *testing.T) {
	domain := util.UniqueString()
	name := util.RandString(10)
	hostname := "ftps.example.com"
	username := "ftpsuser"
	port := 990
	
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create FTPS target
			{
				Config: testAccTargetConfigFtps(domain, name, hostname, port, username),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetCheckExist("buddy_target.test", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "domain", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "name", name),
					resource.TestCheckResourceAttr("buddy_target.test", "type", "FTPS"),
					resource.TestCheckResourceAttr("buddy_target.test", "hostname", hostname),
					resource.TestCheckResourceAttr("buddy_target.test", "port", strconv.Itoa(port)),
					resource.TestCheckResourceAttr("buddy_target.test", "username", username),
				),
			},
		},
	})
}

func TestAccTarget_shell(t *testing.T) {
	domain := util.UniqueString()
	name := util.RandString(10)
	
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create Shell target
			{
				Config: testAccTargetConfigShell(domain, name),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetCheckExist("buddy_target.test", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "domain", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "name", name),
					resource.TestCheckResourceAttr("buddy_target.test", "type", "SHELL"),
				),
			},
		},
	})
}

func TestAccTarget_s3(t *testing.T) {
	domain := util.UniqueString()
	name := util.RandString(10)
	
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create S3 target
			{
				Config: testAccTargetConfigS3(domain, name),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetCheckExist("buddy_target.test", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "domain", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "name", name),
					resource.TestCheckResourceAttr("buddy_target.test", "type", "AMAZON_S3"),
				),
			},
		},
	})
}

func TestAccTarget_gcs(t *testing.T) {
	domain := util.UniqueString()
	name := util.RandString(10)
	
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create GCS target
			{
				Config: testAccTargetConfigGcs(domain, name),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetCheckExist("buddy_target.test", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "domain", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "name", name),
					resource.TestCheckResourceAttr("buddy_target.test", "type", "GOOGLE_CLOUD_STORAGE"),
				),
			},
		},
	})
}

func TestAccTarget_azure(t *testing.T) {
	domain := util.UniqueString()
	name := util.RandString(10)
	
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create Azure Storage target
			{
				Config: testAccTargetConfigAzure(domain, name),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetCheckExist("buddy_target.test", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "domain", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "name", name),
					resource.TestCheckResourceAttr("buddy_target.test", "type", "AZURE_STORAGE"),
				),
			},
		},
	})
}

func TestAccTarget_docker(t *testing.T) {
	domain := util.UniqueString()
	name := util.RandString(10)
	hostname := "registry.example.com"
	username := "dockeruser"
	
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create Docker Registry target
			{
				Config: testAccTargetConfigDocker(domain, name, hostname, username),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetCheckExist("buddy_target.test", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "domain", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "name", name),
					resource.TestCheckResourceAttr("buddy_target.test", "type", "DOCKER_REGISTRY"),
					resource.TestCheckResourceAttr("buddy_target.test", "hostname", hostname),
					resource.TestCheckResourceAttr("buddy_target.test", "username", username),
				),
			},
		},
	})
}

func TestAccTarget_kubernetes(t *testing.T) {
	domain := util.UniqueString()
	name := util.RandString(10)
	hostname := "k8s.example.com"
	
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create Kubernetes target
			{
				Config: testAccTargetConfigKubernetes(domain, name, hostname),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetCheckExist("buddy_target.test", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "domain", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "name", name),
					resource.TestCheckResourceAttr("buddy_target.test", "type", "KUBERNETES"),
					resource.TestCheckResourceAttr("buddy_target.test", "hostname", hostname),
				),
			},
		},
	})
}

func TestAccTarget_allowed_pipelines(t *testing.T) {
	domain := util.UniqueString()
	projectName := util.RandString(10)
	name := util.RandString(10)
	hostname := "example.com"
	username := "testuser"
	port := 22
	
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccTargetCheckDestroy,
		Steps: []resource.TestStep{
			// Create project with pipeline
			{
				Config: testAccProjectConfigWithPipeline(domain, projectName),
			},
			// Create target with allowed pipelines
			{
				Config: testAccProjectConfigWithPipeline(domain, projectName) + testAccTargetConfigWithAllowedPipelines(domain, projectName, name, hostname, port, username),
				Check: resource.ComposeTestCheckFunc(
					testAccTargetCheckExist("buddy_target.test", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "domain", domain),
					resource.TestCheckResourceAttr("buddy_target.test", "project_name", projectName),
					resource.TestCheckResourceAttr("buddy_target.test", "name", name),
					resource.TestCheckResourceAttr("buddy_target.test", "all_pipelines_allowed", "false"),
					resource.TestCheckResourceAttr("buddy_target.test", "allowed_pipelines.#", "1"),
				),
			},
		},
	})
}

func testAccTargetCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_target" {
			continue
		}
		domain, pid, tid, err := util.DecomposeTripleId(rs.Primary.ID)
		if err != nil {
			domain, tid, err = util.DecomposeDoubleId(rs.Primary.ID)
			if err != nil {
				return err
			}
		}
		targetId, err := strconv.Atoi(tid)
		if err != nil {
			return err
		}
		var target *buddy.Target
		var resp *buddy.Response
		if pid != "" {
			target, resp, err = acc.ApiClient.TargetService.GetInProject(domain, pid, targetId)
		} else {
			target, resp, err = acc.ApiClient.TargetService.GetInWorkspace(domain, targetId)
		}
		if err == nil && target != nil {
			return util.ErrorResourceExists()
		}
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}

func testAccTargetCheckExist(n string, domain string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		dom, pid, tid, err := util.DecomposeTripleId(rs.Primary.ID)
		if err != nil {
			dom, tid, err = util.DecomposeDoubleId(rs.Primary.ID)
			if err != nil {
				return err
			}
		}
		if dom != domain {
			return fmt.Errorf("bad domain, expected: %s, got: %s", domain, dom)
		}
		targetId, err := strconv.Atoi(tid)
		if err != nil {
			return err
		}
		var target *buddy.Target
		if pid != "" {
			target, _, err = acc.ApiClient.TargetService.GetInProject(dom, pid, targetId)
		} else {
			target, _, err = acc.ApiClient.TargetService.GetInWorkspace(dom, targetId)
		}
		if err != nil {
			return err
		}
		if target == nil {
			return fmt.Errorf("target not found")
		}
		return nil
	}
}

func testAccTargetConfigWorkspaceSsh(domain string, name string, hostname string, port int, username string) string {
	return fmt.Sprintf(`
resource "buddy_target" "test" {
	domain = "%s"
	name = "%s"
	type = "SSH"
	hostname = "%s"
	port = %d
	username = "%s"
	password = "test123"
}`, domain, name, hostname, port, username)
}

func testAccTargetConfigWorkspaceSshFull(domain string, name string, hostname string, port int, username string, description string, filePath string) string {
	return fmt.Sprintf(`
resource "buddy_variable_ssh_key" "test" {
	domain = "%s"
	display_name = "test-key"
	value = "test-key-value"
	file_place = "NONE"
}

resource "buddy_target" "test" {
	domain = "%s"
	name = "%s"
	type = "SSH"
	hostname = "%s"
	port = %d
	username = "%s"
	description = "%s"
	file_path = "%s"
	auth_mode = "KEY"
	key_id = buddy_variable_ssh_key.test.variable_id
	all_pipelines_allowed = false
	tags = ["production", "secure"]
}`, domain, domain, name, hostname, port, username, description, filePath)
}

func testAccTargetConfigProjectFtp(domain string, projectName string, name string, hostname string, port int, username string) string {
	return fmt.Sprintf(`
resource "buddy_target" "test" {
	domain = "%s"
	project_name = "%s"
	name = "%s"
	type = "FTP"
	hostname = "%s"
	port = %d
	username = "%s"
	password = "ftp123"
}`, domain, projectName, name, hostname, port, username)
}

func testAccTargetConfigSftp(domain string, name string, hostname string, port int, username string) string {
	return fmt.Sprintf(`
resource "buddy_target" "test" {
	domain = "%s"
	name = "%s"
	type = "SFTP"
	hostname = "%s"
	port = %d
	username = "%s"
	password = "sftp123"
}`, domain, name, hostname, port, username)
}

func testAccTargetConfigFtps(domain string, name string, hostname string, port int, username string) string {
	return fmt.Sprintf(`
resource "buddy_target" "test" {
	domain = "%s"
	name = "%s"
	type = "FTPS"
	hostname = "%s"
	port = %d
	username = "%s"
	password = "ftps123"
}`, domain, name, hostname, port, username)
}

func testAccTargetConfigShell(domain string, name string) string {
	return fmt.Sprintf(`
resource "buddy_target" "test" {
	domain = "%s"
	name = "%s"
	type = "SHELL"
}`, domain, name)
}

func testAccTargetConfigS3(domain string, name string) string {
	return fmt.Sprintf(`
resource "buddy_target" "test" {
	domain = "%s"
	name = "%s"
	type = "AMAZON_S3"
}`, domain, name)
}

func testAccTargetConfigGcs(domain string, name string) string {
	return fmt.Sprintf(`
resource "buddy_target" "test" {
	domain = "%s"
	name = "%s"
	type = "GOOGLE_CLOUD_STORAGE"
}`, domain, name)
}

func testAccTargetConfigAzure(domain string, name string) string {
	return fmt.Sprintf(`
resource "buddy_target" "test" {
	domain = "%s"
	name = "%s"
	type = "AZURE_STORAGE"
}`, domain, name)
}

func testAccTargetConfigDocker(domain string, name string, hostname string, username string) string {
	return fmt.Sprintf(`
resource "buddy_target" "test" {
	domain = "%s"
	name = "%s"
	type = "DOCKER_REGISTRY"
	hostname = "%s"
	username = "%s"
	password = "docker123"
}`, domain, name, hostname, username)
}

func testAccTargetConfigKubernetes(domain string, name string, hostname string) string {
	return fmt.Sprintf(`
resource "buddy_target" "test" {
	domain = "%s"
	name = "%s"
	type = "KUBERNETES"
	hostname = "%s"
}`, domain, name, hostname)
}

func testAccTargetConfigWithAllowedPipelines(domain string, projectName string, name string, hostname string, port int, username string) string {
	return fmt.Sprintf(`
resource "buddy_target" "test" {
	domain = "%s"
	project_name = "%s"
	name = "%s"
	type = "SSH"
	hostname = "%s"
	port = %d
	username = "%s"
	password = "test123"
	all_pipelines_allowed = false
	allowed_pipelines = [buddy_pipeline.test.pipeline_id]
}`, domain, projectName, name, hostname, port, username)
}

func testAccProjectConfig(domain string, projectName string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
	domain = "%s"
}

resource "buddy_project" "test" {
	domain = buddy_workspace.foo.domain
	display_name = "%s"
}`, domain, projectName)
}

func testAccProjectConfigWithPipeline(domain string, projectName string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
	domain = "%s"
}

resource "buddy_project" "test" {
	domain = buddy_workspace.foo.domain
	display_name = "%s"
}

resource "buddy_pipeline" "test" {
	domain = buddy_workspace.foo.domain
	project_name = buddy_project.test.name
	name = "test"
}`, domain, projectName)
}