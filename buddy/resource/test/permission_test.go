package test

import (
	"buddy-terraform/buddy/acc"
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strconv"
	"testing"
)

func TestAccPermission(t *testing.T) {
	var permission api.Permission
	domain := util.UniqueString()
	name := util.RandString(5)
	pipelineAccessLevel := api.PermissionAccessLevelRunOnly
	repositoryAccessLevel := api.PermissionAccessLevelReadWrite
	sandboxAccessLevel := api.PermissionAccessLevelReadOnly
	newName := util.RandString(5)
	newPipelineAccessLevel := api.PermissionAccessLevelDenied
	newRepositoryAccessLevel := api.PermissionAccessLevelManage
	newSandboxAccessLevel := api.PermissionAccessLevelReadWrite
	newDescription := util.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acc.PreCheck(t) },
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccPermissionCheckDestroy,
		Steps: []resource.TestStep{
			// create a permission
			{
				Config: testAccPermissionConfig(domain, name, pipelineAccessLevel, repositoryAccessLevel, sandboxAccessLevel),
				Check: resource.ComposeTestCheckFunc(
					testAccPermissionGet("buddy_permission.bar", &permission),
					testAccPermissionAttributes("buddy_permission.bar", &permission, &testAccPermissionExpectedAttributes{
						Name:                  name,
						PipelineAccessLevel:   pipelineAccessLevel,
						RepositoryAccessLevel: repositoryAccessLevel,
						SandboxAccessLevel:    sandboxAccessLevel,
						Type:                  "CUSTOM",
						Description:           "",
					}),
				),
			},
			// update permission
			{
				Config: testAccPermissionUpdateConfig(domain, newName, newPipelineAccessLevel, newRepositoryAccessLevel, newSandboxAccessLevel, newDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccPermissionGet("buddy_permission.bar", &permission),
					testAccPermissionAttributes("buddy_permission.bar", &permission, &testAccPermissionExpectedAttributes{
						Name:                  newName,
						PipelineAccessLevel:   newPipelineAccessLevel,
						RepositoryAccessLevel: newRepositoryAccessLevel,
						SandboxAccessLevel:    newSandboxAccessLevel,
						Type:                  "CUSTOM",
						Description:           newDescription,
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
	Name                  string
	PipelineAccessLevel   string
	RepositoryAccessLevel string
	SandboxAccessLevel    string
	Type                  string
	Description           string
}

func testAccPermissionAttributes(n string, permission *api.Permission, want *testAccPermissionExpectedAttributes) resource.TestCheckFunc {
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
		if err := util.CheckFieldEqualAndSet("sandbox_access_level", attrs["sandbox_access_level"], want.SandboxAccessLevel); err != nil {
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

func testAccPermissionGet(n string, permission *api.Permission) resource.TestCheckFunc {
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

func testAccPermissionConfig(domain string, name string, pipelineAccessLevel string, repositoryAccessLevel string, sandboxAccessLevel string) string {
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
}
`, domain, name, pipelineAccessLevel, repositoryAccessLevel, sandboxAccessLevel)
}

func testAccPermissionUpdateConfig(domain string, name string, pipelineAccessLevel string, repositoryAccessLevel string, sandboxAccessLevel string, description string) string {
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
    description = "%s"
}
`, domain, name, pipelineAccessLevel, repositoryAccessLevel, sandboxAccessLevel, description)
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
		if resp.StatusCode != 404 {
			return err
		}
	}
	return nil
}
