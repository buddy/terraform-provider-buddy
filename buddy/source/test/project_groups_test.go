package test

import (
	"buddy-terraform/buddy/acc"
	"buddy-terraform/buddy/util"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
)

func TestAccSourceProjectGroups(t *testing.T) {
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name1 := "abc" + util.RandString(10)
	name2 := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceProjectGroupsConfig(domain, projectName, name1, name2),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceProjectGroupsAttributes("data.buddy_project_groups.pm", 2),
					testAccSourceProjectGroupsAttributes("data.buddy_project_groups.filter", 1),
				),
			},
		},
	})
}

func testAccSourceProjectGroupsAttributes(n string, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsGroupsCount, _ := strconv.Atoi(attrs["groups.#"])
		attrsGroupId, _ := strconv.Atoi(attrs["groups.0.group_id"])
		if err := util.CheckIntFieldEqual("groups.#", attrsGroupsCount, count); err != nil {
			return err
		}
		if err := util.CheckFieldSet("groups.0.html_url", attrs["groups.0.html_url"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("groups.0.name", attrs["groups.0.name"]); err != nil {
			return err
		}
		if err := util.CheckIntFieldSet("groups.0.group_id", attrsGroupId); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourceProjectGroupsConfig(domain string, projectName string, name1 string, name2 string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "w" {
   domain = "%s"
}

resource "buddy_project" "p" {
   domain = "${buddy_workspace.w.domain}"
   display_name = "%s"
}

resource "buddy_group" "g1" {
   domain = "${buddy_workspace.w.domain}"
   name ="%s"
}

resource "buddy_group" "g2" {
   domain = "${buddy_workspace.w.domain}"
   name ="%s"
}

resource "buddy_permission" "perm" {
   domain = "${buddy_workspace.w.domain}"
   name = "perm"
   pipeline_access_level = "READ_ONLY"
   repository_access_level = "READ_ONLY"
   sandbox_access_level = "READ_ONLY"
}

resource "buddy_project_group" "pg1" {
   domain = "${buddy_workspace.w.domain}"
   project_name = "${buddy_project.p.name}"
   permission_id = "${buddy_permission.perm.permission_id}"
   group_id = "${buddy_group.g1.group_id}"
}

resource "buddy_project_group" "pg2" {
   domain = "${buddy_workspace.w.domain}"
   project_name = "${buddy_project.p.name}"
   permission_id = "${buddy_permission.perm.permission_id}"
   group_id = "${buddy_group.g2.group_id}"
}

data "buddy_project_groups" "pm" {
   domain = "${buddy_workspace.w.domain}"
   project_name = "${buddy_project.p.name}"
   depends_on = [buddy_project_group.pg1, buddy_project_group.pg2]
}

data "buddy_project_groups" "filter" {
   domain = "${buddy_workspace.w.domain}"
   project_name = "${buddy_project.p.name}"
	name_regex = "^abc"
   depends_on = [buddy_project_group.pg1, buddy_project_group.pg2]
}

`, domain, projectName, name1, name2)
}
