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

func TestAccGroup(t *testing.T) {
	var group buddy.Group
	var permission buddy.Permission
	domain := util.UniqueString()
	name := util.RandString(5)
	newName := util.RandString(5)
	newDescription := util.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acc.PreCheck(t) },
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccGroupCheckDestroy,
		Steps: []resource.TestStep{
			// create group
			{
				Config: testAccGroupConfig(domain, name),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupGet("buddy_group.bar", &group),
					testAccGroupAttributes("buddy_group.bar", &group, name, "", false, nil),
				),
			},
			// update group
			{
				Config: testAccGroupUpdateConfig(domain, newName, newDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupGet("buddy_group.bar", &group),
					testAccGroupAttributes("buddy_group.bar", &group, newName, newDescription, false, nil),
				),
			},
			// update group assign
			{
				Config: testAccGroupUpdateProjectAssignConfig(domain, newName, newDescription, false),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupGet("buddy_group.bar", &group),
					testAccPermissionGet("buddy_permission.perm", &permission),
					testAccGroupAttributes("buddy_group.bar", &group, newName, newDescription, false, &permission),
				),
			},
			// update group assign
			{
				Config: testAccGroupUpdateProjectAssignConfig(domain, newName, newDescription, true),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupGet("buddy_group.bar", &group),
					testAccPermissionGet("buddy_permission.perm", &permission),
					testAccGroupAttributes("buddy_group.bar", &group, newName, newDescription, true, &permission),
				),
			},
			// null desc
			{
				Config: testAccGroupConfig(domain, newName),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupGet("buddy_group.bar", &group),
					testAccGroupAttributes("buddy_group.bar", &group, newName, "", false, nil),
				),
			},
			// import group
			{
				ResourceName:            "buddy_group.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"auto_assign_permission_set_id"},
			},
		},
	})
}

func testAccGroupAttributes(n string, group *buddy.Group, name string, description string, autoAssign bool, defPerm *buddy.Permission) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsAutoAssignToProjects, _ := strconv.ParseBool(attrs["auto_assign_to_new_projects"])
		attrsAutoAssignToProjectsPermissionId, _ := strconv.Atoi(attrs["auto_assign_permission_set_id"])
		if err := util.CheckFieldEqualAndSet("Name", group.Name, name); err != nil {
			return err
		}
		if err := util.CheckFieldEqual("Description", group.Description, description); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("group_id", attrs["group_id"], strconv.Itoa(group.Id)); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("html_url", attrs["html_url"], group.HtmlUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqual("description", attrs["description"], group.Description); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("auto_assign_to_new_projects", attrsAutoAssignToProjects, autoAssign); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("AutoAssignToNewProjects", group.AutoAssignToNewProjects, autoAssign); err != nil {
			return err
		}
		if defPerm != nil {
			if autoAssign {
				if err := util.CheckIntFieldEqual("AutoAssignPermissionSetId", group.AutoAssignPermissionSetId, defPerm.Id); err != nil {
					return err
				}
			}
			if err := util.CheckIntFieldEqual("auto_assign_permission_set_id", attrsAutoAssignToProjectsPermissionId, defPerm.Id); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccGroupGet(n string, group *buddy.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		domain, gid, err := util.DecomposeDoubleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		groupId, err := strconv.Atoi(gid)
		if err != nil {
			return err
		}
		g, _, err := acc.ApiClient.GroupService.Get(domain, groupId)
		if err != nil {
			return err
		}
		*group = *g
		return nil
	}
}

func testAccGroupUpdateConfig(domain string, name string, description string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_permission" "perm" {
    domain = "${buddy_workspace.foo.domain}"
    name = "test"
    pipeline_access_level = "READ_ONLY"
    repository_access_level = "READ_ONLY"
	sandbox_access_level = "READ_ONLY"
}

resource "buddy_group" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    name = "%s"
    description = "%s"
}
`, domain, name, description)
}

func testAccGroupUpdateProjectAssignConfig(domain string, name string, description string, autoAssign bool) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_permission" "perm" {
    domain = "${buddy_workspace.foo.domain}"
    name = "test"
    pipeline_access_level = "READ_ONLY"
    repository_access_level = "READ_ONLY"
	sandbox_access_level = "READ_ONLY"
}

resource "buddy_group" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    name = "%s"
    description = "%s"
	auto_assign_to_new_projects = %t
	auto_assign_permission_set_id = "${buddy_permission.perm.permission_id}"
}
`, domain, name, description, autoAssign)
}

func testAccGroupConfig(domain string, name string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_permission" "perm" {
    domain = "${buddy_workspace.foo.domain}"
    name = "test"
    pipeline_access_level = "READ_ONLY"
    repository_access_level = "READ_ONLY"
	sandbox_access_level = "READ_ONLY"
}

resource "buddy_group" "bar" {
   domain = "${buddy_workspace.foo.domain}"
   name = "%s"
}
`, domain, name)
}

func testAccGroupCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_group" {
			continue
		}
		domain, gid, err := util.DecomposeDoubleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		groupId, err := strconv.Atoi(gid)
		if err != nil {
			return err
		}
		group, resp, err := acc.ApiClient.GroupService.Get(domain, groupId)
		if err == nil && group != nil {
			return util.ErrorResourceExists()
		}
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}
