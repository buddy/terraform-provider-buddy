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

var ignoreEnvImportVerify = []string{
	"permissions",
	"allowed_pipeline",
	"allowed_environment",
}

func TestAccEnvironmentPermissions(t *testing.T) {
	var environment buddy.Environment
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name := util.RandString(10)
	email := util.RandEmail()
	groupName := util.RandString(10)
	identifier := util.UniqueString()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccEnvironmentCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentPermissionsUserConfig(domain, projectName, identifier, name, email, groupName, buddy.EnvironmentPermissionAccessLevelDenied, buddy.EnvironmentPermissionAccessLevelUseOnly),
				Check: resource.ComposeTestCheckFunc(
					testAccEnvironmentGet("buddy_environment.env", &environment),
					testAccEnvironmentAttributes("buddy_environment.env", &environment, name, identifier, "", "", true, true, buddy.EnvironmentScopeProject, false, "", "", buddy.EnvironmentPermissionAccessLevelDenied, buddy.EnvironmentPermissionAccessLevelUseOnly, ""),
				),
			},
			{
				Config: testAccEnvironmentPermissionsConfig(domain, projectName, identifier, name, email, groupName, buddy.EnvironmentPermissionAccessLevelDefault),
				Check: resource.ComposeTestCheckFunc(
					testAccEnvironmentGet("buddy_environment.env", &environment),
					testAccEnvironmentAttributes("buddy_environment.env", &environment, name, identifier, "", "", true, true, buddy.EnvironmentScopeProject, false, "", "", buddy.EnvironmentPermissionAccessLevelDefault, "", ""),
				),
			},
			{
				Config: testAccEnvironmentPermissionsUserGroupConfig(domain, projectName, identifier, name, email, groupName, buddy.EnvironmentPermissionAccessLevelManage, buddy.EnvironmentPermissionAccessLevelDefault, buddy.EnvironmentPermissionAccessLevelDenied),
				Check: resource.ComposeTestCheckFunc(
					testAccEnvironmentGet("buddy_environment.env", &environment),
					testAccEnvironmentAttributes("buddy_environment.env", &environment, name, identifier, "", "", true, true, buddy.EnvironmentScopeProject, false, "", "", buddy.EnvironmentPermissionAccessLevelManage, buddy.EnvironmentPermissionAccessLevelDefault, buddy.EnvironmentPermissionAccessLevelDenied),
				),
			},
			{
				Config: testAccEnvironmentPermissionsGroupConfig(domain, projectName, identifier, name, email, groupName, buddy.EnvironmentPermissionAccessLevelUseOnly, buddy.EnvironmentPermissionAccessLevelManage),
				Check: resource.ComposeTestCheckFunc(
					testAccEnvironmentGet("buddy_environment.env", &environment),
					testAccEnvironmentAttributes("buddy_environment.env", &environment, name, identifier, "", "", true, true, buddy.EnvironmentScopeProject, false, "", "", buddy.EnvironmentPermissionAccessLevelUseOnly, "", buddy.EnvironmentPermissionAccessLevelManage),
				),
			},
			// import env
			{
				ResourceName:            "buddy_environment.env",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoreEnvImportVerify,
			},
		},
	})
}

func TestAccEnvironmentWorkspace(t *testing.T) {
	var environment buddy.Environment
	domain := util.UniqueString()
	baseName := util.RandString(10)
	baseIdentifier := util.UniqueString()
	envProjName := util.RandString(10)
	envProjIdentifier := util.UniqueString()
	pipName := util.RandString(10)
	pipIdentifier := util.UniqueString()
	name := util.RandString(10)
	identifier := util.UniqueString()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccEnvironmentCheckDestroy,
		Steps: []resource.TestStep{
			// create env
			{
				Config: testAccEnvironmentWorkspace(domain, baseName, baseIdentifier, envProjName, envProjIdentifier, pipName, pipIdentifier, name, identifier),
				Check: resource.ComposeTestCheckFunc(
					testAccEnvironmentGet("buddy_environment.env", &environment),
					testAccEnvironmentAttributes("buddy_environment.env", &environment, name, identifier, "", "", false, false, buddy.EnvironmentScopeWorkspace, false, baseIdentifier, "", "", "", ""),
				),
			},
			// import env
			{
				ResourceName:            "buddy_environment.env",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoreEnvImportVerify,
			},
		},
	})
}

func TestAccEnvironmentSimple(t *testing.T) {
	var environment buddy.Environment
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	identifier := util.UniqueString()
	newIdentifier := util.UniqueString()
	url := "https://" + util.RandString(10) + ".com"
	newUrl := "https://" + util.RandString(10) + ".com"
	tag := util.RandString(3)
	newTag := util.RandString(3)
	icon := util.RandString(5)
	newIcon := util.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccEnvironmentCheckDestroy,
		Steps: []resource.TestStep{
			// create env
			{
				Config: testAccEnvironmentConfig(domain, projectName, name, identifier, url, false, false, true, icon, tag),
				Check: resource.ComposeTestCheckFunc(
					testAccEnvironmentGet("buddy_environment.env", &environment),
					testAccEnvironmentAttributes("buddy_environment.env", &environment, name, identifier, url, icon, false, false, buddy.EnvironmentScopeProject, true, "", tag, "", "", ""),
				),
			},
			// update env
			{
				Config: testAccEnvironmentConfig(domain, projectName, newName, newIdentifier, newUrl, true, true, false, newIcon, newTag),
				Check: resource.ComposeTestCheckFunc(
					testAccEnvironmentGet("buddy_environment.env", &environment),
					testAccEnvironmentAttributes("buddy_environment.env", &environment, newName, newIdentifier, newUrl, newIcon, true, true, buddy.EnvironmentScopeProject, false, "", newTag, "", "", ""),
				),
			},
			// import env
			{
				ResourceName:            "buddy_environment.env",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoreEnvImportVerify,
			},
		},
	})
}

func testAccEnvironmentAttributes(n string, environment *buddy.Environment, name string, identifier string, url string, icon string, allPipelinesAllowed bool, allEnvsAllowed bool, scope string, baseOnly bool, baseEnvironment string, tag string, othersLevel string, userLevel string, groupLevel string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsAllPipelinesAllowed, _ := strconv.ParseBool(attrs["all_pipelines_allowed"])
		attrsAllEnvsAllowed, _ := strconv.ParseBool(attrs["all_environments_allowed"])
		attrsBaseOnly, _ := strconv.ParseBool(attrs["base_only"])
		if err := util.CheckFieldEqualAndSet("Name", environment.Name, name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], name); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("AllPipelinesAllowed", environment.AllPipelinesAllowed, allPipelinesAllowed); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("all_pipelines_allowed", attrsAllPipelinesAllowed, allPipelinesAllowed); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("AllEnvironmentsAllowed", environment.AllEnvironmentsAllowed, allEnvsAllowed); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("all_environments_allowed", attrsAllEnvsAllowed, allEnvsAllowed); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("BaseOnly", environment.BaseOnly, baseOnly); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("base_only", attrsBaseOnly, baseOnly); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("Scope", environment.Scope, scope); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("scope", attrs["scope"], scope); err != nil {
			return err
		}
		if scope == buddy.EnvironmentScopeProject {
			if err := util.CheckFieldEqualAndSet("Project.Name", environment.Project.Name, attrs["project_name"]); err != nil {
				return err
			}
		}
		if baseEnvironment != "" {
			if err := util.CheckFieldEqualAndSet("BaseEnvironments[0]", environment.BaseEnvironments[0], baseEnvironment); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("base_environments.0", attrs["base_environments.0"], baseEnvironment); err != nil {
				return err
			}
		}
		if icon != "" {
			if err := util.CheckFieldEqualAndSet("Icon", environment.Icon, icon); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("icon", attrs["icon"], icon); err != nil {
				return err
			}
		}
		if identifier != "" {
			if err := util.CheckFieldEqualAndSet("Identifier", environment.Identifier, identifier); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("identifier", attrs["identifier"], identifier); err != nil {
				return err
			}
		} else {
			if err := util.CheckFieldEqualAndSet("identifier", attrs["identifier"], environment.Identifier); err != nil {
				return err
			}
		}
		if url != "" {
			if err := util.CheckFieldEqualAndSet("PublicUrl", environment.PublicUrl, url); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("public_url", attrs["public_url"], url); err != nil {
				return err
			}
		}
		if err := util.CheckFieldEqualAndSet("environment_id", attrs["environment_id"], environment.Id); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("html_url", attrs["html_url"], environment.HtmlUrl); err != nil {
			return err
		}
		if tag != "" {
			if err := util.CheckIntFieldEqual("Tags", len(environment.Tags), 1); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("Tags[0]", environment.Tags[0], tag); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("tags.0", attrs["tags.0"], tag); err != nil {
				return err
			}
		} else {
			if err := util.CheckIntFieldEqual("Tags", len(environment.Tags), 0); err != nil {
				return err
			}
		}

		if othersLevel != "" {
			if err := util.CheckFieldEqualAndSet("Permissions.Others", environment.Permissions.Others, othersLevel); err != nil {
				return err
			}
		} else {
			defaultOtherAccess := buddy.EnvironmentPermissionAccessLevelDefault
			if scope == buddy.EnvironmentScopeWorkspace {
				defaultOtherAccess = buddy.EnvironmentPermissionAccessLevelUseOnly
			}
			if err := util.CheckFieldEqualAndSet("Permissions.Others", environment.Permissions.Others, defaultOtherAccess); err != nil {
				return err
			}
		}
		if userLevel != "" {
			if err := util.CheckFieldEqualAndSet("Permissions.Users[0].AccessLevel", environment.Permissions.Users[0].AccessLevel, userLevel); err != nil {
				return err
			}
		} else {
			count := 0
			if scope == buddy.EnvironmentScopeWorkspace {
				count = 1
			}
			if err := util.CheckIntFieldEqual("len(Permissions.Users)", len(environment.Permissions.Users), count); err != nil {
				return err
			}
		}
		if groupLevel != "" {
			if err := util.CheckFieldEqualAndSet("Permissions.Groups[0].AccessLevel", environment.Permissions.Groups[0].AccessLevel, groupLevel); err != nil {
				return err
			}
		} else {
			if err := util.CheckIntFieldEqual("len(Permissions.Groups)", len(environment.Permissions.Groups), 0); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccEnvironmentGet(n string, environment *buddy.Environment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		domain, _, environmentId, err := util.DecomposeTripleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		e, _, err := acc.ApiClient.EnvironmentService.Get(domain, environmentId)
		if err != nil {
			return err
		}
		*environment = *e
		return nil
	}
}

func testAccEnvironmentPermissionsUserGroupConfig(domain string, projectName string, identifier string, name string, email string, groupName string, otherLevel string, userLevel string, groupLevel string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_member" "a" {
    domain = "${buddy_workspace.foo.domain}"
    email = "%s"
}

resource "buddy_group" "g" {
	domain = "${buddy_workspace.foo.domain}"
	name = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_permission" "a" {
    domain = "${buddy_workspace.foo.domain}"
    name = "perm"
    pipeline_access_level = "READ_WRITE"
    repository_access_level = "READ_ONLY"
	  sandbox_access_level = "READ_ONLY"
}

resource "buddy_project_member" "bar" {
	domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
	member_id = "${buddy_member.a.member_id}"
	permission_id = "${buddy_permission.a.permission_id}"
}

resource "buddy_project_group" "bar" {
	domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
	group_id = "${buddy_group.g.group_id}"
	permission_id = "${buddy_permission.a.permission_id}"
}

resource "buddy_environment" "env" {
    domain = "${buddy_workspace.foo.domain}"
    project_name = "${buddy_project.proj.name}"
    identifier = "%s"
    name = "%s"
    permissions {
      others = "%s"
      user {
        id = "${buddy_project_member.bar.member_id}"
        access_level = "%s"
      }
      group {
		  	id = "${buddy_project_group.bar.group_id}"
			  access_level = "%s"
		  }
    }
}
`, domain, email, groupName, projectName, identifier, name, otherLevel, userLevel, groupLevel)
}

func testAccEnvironmentPermissionsConfig(domain string, projectName string, identifier string, name string, email string, groupName string, otherLevel string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_member" "a" {
    domain = "${buddy_workspace.foo.domain}"
    email = "%s"
}

resource "buddy_group" "g" {
	domain = "${buddy_workspace.foo.domain}"
	name = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_permission" "a" {
    domain = "${buddy_workspace.foo.domain}"
    name = "perm"
    pipeline_access_level = "READ_WRITE"
    repository_access_level = "READ_ONLY"
	  sandbox_access_level = "READ_ONLY"
}

resource "buddy_project_member" "bar" {
	domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
	member_id = "${buddy_member.a.member_id}"
	permission_id = "${buddy_permission.a.permission_id}"
}

resource "buddy_project_group" "bar" {
	domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
	group_id = "${buddy_group.g.group_id}"
	permission_id = "${buddy_permission.a.permission_id}"
}

resource "buddy_environment" "env" {
    domain = "${buddy_workspace.foo.domain}"
    project_name = "${buddy_project.proj.name}"
    identifier = "%s"
    name = "%s"
    permissions {
      others = "%s"
    }
}
`, domain, email, groupName, projectName, identifier, name, otherLevel)
}

func testAccEnvironmentPermissionsGroupConfig(domain string, projectName string, identifier string, name string, email string, groupName string, otherLevel string, groupLevel string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_member" "a" {
    domain = "${buddy_workspace.foo.domain}"
    email = "%s"
}

resource "buddy_group" "g" {
	domain = "${buddy_workspace.foo.domain}"
	name = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_permission" "a" {
    domain = "${buddy_workspace.foo.domain}"
    name = "perm"
    pipeline_access_level = "READ_WRITE"
    repository_access_level = "READ_ONLY"
	  sandbox_access_level = "READ_ONLY"
}

resource "buddy_project_member" "bar" {
	domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
	member_id = "${buddy_member.a.member_id}"
	permission_id = "${buddy_permission.a.permission_id}"
}

resource "buddy_project_group" "bar" {
	domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
	group_id = "${buddy_group.g.group_id}"
	permission_id = "${buddy_permission.a.permission_id}"
}

resource "buddy_environment" "env" {
    domain = "${buddy_workspace.foo.domain}"
    project_name = "${buddy_project.proj.name}"
    identifier = "%s"
    name = "%s"
    permissions {
      others = "%s"
      group {
		  	id = "${buddy_project_group.bar.group_id}"
			  access_level = "%s"
		  }
    }
}
`, domain, email, groupName, projectName, identifier, name, otherLevel, groupLevel)
}

func testAccEnvironmentPermissionsUserConfig(domain string, projectName string, identifier string, name string, email string, groupName string, otherLevel string, userLevel string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_member" "a" {
    domain = "${buddy_workspace.foo.domain}"
    email = "%s"
}

resource "buddy_group" "g" {
	domain = "${buddy_workspace.foo.domain}"
	name = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_permission" "a" {
    domain = "${buddy_workspace.foo.domain}"
    name = "perm"
    pipeline_access_level = "READ_WRITE"
    repository_access_level = "READ_ONLY"
	  sandbox_access_level = "READ_ONLY"
}

resource "buddy_project_member" "bar" {
	domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
	member_id = "${buddy_member.a.member_id}"
	permission_id = "${buddy_permission.a.permission_id}"
}

resource "buddy_project_group" "bar" {
	domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
	group_id = "${buddy_group.g.group_id}"
	permission_id = "${buddy_permission.a.permission_id}"
}

resource "buddy_environment" "env" {
    domain = "${buddy_workspace.foo.domain}"
    project_name = "${buddy_project.proj.name}"
    identifier = "%s"
    name = "%s"
    permissions {
      others = "%s"
      user {
        id = "${buddy_project_member.bar.member_id}"
        access_level = "%s"
      }
    }
}
`, domain, email, groupName, projectName, identifier, name, otherLevel, userLevel)
}

func testAccEnvironmentWorkspace(domain string, baseName string, baseIdentifier string, projEnvName string, projEnvIdentifier string, pipName string, pipIdentifier string, name string, identifier string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_environment" "base" {
		domain = "${buddy_workspace.foo.domain}"
    name = "%s"
    identifier = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "abc"
}

resource "buddy_environment" "env_proj" {
	domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
	name = "%s"
	identifier = "%s"
}

resource "buddy_pipeline" "pip_proj" {
  domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
  name = "%s"
  identifier = "%s"
}

resource "buddy_environment" "env" {
    domain = "${buddy_workspace.foo.domain}"
    name = "%s"
    identifier = "%s"
		base_environments = ["${buddy_environment.base.identifier}"]
		allowed_pipeline {
			project = "${buddy_project.proj.name}"
      pipeline = "${buddy_pipeline.pip_proj.identifier}"
		}
		allowed_environment {
			project = "${buddy_project.proj.name}"
      environment = "${buddy_environment.env_proj.identifier}"
		}
}
`, domain, baseName, baseIdentifier, projEnvName, projEnvIdentifier, pipName, pipIdentifier, name, identifier)
}

func testAccEnvironmentConfig(domain string, projectName string, name string, identifier string, url string, allPipelinesAllowed bool, allEnvsAllowed bool, baseOnly bool, icon string, tag string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_environment" "env" {
    domain = "${buddy_workspace.foo.domain}"
    project_name = "${buddy_project.proj.name}"
    name = "%s"
    identifier = "%s"
    public_url = "%s"
    all_pipelines_allowed = "%t"
		all_environments_allowed = "%t"
		base_only = "%t"
		icon = "%s"
    tags = ["%s"]
}
`, domain, projectName, name, identifier, url, allPipelinesAllowed, allEnvsAllowed, baseOnly, icon, tag)
}

func testAccEnvironmentCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_environment" {
			continue
		}
		domain, _, environmentId, err := util.DecomposeTripleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		environment, resp, err := acc.ApiClient.EnvironmentService.Get(domain, environmentId)
		if err == nil && environment != nil {
			return util.ErrorResourceExists()
		}
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}
