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
//func TestAccSourceWebhooks(t *testing.T) {
//	domain := util.UniqueString()
//	target1 := "https://127.0.0.1"
//	target2 := "https://192.168.1.1"
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			acc.PreCheck(t)
//		},
//		CheckDestroy:      acc.DummyCheckDestroy,
//		ProviderFactories: acc.ProviderFactories,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccSourceWebhooksConfig(domain, target1, target2),
//				Check: resource.ComposeTestCheckFunc(
//					testAccSourceWebhooksAttributes("data.buddy_webhooks.all", 2),
//					testAccSourceWebhooksAttributes("data.buddy_webhooks.filter", 1),
//				),
//			},
//		},
//	})
//}
//
//func testAccSourceWebhooksAttributes(n string, count int) resource.TestCheckFunc {
//	return func(s *terraform.State) error {
//		rs, ok := s.RootModule().Resources[n]
//		if !ok {
//			return fmt.Errorf("not found: %s", n)
//		}
//		attrs := rs.Primary.Attributes
//		attrsWebhooksCount, _ := strconv.Atoi(attrs["webhooks.#"])
//		attrsWebhookId, _ := strconv.Atoi(attrs["webhooks.0.webhook_id"])
//		if err := util.CheckIntFieldEqualAndSet("webhooks.#", attrsWebhooksCount, count); err != nil {
//			return err
//		}
//		if err := util.CheckFieldSet("webhooks.0.html_url", attrs["webhooks.0.html_url"]); err != nil {
//			return err
//		}
//		if err := util.CheckFieldSet("webhooks.0.target_url", attrs["webhooks.0.target_url"]); err != nil {
//			return err
//		}
//		if err := util.CheckIntFieldSet("webhooks.0.webhook_id", attrsWebhookId); err != nil {
//			return err
//		}
//		return nil
//	}
//}
//
//func testAccSourceWebhooksConfig(domain string, target1 string, target2 string) string {
//	return fmt.Sprintf(`
//resource "buddy_workspace" "foo" {
//    domain = "%s"
//}
//
//resource "buddy_webhook" "a" {
//    domain = "${buddy_workspace.foo.domain}"
//    events = ["PUSH"]
//	projects = []
//    target_url = "%s"
//}
//
//resource "buddy_webhook" "b" {
//    domain = "${buddy_workspace.foo.domain}"
//    events = ["PUSH"]
//	projects = []
//    target_url = "%s"
//}
//
//data "buddy_webhooks" "all" {
//    domain = "${buddy_workspace.foo.domain}"
//    depends_on = [buddy_webhook.a, buddy_webhook.b]
//}
//
//data "buddy_webhooks" "filter" {
//    domain = "${buddy_workspace.foo.domain}"
//    depends_on = [buddy_webhook.a, buddy_webhook.b]
//    target_url_regex = "192"
//}
//
//`, domain, target1, target2)
//}
