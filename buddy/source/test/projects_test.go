package test

//
//import (
//	"buddy-terraform/buddy/acc"
//	"buddy-terraform/buddy/util"
//	"fmt"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
//	"strconv"
//	"testing"
//)
//
//func TestAccSourceProjects(t *testing.T) {
//	domain := util.UniqueString()
//	name1 := "aaa" + util.UniqueString()
//	name2 := util.UniqueString() + "bbb"
//	name3 := util.UniqueString()
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			acc.PreCheck(t)
//		},
//		CheckDestroy:      acc.DummyCheckDestroy,
//		ProviderFactories: acc.ProviderFactories,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccSourceProjectsConfig(domain, name1, name2, name3),
//				Check: resource.ComposeTestCheckFunc(
//					testAccSourceProjectsAttributes("data.buddy_projects.a", 1, name1),
//					testAccSourceProjectsAttributes("data.buddy_projects.b", 1, name2),
//					testAccSourceProjectsAttributes("data.buddy_projects.c", 0, ""),
//					testAccSourceProjectsAttributes("data.buddy_projects.d", 3, ""),
//				),
//			},
//		},
//	})
//}
//
//func testAccSourceProjectsAttributes(n string, count int, name string) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		rs, ok := s.RootModule().Resources[n]
//		if !ok {
//			return fmt.Errorf("not found: %s", n)
//		}
//		attrs := rs.Primary.Attributes
//		attrsProjectsCount, _ := strconv.Atoi(attrs["projects.#"])
//		if err := util.CheckIntFieldEqual("projects.#", attrsProjectsCount, count); err != nil {
//			return err
//		}
//		if count > 0 {
//			if name != "" {
//				if err := util.CheckFieldEqualAndSet("projects.0.name", attrs["projects.0.name"], name); err != nil {
//					return err
//				}
//			} else {
//				if err := util.CheckFieldSet("projects.0.name", attrs["projects.0.name"]); err != nil {
//					return err
//				}
//			}
//			if err := util.CheckFieldSet("projects.0.display_name", attrs["projects.0.display_name"]); err != nil {
//				return err
//			}
//			if err := util.CheckFieldSet("projects.0.status", attrs["projects.0.status"]); err != nil {
//				return err
//			}
//			if err := util.CheckFieldSet("projects.0.html_url", attrs["projects.0.html_url"]); err != nil {
//				return err
//			}
//		}
//		return nil
//	}
//}
//
//func testAccSourceProjectsConfig(domain string, name1 string, name2 string, name3 string) string {
//	return fmt.Sprintf(`
//resource "buddy_workspace" "w" {
//    domain = "%s"
//}
//
//resource "buddy_project" "a" {
//    domain = "${buddy_workspace.w.domain}"
//    display_name = "%s"
//}
//
//resource "buddy_project" "b" {
//    domain = "${buddy_workspace.w.domain}"
//    display_name = "%s"
//}
//
//resource "buddy_project" "c" {
//    domain = "${buddy_workspace.w.domain}"
//    display_name = "%s"
//}
//
//data "buddy_projects" "a" {
//	depends_on = [buddy_project.a, buddy_project.b, buddy_project.c]
//    domain = "${buddy_workspace.w.domain}"
//    name_regex = "^aaa"
//}
//
//data "buddy_projects" "b" {
//	depends_on = [buddy_project.a, buddy_project.b, buddy_project.c]
//    domain = "${buddy_workspace.w.domain}"
//    display_name_regex = "bbb$"
//}
//
//data "buddy_projects" "c" {
//	depends_on = [buddy_project.a, buddy_project.b, buddy_project.c]
//    domain = "${buddy_workspace.w.domain}"
//    status = "CLOSED"
//}
//
//data "buddy_projects" "d" {
//	depends_on = [buddy_project.a, buddy_project.b, buddy_project.c]
//    domain = "${buddy_workspace.w.domain}"
//}
//`, domain, name1, name2, name3)
//}
