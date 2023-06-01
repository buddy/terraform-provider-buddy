package test

import (
	"buddy-terraform/buddy/acc"
	"buddy-terraform/buddy/util"
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
)

func TestAccSourceWebhook(t *testing.T) {
	domain := util.UniqueString()
	event := buddy.WebhookEventPush
	projectName := util.UniqueString()
	targetUrl := "https://127.0.0.1"
	secretKey := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceWebhookConfig(domain, event, projectName, targetUrl, secretKey),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceWebhookAttributes("data.buddy_webhook.id", targetUrl),
					testAccSourceWebhookAttributes("data.buddy_webhook.url", targetUrl),
				),
			},
		},
	})
}

func testAccSourceWebhookAttributes(n string, targetUrl string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsWebhookId, _ := strconv.Atoi(attrs["webhook_id"])
		if err := util.CheckFieldEqualAndSet("target_url", attrs["target_url"], targetUrl); err != nil {
			return err
		}
		if err := util.CheckIntFieldSet("webhook_id", attrsWebhookId); err != nil {
			return err
		}
		if err := util.CheckFieldSet("html_url", attrs["html_url"]); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourceWebhookConfig(domain string, event string, projectName string, targetUrl string, secretKey string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_project" "proj" {
   domain = "${buddy_workspace.foo.domain}"
   display_name = "%s"
}

resource "buddy_webhook" "web" {
   domain = "${buddy_workspace.foo.domain}"
   events = ["%s"]
   target_url = "%s"
   secret_key = "%s"
   projects = ["${buddy_project.proj.name}"]
}

data "buddy_webhook" "id" {
   domain = "${buddy_workspace.foo.domain}"
   webhook_id = "${buddy_webhook.web.webhook_id}"
}

data "buddy_webhook" "url" {
   domain = "${buddy_workspace.foo.domain}"
   target_url = "${buddy_webhook.web.target_url}"
}
`, domain, projectName, event, targetUrl, secretKey)
}
