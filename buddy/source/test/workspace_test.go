package test

import (
	"buddy-terraform/buddy/acc"
	"buddy-terraform/buddy/util"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strconv"
	"testing"
)

func TestAccSourceWorkspace(t *testing.T) {
	domain := util.UniqueString()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:      acc.DummyCheckDestroy,
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceWorkspaceConfig(domain),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceWorkspaceAttributes("data.buddy_workspace.domain", domain),
					testAccSourceWorkspaceAttributes("data.buddy_workspace.name", domain),
				),
			},
		},
	})
}

func testAccSourceWorkspaceAttributes(n string, domain string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsWorkspaceId, _ := strconv.Atoi(attrs["workspace_id"])
		if err := util.CheckFieldEqualAndSet("domain", attrs["domain"], domain); err != nil {
			return err
		}
		if err := util.CheckIntFieldSet("workspace_id", attrsWorkspaceId); err != nil {
			return err
		}
		if err := util.CheckFieldSet("html_url", attrs["html_url"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("name", attrs["name"]); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourceWorkspaceConfig(domain string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

data "buddy_workspace" "domain" {
    domain = "${buddy_workspace.foo.domain}"
}

data "buddy_workspace" "name" {
    name = "${buddy_workspace.foo.name}"
}
`, domain)
}
