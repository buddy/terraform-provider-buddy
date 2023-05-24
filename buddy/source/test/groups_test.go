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
//func TestAccSourceGroups(t *testing.T) {
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			acc.PreCheck(t)
//		},
//		ProviderFactories: acc.ProviderFactories,
//		CheckDestroy:      acc.DummyCheckDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccSourceGroupsConfig(),
//				Check: resource.ComposeTestCheckFunc(
//					testAccSourceGroupsAttributes("data.buddy_groups.all", 2),
//					testAccSourceGroupsAttributes("data.buddy_groups.name", 1),
//				),
//			},
//		},
//	})
//}
//
//func testAccSourceGroupsAttributes(n string, count int) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		rs, ok := s.RootModule().Resources[n]
//		if !ok {
//			return fmt.Errorf("not found: %s", n)
//		}
//		attrs := rs.Primary.Attributes
//		attrsGroupsCount, _ := strconv.Atoi(attrs["groups.#"])
//		attrsGroupId, _ := strconv.Atoi(attrs["groups.0.group_id"])
//		if err := util.CheckIntFieldEqual("groups.#", attrsGroupsCount, count); err != nil {
//			return err
//		}
//		if err := util.CheckIntFieldSet("groups.0.group_id", attrsGroupId); err != nil {
//			return err
//		}
//		if err := util.CheckFieldSet("groups.0.name", attrs["groups.0.name"]); err != nil {
//			return err
//		}
//		if err := util.CheckFieldSet("groups.0.html_url", attrs["groups.0.html_url"]); err != nil {
//			return err
//		}
//		return nil
//	}
//}
//
//func testAccSourceGroupsConfig() string {
//	return fmt.Sprintf(`
//resource "buddy_workspace" "w" {
//    domain = "%s"
//}
//
//resource "buddy_group" "a" {
//    domain = "${buddy_workspace.w.domain}"
//    name = "abcdef"
//}
//
//resource "buddy_group" "b" {
//    domain = "${buddy_workspace.w.domain}"
//    name = "test"
//}
//
//data "buddy_groups" "all" {
//    domain = "${buddy_workspace.w.domain}"
//    depends_on = [buddy_group.a, buddy_group.b]
//}
//
//data "buddy_groups" "name" {
//    domain = "${buddy_workspace.w.domain}"
//    depends_on = [buddy_group.a, buddy_group.b]
//    name_regex = "^abc"
//}
//`, util.UniqueString())
//}
