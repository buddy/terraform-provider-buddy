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
//func TestAccSourceVariable(t *testing.T) {
//	domain := util.UniqueString()
//	key := util.RandString(10)
//	val := util.RandString(10)
//	desc := util.RandString(10)
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			acc.PreCheck(t)
//		},
//		CheckDestroy:      acc.DummyCheckDestroy,
//		ProviderFactories: acc.ProviderFactories,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccSourceVariableConfig(domain, key, val, desc, false, true),
//				Check: resource.ComposeTestCheckFunc(
//					testAccSourceVariableAttributes("data.buddy_variable.id", key, val, desc, false, true),
//					testAccSourceVariableAttributes("data.buddy_variable.key", key, val, desc, false, true),
//				),
//			},
//		},
//	})
//}
//
//func testAccSourceVariableAttributes(n string, key string, val string, desc string, encrypred bool, settable bool) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		rs, ok := s.RootModule().Resources[n]
//		if !ok {
//			return fmt.Errorf("not found: %s", n)
//		}
//		attrs := rs.Primary.Attributes
//		attrsVariableId, _ := strconv.Atoi(attrs["variable_id"])
//		attrsEncrypted, _ := strconv.ParseBool(attrs["encrypted"])
//		attrsSettable, _ := strconv.ParseBool(attrs["settable"])
//		if err := util.CheckFieldEqualAndSet("key", attrs["key"], key); err != nil {
//			return err
//		}
//		if err := util.CheckBoolFieldEqual("encrypted", attrsEncrypted, encrypred); err != nil {
//			return err
//		}
//		if err := util.CheckBoolFieldEqual("settable", attrsSettable, settable); err != nil {
//			return err
//		}
//		if err := util.CheckFieldEqualAndSet("description", attrs["description"], desc); err != nil {
//			return err
//		}
//		if err := util.CheckFieldEqualAndSet("value", attrs["value"], val); err != nil {
//			return err
//		}
//		if err := util.CheckIntFieldSet("variable_id", attrsVariableId); err != nil {
//			return err
//		}
//		return nil
//	}
//}
//
//func testAccSourceVariableConfig(domain string, key string, val string, desc string, encrypred bool, settable bool) string {
//	return fmt.Sprintf(`
//resource "buddy_workspace" "foo" {
//    domain = "%s"
//}
//
//resource "buddy_variable" "var" {
//    domain = "${buddy_workspace.foo.domain}"
//    key = "%s"
//    value = "%s"
//	encrypted = %s
//	settable = %s
//	description = "%s"
//}
//
//data "buddy_variable" "id" {
//    domain = "${buddy_workspace.foo.domain}"
//    variable_id = "${buddy_variable.var.variable_id}"
//}
//
//data "buddy_variable" "key" {
//    domain = "${buddy_workspace.foo.domain}"
//    key = "${buddy_variable.var.key}"
//}
//`, domain, key, val, strconv.FormatBool(encrypred), strconv.FormatBool(settable), desc)
//}
