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
	"var",
}

func TestAccEnvironmentVariables(t *testing.T) {
	var environment buddy.Environment
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name := util.RandString(10)
	identifier := util.UniqueString()
	key := util.UniqueString()
	val := util.RandString(10)
	newVal := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccEnvironmentCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEnvironmentVariablesConfig(domain, projectName, identifier, name, key, val),
				Check: resource.ComposeTestCheckFunc(
					testAccEnvironmentGet("buddy_environment.env", &environment),
					testAccEnvironmentAttributes("buddy_environment.env", &environment, name, "", "", true, "", "", "", "", key, val),
				),
			},
			{
				Config: testAccEnvironmentVariablesConfig(domain, projectName, identifier, name, key, newVal),
				Check: resource.ComposeTestCheckFunc(
					testAccEnvironmentGet("buddy_environment.env", &environment),
					testAccEnvironmentAttributes("buddy_environment.env", &environment, name, "", "", true, "", "", "", "", key, newVal),
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
					testAccEnvironmentAttributes("buddy_environment.env", &environment, name, "", "", true, "", buddy.EnvironmentPermissionAccessLevelDenied, buddy.EnvironmentPermissionAccessLevelUseOnly, "", "", ""),
				),
			},
			{
				Config: testAccEnvironmentPermissionsConfig(domain, projectName, identifier, name, email, groupName, buddy.EnvironmentPermissionAccessLevelDefault),
				Check: resource.ComposeTestCheckFunc(
					testAccEnvironmentGet("buddy_environment.env", &environment),
					testAccEnvironmentAttributes("buddy_environment.env", &environment, name, "", "", true, "", buddy.EnvironmentPermissionAccessLevelDefault, "", "", "", ""),
				),
			},
			{
				Config: testAccEnvironmentPermissionsUserGroupConfig(domain, projectName, identifier, name, email, groupName, buddy.EnvironmentPermissionAccessLevelManage, buddy.EnvironmentPermissionAccessLevelDefault, buddy.EnvironmentPermissionAccessLevelDenied),
				Check: resource.ComposeTestCheckFunc(
					testAccEnvironmentGet("buddy_environment.env", &environment),
					testAccEnvironmentAttributes("buddy_environment.env", &environment, name, "", "", true, "", buddy.EnvironmentPermissionAccessLevelManage, buddy.EnvironmentPermissionAccessLevelDefault, buddy.EnvironmentPermissionAccessLevelDenied, "", ""),
				),
			},
			{
				Config: testAccEnvironmentPermissionsGroupConfig(domain, projectName, identifier, name, email, groupName, buddy.EnvironmentPermissionAccessLevelUseOnly, buddy.EnvironmentPermissionAccessLevelManage),
				Check: resource.ComposeTestCheckFunc(
					testAccEnvironmentGet("buddy_environment.env", &environment),
					testAccEnvironmentAttributes("buddy_environment.env", &environment, name, "", "", true, "", buddy.EnvironmentPermissionAccessLevelUseOnly, "", buddy.EnvironmentPermissionAccessLevelManage, "", ""),
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
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccEnvironmentCheckDestroy,
		Steps: []resource.TestStep{
			// create env
			{
				Config: testAccEnvironmentConfig(domain, projectName, name, identifier, url, false, tag),
				Check: resource.ComposeTestCheckFunc(
					testAccEnvironmentGet("buddy_environment.env", &environment),
					testAccEnvironmentAttributes("buddy_environment.env", &environment, name, identifier, url, false, tag, "", "", "", "", ""),
				),
			},
			// update env
			{
				Config: testAccEnvironmentConfig(domain, projectName, newName, newIdentifier, newUrl, false, newTag),
				Check: resource.ComposeTestCheckFunc(
					testAccEnvironmentGet("buddy_environment.env", &environment),
					testAccEnvironmentAttributes("buddy_environment.env", &environment, newName, newIdentifier, newUrl, false, newTag, "", "", "", "", ""),
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

func testAccEnvironmentAttributes(n string, environment *buddy.Environment, name string, identifier string, url string, allPipelinesAllowed bool, tag string, othersLevel string, userLevel string, groupLevel string, varKey string, varVal string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsAllPipelinesAllowed, _ := strconv.ParseBool(attrs["all_pipelines_allowed"])
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
		if err := util.CheckFieldEqualAndSet("project.0.name", attrs["project.0.name"], environment.Project.Name); err != nil {
			return err
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
			if err := util.CheckFieldEqualAndSet("Permissions.Others", environment.Permissions.Others, buddy.EnvironmentPermissionAccessLevelDefault); err != nil {
				return err
			}
		}
		if userLevel != "" {
			if err := util.CheckFieldEqualAndSet("Permissions.Users[0].AccessLevel", environment.Permissions.Users[0].AccessLevel, userLevel); err != nil {
				return err
			}
		} else {
			if err := util.CheckIntFieldEqual("len(Permissions.Users)", len(environment.Permissions.Users), 0); err != nil {
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
		if varKey != "" {
			if err := util.CheckIntFieldEqual("Variables", len(environment.Variables), 1); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("Variables[0].Key", environment.Variables[0].Key, varKey); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("Variables[0].Value", environment.Variables[0].Value, varVal); err != nil {
				return err
			}
		} else {
			if err := util.CheckIntFieldEqual("Variables", len(environment.Variables), 0); err != nil {
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
		domain, projectName, environmentId, err := util.DecomposeTripleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		e, _, err := acc.ApiClient.EnvironmentService.Get(domain, projectName, environmentId)
		if err != nil {
			return err
		}
		*environment = *e
		return nil
	}
}

func testAccEnvironmentVariablesConfig(domain string, projectName string, identifier string, name string, key string, value string) string {
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
    identifier = "%s"
    name = "%s"
    var {
      key = "%s"
      value = "%s"
    }
}
`, domain, projectName, identifier, name, key, value)
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

func testAccEnvironmentConfig(domain string, projectName string, name string, identifier string, url string, allPipelinesAllowed bool, tag string) string {
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
    tags = ["%s"]
}
`, domain, projectName, name, identifier, url, allPipelinesAllowed, tag)
}

func testAccEnvironmentCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_environment" {
			continue
		}
		domain, projectName, environmentId, err := util.DecomposeTripleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		environment, resp, err := acc.ApiClient.EnvironmentService.Get(domain, projectName, environmentId)
		if err == nil && environment != nil {
			return util.ErrorResourceExists()
		}
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}
