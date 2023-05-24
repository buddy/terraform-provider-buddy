package test

//
//import (
//	"buddy-terraform/buddy/acc"
//	"buddy-terraform/buddy/util"
//	"fmt"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
//	"testing"
//)
//
//func TestAccSourceProject(t *testing.T) {
//	domain := util.UniqueString()
//	name := util.UniqueString()
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			acc.PreCheck(t)
//		},
//		CheckDestroy:      acc.DummyCheckDestroy,
//		ProviderFactories: acc.ProviderFactories,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccSourceProjectConfig(domain, name),
//				Check: resource.ComposeTestCheckFunc(
//					testAccSourceProjectAttributes("data.buddy_project.name", name),
//					testAccSourceProjectAttributes("data.buddy_project.display_name", name),
//				),
//			},
//		},
//	})
//}
//
//func testAccSourceProjectAttributes(n string, name string) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		rs, ok := s.RootModule().Resources[n]
//		if !ok {
//			return fmt.Errorf("not found: %s", n)
//		}
//		attrs := rs.Primary.Attributes
//		if err := util.CheckFieldEqualAndSet("display_name", attrs["display_name"], name); err != nil {
//			return err
//		}
//		if err := util.CheckFieldEqualAndSet("name", attrs["name"], name); err != nil {
//			return err
//		}
//		if err := util.CheckFieldSet("html_url", attrs["html_url"]); err != nil {
//			return err
//		}
//		if err := util.CheckFieldEqualAndSet("status", attrs["status"], "ACTIVE"); err != nil {
//			return err
//		}
//		return nil
//	}
//}
//
//func testAccSourceProjectConfig(domain string, name string) string {
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
//data "buddy_project" "name" {
//    domain = "${buddy_workspace.foo.domain}"
//    name = "${buddy_project.proj.name}"
//}
//
//data "buddy_project" "display_name" {
//    domain = "${buddy_workspace.foo.domain}"
//    display_name = "${buddy_project.proj.display_name}"
//}
//`, domain, name)
//}
