package test

import (
	"buddy-terraform/buddy/acc"
	"buddy-terraform/buddy/util"
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strconv"
	"testing"
)

func TestAccProjectGroup(t *testing.T) {
	var group buddy.ProjectGroup
	domain := util.UniqueString()
	nameA := util.RandString(10)
	nameB := util.RandString(10)
	projectDisplayNameA := util.RandString(10)
	projectDisplayNameB := util.RandString(10)
	permissionNameA := util.RandString(10)
	permissionNameB := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccProjectGroupCheckDestroy,
		Steps: []resource.TestStep{
			// create
			{
				Config: testAccProjectGroupConfig(domain, nameA, nameB, projectDisplayNameA, projectDisplayNameB, permissionNameA, permissionNameB, "a", "a", "a"),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectGroupGet("buddy_project_group.bar", &group),
					testAccProjectGroupAttributes("buddy_project_group.bar", &group, permissionNameA),
				),
			},
			// update group
			{
				Config: testAccProjectGroupConfig(domain, nameA, nameB, projectDisplayNameA, projectDisplayNameB, permissionNameA, permissionNameB, "a", "b", "a"),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectGroupGet("buddy_project_group.bar", &group),
					testAccProjectGroupAttributes("buddy_project_group.bar", &group, permissionNameA),
				),
			},
			// update project
			{
				Config: testAccProjectGroupConfig(domain, nameA, nameB, projectDisplayNameA, projectDisplayNameB, permissionNameA, permissionNameB, "b", "b", "a"),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectGroupGet("buddy_project_group.bar", &group),
					testAccProjectGroupAttributes("buddy_project_group.bar", &group, permissionNameA),
				),
			},
			// update permission
			{
				Config: testAccProjectGroupConfig(domain, nameA, nameB, projectDisplayNameA, projectDisplayNameB, permissionNameA, permissionNameB, "b", "b", "b"),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectGroupGet("buddy_project_group.bar", &group),
					testAccProjectGroupAttributes("buddy_project_group.bar", &group, permissionNameB),
				),
			},
			// import
			{
				ResourceName:      "buddy_project_group.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProjectGroupAttributes(n string, group *buddy.ProjectGroup, permissionName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsGroupId, _ := strconv.Atoi(attrs["group_id"])
		attrsPermissionId, _ := strconv.Atoi(attrs["permission_id"])
		attrsPermissionPermissionId, _ := strconv.Atoi(attrs["permission.0.permission_id"])
		if err := util.CheckIntFieldEqualAndSet("group_id", attrsGroupId, group.Id); err != nil {
			return err
		}
		if err := util.CheckIntFieldEqualAndSet("permission_id", attrsPermissionId, group.PermissionSet.Id); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("html_url", attrs["html_url"], group.HtmlUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], group.Name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("permission.0.html_url", attrs["permission.0.html_url"], group.PermissionSet.HtmlUrl); err != nil {
			return err
		}
		if err := util.CheckIntFieldEqualAndSet("permission.0.permission_id", attrsPermissionPermissionId, group.PermissionSet.Id); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("permission.0.name", attrs["permission.0.name"], group.PermissionSet.Name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("permission.0.name", attrs["permission.0.name"], permissionName); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("permission.0.type", attrs["permission.0.type"], group.PermissionSet.Type); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("permission.0.pipeline_access_level", attrs["permission.0.pipeline_access_level"], group.PermissionSet.PipelineAccessLevel); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("permission.0.repository_access_level", attrs["permission.0.repository_access_level"], group.PermissionSet.RepositoryAccessLevel); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("permission.0.sandbox_access_level", attrs["permission.0.sandbox_access_level"], group.PermissionSet.SandboxAccessLevel); err != nil {
			return err
		}
		return nil
	}
}

func testAccProjectGroupGet(n string, group *buddy.ProjectGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		domain, projectName, gid, err := util.DecomposeTripleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		groupId, err := strconv.Atoi(gid)
		if err != nil {
			return err
		}
		g, _, err := acc.ApiClient.ProjectGroupService.GetProjectGroup(domain, projectName, groupId)
		if err != nil {
			return err
		}
		*group = *g
		return nil
	}
}

func testAccProjectGroupConfig(domain string, nameA string, nameB string, projectDisplayNameA string, projectDisplayNameB string, permissionNameA string, permissionNameB string, whichProject string, whichGroup string, whichPermission string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_group" "a" {
    domain = "${buddy_workspace.foo.domain}"
    name = "%s"
}

resource "buddy_group" "b" {
    domain = "${buddy_workspace.foo.domain}"
    name = "%s"
}

resource "buddy_project" "a" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_project" "b" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_permission" "a" {
    domain = "${buddy_workspace.foo.domain}"
    name = "%s"
    pipeline_access_level = "%s"
    repository_access_level = "%s"
	sandbox_access_level = "%s"
}

resource "buddy_permission" "b" {
    domain = "${buddy_workspace.foo.domain}"
    name = "%s"
    pipeline_access_level = "%s"
    repository_access_level = "%s"
	sandbox_access_level = "%s"
}

resource "buddy_project_group" "bar" {
	domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.%s.name}"
	group_id = "${buddy_group.%s.group_id}"
	permission_id = "${buddy_permission.%s.permission_id}"
}

`,
		domain,
		nameA,
		nameB,
		projectDisplayNameA,
		projectDisplayNameB,
		permissionNameA,
		buddy.PermissionAccessLevelReadOnly,
		buddy.PermissionAccessLevelReadOnly,
		buddy.PermissionAccessLevelReadOnly,
		permissionNameB,
		buddy.PermissionAccessLevelReadWrite,
		buddy.PermissionAccessLevelReadWrite,
		buddy.PermissionAccessLevelReadWrite,
		whichProject,
		whichGroup,
		whichPermission,
	)
}

func testAccProjectGroupCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_project_group" {
			continue
		}
		domain, projectName, gid, err := util.DecomposeTripleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		groupId, err := strconv.Atoi(gid)
		if err != nil {
			return err
		}
		group, resp, err := acc.ApiClient.ProjectGroupService.GetProjectGroup(domain, projectName, groupId)
		if err == nil && group != nil {
			return util.ErrorResourceExists()
		}
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}
