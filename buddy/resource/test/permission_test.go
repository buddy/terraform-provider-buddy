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

func TestAccPermission(t *testing.T) {
	var permission buddy.Permission
	domain := util.UniqueString()
	name := util.RandString(5)
	pipelineAccessLevel := buddy.PermissionAccessLevelRunOnly
	repositoryAccessLevel := buddy.PermissionAccessLevelReadWrite
	sandboxAccessLevel := buddy.PermissionAccessLevelReadOnly
	targetAccessLevel := buddy.PermissionAccessLevelReadOnly
	environmentAccessLevel := buddy.PermissionAccessLevelUseOnly
	newName := util.RandString(5)
	newPipelineAccessLevel := buddy.PermissionAccessLevelReadWrite
	newRepositoryAccessLevel := buddy.PermissionAccessLevelManage
	newSandboxAccessLevel := buddy.PermissionAccessLevelReadWrite
	newDescription := util.RandString(5)
	newProjectTeamAccessLevel := buddy.PermissionAccessLevelManage
	newTargetAccessLevel := buddy.PermissionAccessLevelManage
	newEnvironmentAccessLevel := buddy.PermissionAccessLevelManage
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccPermissionCheckDestroy,
		Steps: []resource.TestStep{
			// create a permission
			{
				Config: testAccPermissionConfig(domain, name, pipelineAccessLevel, repositoryAccessLevel, sandboxAccessLevel, targetAccessLevel, environmentAccessLevel),
				Check: resource.ComposeTestCheckFunc(
					testAccPermissionGet("buddy_permission.bar", &permission),
					testAccPermissionAttributes("buddy_permission.bar", &permission, &testAccPermissionExpectedAttributes{
						Name:                   name,
						PipelineAccessLevel:    pipelineAccessLevel,
						RepositoryAccessLevel:  repositoryAccessLevel,
						SandboxAccessLevel:     sandboxAccessLevel,
						TargetAccessLevel:      targetAccessLevel,
						EnvironmentAccessLevel: environmentAccessLevel,
						Type:                   "CUSTOM",
						Description:            "",
					}),
				),
			},
			// update permission
			{
				Config: testAccPermissionUpdateConfig(domain, newName, newPipelineAccessLevel, newRepositoryAccessLevel, newProjectTeamAccessLevel, newSandboxAccessLevel, newTargetAccessLevel, newEnvironmentAccessLevel, newDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccPermissionGet("buddy_permission.bar", &permission),
					testAccPermissionAttributes("buddy_permission.bar", &permission, &testAccPermissionExpectedAttributes{
						Name:                   newName,
						PipelineAccessLevel:    newPipelineAccessLevel,
						RepositoryAccessLevel:  newRepositoryAccessLevel,
						SandboxAccessLevel:     newSandboxAccessLevel,
						ProjectTeamAccessLevel: newProjectTeamAccessLevel,
						TargetAccessLevel:      newTargetAccessLevel,
						EnvironmentAccessLevel: newEnvironmentAccessLevel,
						Type:                   "CUSTOM",
						Description:            newDescription,
					}),
				),
			},
			// import permission
			{
				ResourceName:      "buddy_permission.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

type testAccPermissionExpectedAttributes struct {
	Name                   string
	PipelineAccessLevel    string
	RepositoryAccessLevel  string
	SandboxAccessLevel     string
	ProjectTeamAccessLevel string
	TargetAccessLevel      string
	EnvironmentAccessLevel string
	Type                   string
	Description            string
}

func testAccPermissionAttributes(n string, permission *buddy.Permission, want *testAccPermissionExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsPermissionId, _ := strconv.Atoi(attrs["permission_id"])
		if err := util.CheckFieldEqualAndSet("Name", permission.Name, want.Name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("PipelineAccessLevel", permission.PipelineAccessLevel, want.PipelineAccessLevel); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("RepositoryAccessLevel", permission.RepositoryAccessLevel, want.RepositoryAccessLevel); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("SandboxAccessLevel", permission.SandboxAccessLevel, want.SandboxAccessLevel); err != nil {
			return err
		}
		if want.ProjectTeamAccessLevel != "" {
			if err := util.CheckFieldEqualAndSet("ProjectTeamAccessLevel", permission.ProjectTeamAccessLevel, want.ProjectTeamAccessLevel); err != nil {
				return err
			}
		}
		if err := util.CheckFieldEqualAndSet("Type", permission.Type, want.Type); err != nil {
			return err
		}
		if err := util.CheckFieldEqual("Description", permission.Description, want.Description); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], want.Name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("pipeline_access_level", attrs["pipeline_access_level"], want.PipelineAccessLevel); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("repository_access_level", attrs["repository_access_level"], want.RepositoryAccessLevel); err != nil {
			return err
		}
		if want.ProjectTeamAccessLevel != "" {
			if err := util.CheckFieldEqualAndSet("project_team_access_level", attrs["project_team_access_level"], want.ProjectTeamAccessLevel); err != nil {
				return err
			}
		}
		if err := util.CheckFieldEqualAndSet("sandbox_access_level", attrs["sandbox_access_level"], want.SandboxAccessLevel); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("target_access_level", attrs["target_access_level"], want.TargetAccessLevel); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("environment_access_level", attrs["environment_access_level"], want.EnvironmentAccessLevel); err != nil {
			return err
		}
		if err := util.CheckIntFieldEqualAndSet("permission_id", attrsPermissionId, permission.Id); err != nil {
			return err
		}
		if err := util.CheckFieldEqual("description", attrs["description"], want.Description); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("html_url", attrs["html_url"], permission.HtmlUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("type", attrs["type"], permission.Type); err != nil {
			return err
		}
		return nil
	}
}

func testAccPermissionGet(n string, permission *buddy.Permission) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		domain, pid, err := util.DecomposeDoubleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		permissionId, err := strconv.Atoi(pid)
		if err != nil {
			return err
		}
		p, _, err := acc.ApiClient.PermissionService.Get(domain, permissionId)
		if err != nil {
			return err
		}
		*permission = *p
		return nil
	}
}

func testAccPermissionConfig(domain string, name string, pipelineAccessLevel string, repositoryAccessLevel string, sandboxAccessLevel string, targetAccessLevel string, environmentAccessLevel string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_permission" "bar" {
   domain = "${buddy_workspace.foo.domain}"
   name = "%s"
   pipeline_access_level = "%s"
   repository_access_level = "%s"
   sandbox_access_level = "%s"
	 target_access_level = "%s"
	 environment_access_level = "%s"
}
`, domain, name, pipelineAccessLevel, repositoryAccessLevel, sandboxAccessLevel, targetAccessLevel, environmentAccessLevel)
}

func testAccPermissionUpdateConfig(domain string, name string, pipelineAccessLevel string, repositoryAccessLevel string, projectTeamAccessLevel string, sandboxAccessLevel string, targetAccessLevel string, environmentAccessLevel string, description string) string {
	return fmt.Sprintf(`

	resource "buddy_workspace" "foo" {
	   domain = "%s"
	}

	resource "buddy_permission" "bar" {
	   domain = "${buddy_workspace.foo.domain}"
	   name = "%s"
	   pipeline_access_level = "%s"
	   repository_access_level = "%s"
		 sandbox_access_level = "%s"
		 project_team_access_level = "%s"
	 	 target_access_level = "%s"
	 	 environment_access_level = "%s"
	   description = "%s"
	}

`, domain, name, pipelineAccessLevel, repositoryAccessLevel, sandboxAccessLevel, projectTeamAccessLevel, targetAccessLevel, environmentAccessLevel, description)
}

func testAccPermissionCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_permission" {
			continue
		}
		domain, pid, err := util.DecomposeDoubleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		permissionId, err := strconv.Atoi(pid)
		if err != nil {
			return err
		}
		permission, resp, err := acc.ApiClient.PermissionService.Get(domain, permissionId)
		if err == nil && permission != nil {
			return util.ErrorResourceExists()
		}
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}
