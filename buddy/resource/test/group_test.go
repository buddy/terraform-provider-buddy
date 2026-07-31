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

func TestAccGroup(t *testing.T) {
	var group buddy.Group
	var permission buddy.Permission
	domain := util.UniqueString()
	name := util.RandString(5)
	newName := util.RandString(5)
	newNote := util.RandString(5)
	newAgentNote := util.RandString(5)
	legacyDescription := util.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccGroupCheckDestroy,
		Steps: []resource.TestStep{
			// create group
			{
				Config: testAccGroupConfig(domain, name),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupGet("buddy_group.bar", &group),
					testAccGroupAttributes("buddy_group.bar", &group, name, "", "", false, nil),
				),
			},
			// update group
			{
				Config: testAccGroupUpdateConfig(domain, newName, newNote, newAgentNote),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupGet("buddy_group.bar", &group),
					testAccGroupAttributes("buddy_group.bar", &group, newName, newNote, newAgentNote, false, nil),
				),
			},
			// update group assign
			{
				Config: testAccGroupUpdateProjectAssignConfig(domain, newName, newNote, newAgentNote, false),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupGet("buddy_group.bar", &group),
					testAccPermissionGet("buddy_permission.perm", &permission),
					testAccGroupAttributes("buddy_group.bar", &group, newName, newNote, newAgentNote, false, &permission),
				),
			},
			// update group assign
			{
				Config: testAccGroupUpdateProjectAssignConfig(domain, newName, newNote, newAgentNote, true),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupGet("buddy_group.bar", &group),
					testAccPermissionGet("buddy_permission.perm", &permission),
					testAccGroupAttributes("buddy_group.bar", &group, newName, newNote, newAgentNote, true, &permission),
				),
			},
			// null group assign
			{
				Config: testAccGroupUpdateProjectAssignConfig(domain, newName, newNote, newAgentNote, false),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupGet("buddy_group.bar", &group),
					testAccPermissionGet("buddy_permission.perm", &permission),
					testAccGroupAttributes("buddy_group.bar", &group, newName, newNote, newAgentNote, false, &permission),
				),
			},
			// deprecated description feeds note, agent note is kept, it is optional & computed
			{
				Config: testAccGroupUpdateDescriptionConfig(domain, newName, legacyDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupGet("buddy_group.bar", &group),
					testAccGroupAttributes("buddy_group.bar", &group, newName, legacyDescription, newAgentNote, false, nil),
				),
			},
			// migrate from deprecated description to note
			{
				Config: testAccGroupUpdateConfig(domain, newName, newNote, newAgentNote),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupGet("buddy_group.bar", &group),
					testAccGroupAttributes("buddy_group.bar", &group, newName, newNote, newAgentNote, false, nil),
				),
			},
			// note & agent note dropped from the config keep their values, they are optional & computed
			{
				Config: testAccGroupConfig(domain, newName),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupGet("buddy_group.bar", &group),
					testAccGroupAttributes("buddy_group.bar", &group, newName, newNote, newAgentNote, false, nil),
				),
			},
			// note & agent note explicitly emptied are cleared
			{
				Config: testAccGroupUpdateConfig(domain, newName, "", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupGet("buddy_group.bar", &group),
					testAccGroupAttributes("buddy_group.bar", &group, newName, "", "", false, nil),
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

func testAccGroupAttributes(n string, group *buddy.Group, name string, note string, agentNote string, autoAssign bool, defPerm *buddy.Permission) resource.TestCheckFunc {
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
		if err := util.CheckFieldEqual("Note", group.Note, note); err != nil {
			return err
		}
		if err := util.CheckFieldEqual("AgentNote", group.AgentNote, agentNote); err != nil {
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
		if err := util.CheckFieldEqual("note", attrs["note"], note); err != nil {
			return err
		}
		if err := util.CheckFieldEqual("agent_note", attrs["agent_note"], agentNote); err != nil {
			return err
		}
		// deprecated description mirrors note
		if err := util.CheckFieldEqual("description", attrs["description"], note); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("auto_assign_to_new_projects", attrsAutoAssignToProjects, autoAssign); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("AutoAssignToNewProjects", group.AutoAssignToNewProjects, autoAssign); err != nil {
			return err
		}
		if defPerm != nil && autoAssign {
			if err := util.CheckIntFieldEqual("AutoAssignPermissionSetId", group.AutoAssignPermissionSetId, defPerm.Id); err != nil {
				return err
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

func testAccGroupUpdateConfig(domain string, name string, note string, agentNote string) string {
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
	   note = "%s"
	   agent_note = "%s"
	}

`, domain, name, note, agentNote)
}

func testAccGroupUpdateDescriptionConfig(domain string, name string, description string) string {
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

func testAccGroupUpdateProjectAssignConfig(domain string, name string, note string, agentNote string, autoAssign bool) string {
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
	   note = "%s"
	   agent_note = "%s"
		auto_assign_to_new_projects = %t
		auto_assign_permission_set_id = "${buddy_permission.perm.permission_id}"
	}

`, domain, name, note, agentNote, autoAssign)
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
