package test

import (
	"buddy-terraform/buddy/acc"
	"buddy-terraform/buddy/util"
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccSourceIntegration(t *testing.T) {
	domain := util.UniqueString()
	name := util.RandString(10)
	typ := buddy.IntegrationTypeAmazon
	scope := buddy.IntegrationScopeAdmin
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceIntegrationConfig(domain, name, typ, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceIntegrationAttributes("data.buddy_integration.id", name, typ),
					testAccSourceIntegrationAttributes("data.buddy_integration.name", name, typ),
				),
			},
		},
	})
}

func testAccSourceIntegrationAttributes(n string, name string, typ string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("type", attrs["type"], typ); err != nil {
			return err
		}
		if err := util.CheckFieldSet("integration_id", attrs["integration_id"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("html_url", attrs["html_url"]); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourceIntegrationConfig(domain string, name string, typ string, scope string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_integration" "int" {
   domain = "${buddy_workspace.foo.domain}"
   name = "%s"
   type = "%s"
   scope = "%s"
   access_key = "ABC1234567890"
   secret_key = "ABC1234567890"
}

data "buddy_integration" "id" {
   domain = "${buddy_workspace.foo.domain}"
   integration_id = "${buddy_integration.int.integration_id}"
}

data "buddy_integration" "name" {
   domain = "${buddy_workspace.foo.domain}"
   name = "${buddy_integration.int.name}"
}
`, domain, name, typ, scope)
}
