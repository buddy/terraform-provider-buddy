package test

//
//import (
//	"buddy-terraform/buddy/acc"
//	"buddy-terraform/buddy/util"
//	"fmt"
//	"github.com/buddy/api-go-sdk/buddy"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
//	"strconv"
//	"testing"
//)
//
//func TestAccSourceProjectMember(t *testing.T) {
//	domain := util.UniqueString()
//	memberEmail := util.RandEmail()
//	projectName := util.UniqueString()
//	permissionName := util.RandString(10)
//	pipelineAccessLevel := buddy.PermissionAccessLevelRunOnly
//	repoAccessLevel := buddy.PermissionAccessLevelReadWrite
//	sandboxAccessLevel := buddy.PermissionAccessLevelReadWrite
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			acc.PreCheck(t)
//		},
//		CheckDestroy:      acc.DummyCheckDestroy,
//		ProviderFactories: acc.ProviderFactories,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccSourceProjectMemberConfig(domain, memberEmail, projectName, permissionName, pipelineAccessLevel, repoAccessLevel, sandboxAccessLevel),
//				Check: resource.ComposeTestCheckFunc(
//					testAccSourceProjectMemberAttributes("data.buddy_project_member.bar", memberEmail, permissionName, pipelineAccessLevel, repoAccessLevel, sandboxAccessLevel),
//				),
//			},
//		},
//	})
//}
//
//func testAccSourceProjectMemberAttributes(n string, email string, permName string, pipelineAccessLevel string, repoAccessLevel string, sandboxAccessLevel string) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		rs, ok := s.RootModule().Resources[n]
//		if !ok {
//			return fmt.Errorf("not found: %s", n)
//		}
//		attrs := rs.Primary.Attributes
//		attrsPermissionPermissionId, _ := strconv.Atoi(attrs["permission.0.permission_id"])
//		if err := util.CheckFieldSet("html_url", attrs["html_url"]); err != nil {
//			return err
//		}
//		if err := util.CheckFieldSet("avatar_url", attrs["avatar_url"]); err != nil {
//			return err
//		}
//		if err := util.CheckFieldEqualAndSet("email", attrs["email"], email); err != nil {
//			return err
//		}
//		if err := util.CheckFieldSet("permission.0.html_url", attrs["permission.0.html_url"]); err != nil {
//			return err
//		}
//		if err := util.CheckIntFieldSet("permission.0.permission_id", attrsPermissionPermissionId); err != nil {
//			return err
//		}
//		if err := util.CheckFieldEqualAndSet("permission.0.name", attrs["permission.0.name"], permName); err != nil {
//			return err
//		}
//		if err := util.CheckFieldEqualAndSet("permission.0.type", attrs["permission.0.type"], "CUSTOM"); err != nil {
//			return err
//		}
//		if err := util.CheckFieldEqualAndSet("permission.0.pipeline_access_level", attrs["permission.0.pipeline_access_level"], pipelineAccessLevel); err != nil {
//			return err
//		}
//		if err := util.CheckFieldEqualAndSet("permission.0.repository_access_level", attrs["permission.0.repository_access_level"], repoAccessLevel); err != nil {
//			return err
//		}
//		if err := util.CheckFieldEqualAndSet("permission.0.sandbox_access_level", attrs["permission.0.sandbox_access_level"], sandboxAccessLevel); err != nil {
//			return err
//		}
//		return nil
//	}
//}
//
//func testAccSourceProjectMemberConfig(domain string, email string, projectName string, permissionName string, pipelineAccessLevel string, repoAccessLevel string, sandboxAccessLevel string) string {
//	return fmt.Sprintf(`
//resource "buddy_workspace" "foo" {
//    domain = "%s"
//}
//
//resource "buddy_project" "proj" {
//    domain = "${buddy_workspace.foo.domain}"
//    display_name = "%s"
//}
//
//resource "buddy_member" "mem" {
//    domain = "${buddy_workspace.foo.domain}"
//    email = "%s"
//}
//
//resource "buddy_permission" "perm" {
//    domain = "${buddy_workspace.foo.domain}"
//    name = "%s"
//    pipeline_access_level = "%s"
//    repository_access_level = "%s"
//	sandbox_access_level = "%s"
//}
//
//resource "buddy_project_member" "bpm" {
//    domain = "${buddy_workspace.foo.domain}"
//	project_name = "${buddy_project.proj.name}"
//	member_id = "${buddy_member.mem.member_id}"
//	permission_id = "${buddy_permission.perm.permission_id}"
//}
//
//data "buddy_project_member" "bar" {
//    domain = "${buddy_workspace.foo.domain}"
//	project_name = "${buddy_project.proj.name}"
//	member_id = "${buddy_member.mem.member_id}"
//    depends_on = [buddy_project_member.bpm]
//}
//`, domain, projectName, email, permissionName, pipelineAccessLevel, repoAccessLevel, sandboxAccessLevel)
//}
