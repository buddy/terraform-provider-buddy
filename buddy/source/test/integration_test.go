package test

import (
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"
)

func TestAccSourceIntegration(t *testing.T) {
	domain := util.UniqueString()
	name := util.RandString(10)
	typ := buddy.IntegrationTypeAmazon
	scope := buddy.IntegrationScopeWorkspace
	identifier := util.RandString(10)
	note := util.RandString(10)
	agentNote := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceIntegrationConfig(domain, name, typ, scope, identifier, note, agentNote),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceIntegrationAttributes("data.buddy_integration.id", name, typ, identifier, note, agentNote),
					testAccSourceIntegrationAttributes("data.buddy_integration.name", name, typ, identifier, note, agentNote),
				),
			},
		},
	})
}

func TestAccSourceIntegration_project(t *testing.T) {
	domain := util.UniqueString()
	projectDisplayName := util.RandString(10)
	name := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceIntegrationProjectConfig(domain, projectDisplayName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceIntegrationProjectAttributes("data.buddy_integration.id", name),
					testAccSourceIntegrationProjectAttributes("data.buddy_integration.name", name),
				),
			},
		},
	})
}

func testAccSourceIntegrationProjectAttributes(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		projectName := s.RootModule().Resources["buddy_project.proj"].Primary.Attributes["name"]
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("type", attrs["type"], buddy.IntegrationTypeShopify); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("project_name", attrs["project_name"], projectName); err != nil {
			return err
		}
		if err := util.CheckFieldSet("integration_id", attrs["integration_id"]); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourceIntegrationAttributes(n string, name string, typ string, identifier string, note string, agentNote string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("note", attrs["note"], note); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("agent_note", attrs["agent_note"], agentNote); err != nil {
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
		if err := util.CheckFieldEqualAndSet("identifier", attrs["identifier"], identifier); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourceIntegrationConfig(domain string, name string, typ string, scope string, identifier string, note string, agentNote string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_integration" "int" {
   domain = "${buddy_workspace.foo.domain}"
   name = "%s"
   type = "%s"
   scope = "%s"
   identifier = "%s"
   access_key = "ABC1234567890"
   secret_key = "ABC1234567890"
   note = "%s"
   agent_note = "%s"
}

data "buddy_integration" "id" {
   domain = "${buddy_workspace.foo.domain}"
   integration_id = "${buddy_integration.int.integration_id}"
}

data "buddy_integration" "name" {
   domain = "${buddy_workspace.foo.domain}"
   name = "${buddy_integration.int.name}"
}
`, domain, name, typ, scope, identifier, note, agentNote)
}

func testAccSourceIntegrationProjectConfig(domain string, projectDisplayName string, name string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_project" "proj" {
   domain = "${buddy_workspace.foo.domain}"
   display_name = "%s"
}

resource "buddy_integration" "int" {
   domain = "${buddy_workspace.foo.domain}"
   name = "%s"
   type = "%s"
   scope = "%s"
   project_name = "${buddy_project.proj.name}"
   shop = "ABC"
   token = "abcdefghijklmnoprst"
}

data "buddy_integration" "id" {
   domain = "${buddy_workspace.foo.domain}"
   integration_id = "${buddy_integration.int.integration_id}"
}

data "buddy_integration" "name" {
   domain = "${buddy_workspace.foo.domain}"
   name = "${buddy_integration.int.name}"
   project_name = "${buddy_project.proj.name}"
}
`, domain, projectDisplayName, name, buddy.IntegrationTypeShopify, buddy.IntegrationScopeProject)
}
