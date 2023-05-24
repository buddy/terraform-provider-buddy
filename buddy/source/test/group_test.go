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
//func TestAccSourceGroup(t *testing.T) {
//	domain := util.UniqueString()
//	name := util.RandString(5)
//	desc := util.RandString(5)
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			acc.PreCheck(t)
//		},
//		CheckDestroy:      acc.DummyCheckDestroy,
//		ProviderFactories: acc.ProviderFactories,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccSourceGroupConfig(domain, name, desc),
//				Check: resource.ComposeTestCheckFunc(
//					testAccSourceGroupAttributes("data.buddy_group.id", name, desc),
//					testAccSourceGroupAttributes("data.buddy_group.name", name, desc),
//				),
//			},
//		},
//	})
//}
//
//func testAccSourceGroupAttributes(n string, name string, desc string) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		rs, ok := s.RootModule().Resources[n]
//		if !ok {
//			return fmt.Errorf("not found: %s", n)
//		}
//		attrs := rs.Primary.Attributes
//		attrsGroupId, _ := strconv.Atoi(attrs["group_id"])
//		if err := util.CheckFieldEqualAndSet("Name", attrs["name"], name); err != nil {
//			return err
//		}
//		if err := util.CheckIntFieldSet("group_id", attrsGroupId); err != nil {
//			return err
//		}
//		if err := util.CheckFieldSet("html_url", attrs["html_url"]); err != nil {
//			return err
//		}
//		if err := util.CheckFieldEqualAndSet("description", attrs["description"], desc); err != nil {
//			return err
//		}
//		return nil
//	}
//}
//
//func testAccSourceGroupConfig(domain string, name string, desc string) string {
//	return fmt.Sprintf(`
//resource "buddy_workspace" "foo" {
//    domain = "%s"
//}
//
//resource "buddy_group" "bar" {
//    domain = "${buddy_workspace.foo.domain}"
//    name = "%s"
//	description = "%s"
//}
//
//data "buddy_group" "id" {
//	domain = "${buddy_workspace.foo.domain}"
//	group_id = "${buddy_group.bar.group_id}"
//}
//
//data "buddy_group" "name" {
//	domain = "${buddy_workspace.foo.domain}"
//	name = "${buddy_group.bar.name}"
//}
//`, domain, name, desc)
//}
