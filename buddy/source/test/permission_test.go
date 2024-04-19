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

func TestAccSourcePermission(t *testing.T) {
	domain := util.UniqueString()
	name := util.RandString(10)
	pipelineAccessLevel := buddy.PermissionAccessLevelReadWrite
	repositoryAccessLevel := buddy.PermissionAccessLevelReadOnly
	sandboxAccessLevel := buddy.PermissionAccessLevelDenied
	projectTeamAccessLevel := buddy.PermissionAccessLevelReadOnly
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourcePermissionConfig(domain, name, pipelineAccessLevel, repositoryAccessLevel, sandboxAccessLevel, projectTeamAccessLevel),
				Check: resource.ComposeTestCheckFunc(
					testAccSourcePermissionAttributes("data.buddy_permission.id", name, pipelineAccessLevel, repositoryAccessLevel, sandboxAccessLevel, projectTeamAccessLevel),
					testAccSourcePermissionAttributes("data.buddy_permission.name", name, pipelineAccessLevel, repositoryAccessLevel, sandboxAccessLevel, projectTeamAccessLevel),
				),
			},
		},
	})
}

func testAccSourcePermissionAttributes(n string, name string, pipelineAccessLevel string, repositoryAccessLevel string, sandboxAccessLevel string, projectTeamAccessLevel string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsPermissionId, _ := strconv.Atoi(attrs["permission_id"])
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("pipeline_access_level", attrs["pipeline_access_level"], pipelineAccessLevel); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("repository_access_level", attrs["repository_access_level"], repositoryAccessLevel); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("sandbox_access_level", attrs["sandbox_access_level"], sandboxAccessLevel); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("project_team_access_level", attrs["project_team_access_level"], projectTeamAccessLevel); err != nil {
			return err
		}
		if err := util.CheckIntFieldSet("permission_id", attrsPermissionId); err != nil {
			return err
		}
		if err := util.CheckFieldSet("html_url", attrs["html_url"]); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("type", attrs["type"], "CUSTOM"); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourcePermissionConfig(domain string, name string, pipelineAccessLevel string, repositoryAccessLevel string, sandboxAccessLevel string, projectTeamAccessLevel string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_permission" "perm" {
   domain = "${buddy_workspace.foo.domain}"
   name = "%s"
   pipeline_access_level = "%s"
   repository_access_level = "%s"
	sandbox_access_level = "%s"
	project_team_access_level = "%s"
}

data "buddy_permission" "id" {
   domain = "${buddy_workspace.foo.domain}"
   permission_id = "${buddy_permission.perm.permission_id}"
}

data "buddy_permission" "name" {
   domain = "${buddy_workspace.foo.domain}"
   name = "${buddy_permission.perm.name}"
}
`, domain, name, pipelineAccessLevel, repositoryAccessLevel, sandboxAccessLevel, projectTeamAccessLevel)
}
