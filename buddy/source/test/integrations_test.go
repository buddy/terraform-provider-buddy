package test

import (
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"
)

func TestAccSourceIntegrations(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceIntegrationsConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceIntegrationsAttributes("data.buddy_integrations.all", 2),
					testAccSourceIntegrationsAttributes("data.buddy_integrations.name", 1),
					testAccSourceIntegrationsAttributes("data.buddy_integrations.type", 1),
				),
			},
		},
	})
}

func testAccSourceIntegrationsAttributes(n string, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsIntegrationsCount, _ := strconv.Atoi(attrs["integrations.#"])
		if err := util.CheckIntFieldEqualAndSet("integrations.#", attrsIntegrationsCount, count); err != nil {
			return err
		}
		if err := util.CheckFieldSet("integrations.0.html_url", attrs["integrations.0.html_url"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("integrations.0.integration_id", attrs["integrations.0.integration_id"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("integrations.0.name", attrs["integrations.0.name"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("integrations.0.type", attrs["integrations.0.type"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("integrations.0.identifier", attrs["integrations.0.identifier"]); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourceIntegrationsConfig() string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_integration" "a" {
   domain = "${buddy_workspace.foo.domain}"
   name = "abcdef"
   type = "%s"
   scope = "%s"
   access_key = "ABC1234567890"
   secret_key = "ABC1234567890"
}

resource "buddy_integration" "b" {
   domain = "${buddy_workspace.foo.domain}"
   name = "zzzzz"
   type = "%s"
   scope = "%s"
   token = "abcdefghijklmnoprst"
}

data "buddy_integrations" "all" {
   domain = "${buddy_workspace.foo.domain}"
   depends_on = [buddy_integration.a, buddy_integration.b]
}

data "buddy_integrations" "name" {
   domain = "${buddy_workspace.foo.domain}"
   depends_on = [buddy_integration.a, buddy_integration.b]
   name_regex = "^abc"
}

data "buddy_integrations" "type" {
   domain = "${buddy_workspace.foo.domain}"
   type = "AMAZON"
   depends_on = [buddy_integration.a, buddy_integration.b]
}
`, util.UniqueString(), buddy.IntegrationTypeAmazon, buddy.IntegrationScopeWorkspace, buddy.IntegrationTypeDigitalOcean, buddy.IntegrationScopeWorkspace)
}
