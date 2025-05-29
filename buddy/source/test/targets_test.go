package test

import (
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"
)

func TestAccSourceTargets_all(t *testing.T) {
	domain := util.UniqueString()
	name1 := util.RandString(10)
	identifier1 := util.UniqueString()
	name2 := util.RandString(10)
	identifier2 := util.UniqueString()
	host := "1.1.1.1"
	port := "44"
	username := util.RandString(10)
	password := util.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceTargetsAllConfig(domain, name1, identifier1, name2, identifier2, host, port, username, password),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.buddy_targets.all", "targets.#", "2"),
				),
			},
		},
	})
}

func TestAccSourceTargets_byProject(t *testing.T) {
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name1 := util.RandString(10)
	identifier1 := util.UniqueString()
	name2 := util.RandString(10)
	identifier2 := util.UniqueString()
	name3 := util.RandString(10)
	identifier3 := util.UniqueString()
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
		Steps: []resource.TestStep{
			{
				Config: testAccSourceTargetsByProjectConfig(domain, projectName, name1, identifier1, name2, identifier2, name3, identifier3, host, port, path, username, key),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.buddy_targets.project", "targets.#", "2"),
					resource.TestCheckResourceAttr("data.buddy_targets.project", "targets.0.name", name1),
					resource.TestCheckResourceAttr("data.buddy_targets.project", "targets.1.name", name2),
				),
			},
		},
	})
}

func TestAccSourceTargets_byPipeline(t *testing.T) {
	domain := util.UniqueString()
	projectName := util.UniqueString()
	pipelineName := util.UniqueString()
	name1 := util.RandString(10)
	identifier1 := util.UniqueString()
	name2 := util.RandString(10)
	identifier2 := util.UniqueString()
	host := "1.1.1.1"
	port := "44"
	path := util.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceTargetsByPipelineConfig(domain, projectName, pipelineName, name1, identifier1, name2, identifier2, host, port, path),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.buddy_targets.pipeline", "targets.#", "1"),
					resource.TestCheckResourceAttr("data.buddy_targets.pipeline", "targets.0.name", name1),
				),
			},
		},
	})
}

func TestAccSourceTargets_byEnvironment(t *testing.T) {
	domain := util.UniqueString()
	projectName := util.UniqueString()
	envName := util.UniqueString()
	envId := util.UniqueString()
	name1 := util.RandString(10)
	identifier1 := util.UniqueString()
	name2 := util.RandString(10)
	identifier2 := util.UniqueString()
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
		Steps: []resource.TestStep{
			{
				Config: testAccSourceTargetsByEnvironmentConfig(domain, projectName, envName, envId, name1, identifier1, name2, identifier2, host, port, path, username, asset),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.buddy_targets.environment", "targets.#", "1"),
					resource.TestCheckResourceAttr("data.buddy_targets.environment", "targets.0.name", name1),
				),
			},
		},
	})
}

func testAccSourceTargetsAllConfig(domain string, name1 string, identifier1 string, name2 string, identifier2 string, host string, port string, username string, password string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_target" "test1" {
    domain     = buddy_workspace.test.domain
    name       = "%s"
    identifier = "%s"
    type       = "%s"
    host       = "%s"
    port       = "%s"
    auth {
        method   = "%s"
        username = "%s"
        password = "%s"
    }
}

resource "buddy_target" "test2" {
    domain     = buddy_workspace.test.domain
    name       = "%s"
    identifier = "%s"
    type       = "%s"
    host       = "%s"
    port       = "%s"
    auth {
        method   = "%s"
        username = "%s"
        password = "%s"
    }
}

data "buddy_targets" "all" {
    domain = buddy_workspace.test.domain
    depends_on = [buddy_target.test1, buddy_target.test2]
}`, domain, name1, identifier1, buddy.TargetTypeSsh, host, port, buddy.TargetAuthMethodPassword, username, password, name2, identifier2, buddy.TargetTypeSsh, host, port, buddy.TargetAuthMethodPassword, username, password)
}

func testAccSourceTargetsByProjectConfig(domain string, projectName string, name1 string, identifier1 string, name2 string, identifier2 string, name3 string, identifier3 string, host string, port string, path string, username string, key string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_project" "test" {
    domain       = buddy_workspace.test.domain
    display_name = "%s"
}

resource "buddy_target" "project1" {
    domain       = buddy_workspace.test.domain
    project_name = buddy_project.test.name
    name         = "%s"
    identifier   = "%s"
    type         = "%s"
    host         = "%s"
    port         = "%s"
    path         = "%s"
    auth {
        method   = "%s"
        username = "%s"
        key      = "%s"
    }
}

resource "buddy_target" "project2" {
    domain       = buddy_workspace.test.domain
    project_name = buddy_project.test.name
    name         = "%s"
    identifier   = "%s"
    type         = "%s"
    host         = "%s"
    port         = "%s"
    path         = "%s"
    auth {
        method   = "%s"
        username = "%s"
        key      = "%s"
    }
}

resource "buddy_target" "workspace" {
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
}

data "buddy_targets" "project" {
    domain       = buddy_workspace.test.domain
    project_name = buddy_project.test.name
    depends_on   = [buddy_target.project1, buddy_target.project2, buddy_target.workspace]
}`, domain, projectName, name1, identifier1, buddy.TargetTypeSsh, host, port, path, buddy.TargetAuthMethodSshKey, username, key, name2, identifier2, buddy.TargetTypeSsh, host, port, path, buddy.TargetAuthMethodSshKey, username, key, name3, identifier3, buddy.TargetTypeSsh, host, port, path, buddy.TargetAuthMethodSshKey, username, key)
}

func testAccSourceTargetsByPipelineConfig(domain string, projectName string, pipelineName string, name1 string, identifier1 string, name2 string, identifier2 string, host string, port string, path string) string {
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

resource "buddy_target" "pipeline" {
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
}

resource "buddy_target" "workspace" {
    domain     = buddy_workspace.test.domain
    name       = "%s"
    identifier = "%s"
    type       = "%s"
    host       = "%s"
    port       = "%s"
    path       = "%s"
    auth {
        method = "%s"
    }
}

data "buddy_targets" "pipeline" {
    domain      = buddy_workspace.test.domain
    pipeline_id = buddy_pipeline.test.pipeline_id
    depends_on  = [buddy_target.pipeline, buddy_target.workspace]
}`, domain, projectName, pipelineName, name1, identifier1, buddy.TargetTypeSsh, host, port, path, buddy.TargetAuthMethodProxyCredentials, name2, identifier2, buddy.TargetTypeSsh, host, port, path, buddy.TargetAuthMethodProxyCredentials)
}

func testAccSourceTargetsByEnvironmentConfig(domain string, projectName string, envName string, envId string, name1 string, identifier1 string, name2 string, identifier2 string, host string, port string, path string, username string, asset string) string {
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

resource "buddy_target" "environment" {
    domain         = buddy_workspace.test.domain
    environment_id = buddy_environment.test.environment_id
    name           = "%s"
    identifier     = "%s"
    type           = "%s"
    host           = "%s"
    port           = "%s"
    path           = "%s"
    auth {
        method   = "%s"
        username = "%s"
        asset    = "%s"
    }
}

resource "buddy_target" "workspace" {
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
        asset    = "%s"
    }
}

data "buddy_targets" "environment" {
    domain         = buddy_workspace.test.domain
    environment_id = buddy_environment.test.environment_id
    depends_on     = [buddy_target.environment, buddy_target.workspace]
}`, domain, projectName, envName, envId, buddy.EnvironmentTypeDev, name1, identifier1, buddy.TargetTypeSsh, host, port, path, buddy.TargetAuthMethodAssetsKey, username, asset, name2, identifier2, buddy.TargetTypeSsh, host, port, path, buddy.TargetAuthMethodAssetsKey, username, asset)
}
