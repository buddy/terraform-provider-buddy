package test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"
)

func TestAccSourceProjectMembers(t *testing.T) {
	domain := util.UniqueString()
	projectName := util.UniqueString()
	email1 := util.RandEmail()
	email2 := util.RandEmail()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceProjectMembersConfig(domain, projectName, email1, email2),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceProjectMembersAttributes("data.buddy_project_members.pm", 3),
					testAccSourceProjectMembersAttributes("data.buddy_project_members.filter", 1),
				),
			},
		},
	})
}

func testAccSourceProjectMembersAttributes(n string, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsMembersCount, _ := strconv.Atoi(attrs["members.#"])
		attrsMemberId, _ := strconv.Atoi(attrs["members.0.member_id"])
		if err := util.CheckIntFieldEqual("members.#", attrsMembersCount, count); err != nil {
			return err
		}
		if err := util.CheckFieldSet("members.0.html_url", attrs["members.0.html_url"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("members.0.avatar_url", attrs["members.0.avatar_url"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("members.0.email", attrs["members.0.email"]); err != nil {
			return err
		}
		if err := util.CheckIntFieldSet("members.0.member_id", attrsMemberId); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourceProjectMembersConfig(domain string, projectName string, email1 string, email2 string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "w" {
   domain = "%s"
}

resource "buddy_project" "p" {
   domain = "${buddy_workspace.w.domain}"
   display_name = "%s"
}

resource "buddy_profile" "me" {
	name = "abcdef"
}

resource "buddy_member" "m1" {
   domain = "${buddy_workspace.w.domain}"
   email ="%s"
}

resource "buddy_member" "m2" {
   domain = "${buddy_workspace.w.domain}"
   email ="%s"
}

resource "buddy_permission" "perm" {
   domain = "${buddy_workspace.w.domain}"
   name = "perm"
   pipeline_access_level = "READ_ONLY"
   repository_access_level = "READ_ONLY"
   sandbox_access_level = "READ_ONLY"
}

resource "buddy_project_member" "pm1" {
   domain = "${buddy_workspace.w.domain}"
   project_name = "${buddy_project.p.name}"
   permission_id = "${buddy_permission.perm.permission_id}"
   member_id = "${buddy_member.m1.member_id}"
}

resource "buddy_project_member" "pm2" {
   domain = "${buddy_workspace.w.domain}"
   project_name = "${buddy_project.p.name}"
   permission_id = "${buddy_permission.perm.permission_id}"
   member_id = "${buddy_member.m2.member_id}"
}

data "buddy_project_members" "pm" {
   domain = "${buddy_workspace.w.domain}"
   project_name = "${buddy_project.p.name}"
   depends_on = [buddy_project_member.pm1, buddy_project_member.pm2, buddy_profile.me]
}

data "buddy_project_members" "filter" {
   domain = "${buddy_workspace.w.domain}"
   project_name = "${buddy_project.p.name}"
	name_regex = "^abc"
   depends_on = [buddy_project_member.pm1, buddy_project_member.pm2, buddy_profile.me]
}

`, domain, projectName, email1, email2)
}
