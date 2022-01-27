package test

import (
	"buddy-terraform/buddy/acc"
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strconv"
	"testing"
)

func TestAccWebhook(t *testing.T) {
	var webhook api.Webhook
	domain := util.UniqueString()
	event := api.WebhookEventPush
	newEvent := api.WebhookEventExecutionSuccessful
	projectName := util.UniqueString()
	targetUrl := "https://127.0.0.1"
	newTargetUrl := "https://aaaa.com"
	var secretKey string
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccWebhookCheckDestroy,
		Steps: []resource.TestStep{
			// create webhook
			{
				Config: testAccWebhookConfig(domain, projectName, event, targetUrl, secretKey),
				Check: resource.ComposeTestCheckFunc(
					testAccWebhookGet("buddy_webhook.bar", &webhook),
					testAccWebhookAttributes("buddy_webhook.bar", &webhook, projectName, event, targetUrl, secretKey),
				),
			},
			// edit webhook
			{
				Config: testAccWebhookUpdateConfig(domain, projectName, newEvent, newTargetUrl),
				Check: resource.ComposeTestCheckFunc(
					testAccWebhookGet("buddy_webhook.bar", &webhook),
					testAccWebhookAttributes("buddy_webhook.bar", &webhook, "", newEvent, newTargetUrl, ""),
				),
			},
			// import webhook
			{
				ResourceName:      "buddy_webhook.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccWebhookAttributes(n string, webhook *api.Webhook, projectName string, event string, targetUrl string, secretKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsWebhookId, _ := strconv.Atoi(attrs["webhook_id"])
		if err := util.CheckFieldEqualAndSet("TargetUrl", webhook.TargetUrl, targetUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqual("SecretKey", webhook.SecretKey, secretKey); err != nil {
			return err
		}
		if len(webhook.Events) != 1 {
			return fmt.Errorf("expected \"Events\" to have one element")
		}
		if err := util.CheckFieldEqualAndSet("Events", webhook.Events[0], event); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("events.0", attrs["events.0"], event); err != nil {
			return err
		}
		if projectName == "" {
			if len(webhook.Projects) != 0 {
				return fmt.Errorf("expected \"Projects\" to have 0 elements")
			}
			if _, ok := attrs["projects.0"]; ok {
				return fmt.Errorf("expected \"projects\" to have 0 elements")
			}
		} else {
			if len(webhook.Projects) != 1 {
				return fmt.Errorf("expected \"Projects\" to have 1 element")
			}
			if err := util.CheckFieldEqualAndSet("Projects", webhook.Projects[0], projectName); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("projects.0", attrs["projects.0"], projectName); err != nil {
				return err
			}
		}
		if err := util.CheckFieldEqualAndSet("target_url", attrs["target_url"], targetUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqual("secret_key", attrs["secret_key"], secretKey); err != nil {
			return err
		}
		if err := util.CheckIntFieldEqualAndSet("webhook_id", attrsWebhookId, webhook.Id); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("html_url", attrs["html_url"], webhook.HtmlUrl); err != nil {
			return err
		}
		return nil
	}
}

func testAccWebhookGet(n string, webhook *api.Webhook) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		domain, wid, err := util.DecomposeDoubleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		webhookId, err := strconv.Atoi(wid)
		if err != nil {
			return err
		}
		w, _, err := acc.ApiClient.WebhookService.Get(domain, webhookId)
		if err != nil {
			return err
		}
		*webhook = *w
		return nil
	}
}

func testAccWebhookConfig(domain string, projectName string, event string, targetUrl string, secretKey string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_webhook" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    events = ["%s"]
    target_url = "%s"
    secret_key = "%s"
    projects = ["${buddy_project.proj.name}"]
}

`, domain, projectName, event, targetUrl, secretKey)
}

func testAccWebhookUpdateConfig(domain string, projectName string, event string, targetUrl string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_webhook" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    events = ["%s"]
	projects = []
    target_url = "%s"
}
`, domain, projectName, event, targetUrl)
}

func testAccWebhookCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_webhook" {
			continue
		}
		domain, wid, err := util.DecomposeDoubleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		webhookId, err := strconv.Atoi(wid)
		if err != nil {
			return err
		}
		webhook, resp, err := acc.ApiClient.WebhookService.Get(domain, webhookId)
		if err == nil && webhook != nil {
			return util.ErrorResourceExists()
		}
		if resp.StatusCode != 404 {
			return err
		}
	}
	return nil
}