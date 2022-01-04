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

func TestAccSourceWorkspaces(t *testing.T) {
	domain1 := util.UniqueString() + "aaa"
	domain2 := "bbb" + util.UniqueString()
	domain3 := util.UniqueString()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:      acc.DummyCheckDestroy,
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceWorkspacesConfig(domain1, domain2, domain3),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceWorkspacesAttributes("data.buddy_workspaces.a", 1, domain1),
					testAccSourceWorkspacesAttributes("data.buddy_workspaces.b", 1, domain2),
				),
			},
		},
	})
}

func testAccSourceWorkspacesAttributes(n string, count int, domain string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsWorkspacesCount, _ := strconv.Atoi(attrs["workspaces.#"])
		attrsWorkspaceId, _ := strconv.Atoi(attrs["workspaces.0.workspace_id"])
		if err := util.CheckIntFieldEqualAndSet("workspaces.#", attrsWorkspacesCount, count); err != nil {
			return err
		}
		if domain != "" {
			if err := util.CheckFieldEqualAndSet("workspaces.0.domain", attrs["workspaces.0.domain"], domain); err != nil {
				return err
			}
		} else {
			if err := util.CheckFieldSet("workspaces.0.domain", attrs["workspaces.0.domain"]); err != nil {
				return err
			}
		}
		if err := util.CheckFieldSet("workspaces.0.html_url", attrs["workspaces.0.html_url"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("workspaces.0.name", attrs["workspaces.0.name"]); err != nil {
			return err
		}
		if err := util.CheckIntFieldSet("workspaces.0.workspace_id", attrsWorkspaceId); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourceWorkspacesConfig(domain1 string, domain2 string, domain3 string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "a" {
    domain = "%s"
}

resource "buddy_workspace" "b" {
    domain = "%s"
}

resource "buddy_workspace" "c" {
    domain = "%s"
}

data "buddy_workspaces" "a" {
    depends_on = [buddy_workspace.a, buddy_workspace.b, buddy_workspace.c]
    name_regex = "aaa$"
}

data "buddy_workspaces" "b" {
    depends_on = [buddy_workspace.a, buddy_workspace.b, buddy_workspace.c]
    domain_regex = "^bbb"
}

`, domain1, domain2, domain3)
}
