package test

import (
	"buddy-terraform/buddy/acc"
	"buddy-terraform/buddy/util"
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strconv"
	"testing"
)

func TestAccWorkspace(t *testing.T) {
	var workspace buddy.Workspace
	domain := util.UniqueString()
	salt := util.RandString(10)
	name := "A" + util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccWorkspaceCheckDestroy,
		Steps: []resource.TestStep{
			// create
			{
				Config: testAccWorkspaceConfig(domain),
				Check: resource.ComposeTestCheckFunc(
					testAccWorkspaceGet("buddy_workspace.foo", &workspace),
					testAccWorkspaceAttributes("buddy_workspace.foo", &workspace, domain, ""),
				),
			},
			// update
			{
				Config: testAccWorkspaceConfigFull(domain, salt, name),
				Check: resource.ComposeTestCheckFunc(
					testAccWorkspaceGet("buddy_workspace.foo", &workspace),
					testAccWorkspaceAttributes("buddy_workspace.foo", &workspace, domain, name),
				),
			},
			{
				ResourceName:            "buddy_workspace.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"encryption_salt"},
			},
		},
	})
}

func testAccWorkspaceCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_workspace" {
			continue
		}
		workspace, resp, err := acc.ApiClient.WorkspaceService.Get(rs.Primary.ID)
		if err == nil && workspace != nil {
			return util.ErrorResourceExists()
		}
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}

func testAccWorkspaceAttributes(n string, workspace *buddy.Workspace, domain string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsWorkspaceId, _ := strconv.Atoi(attrs["workspace_id"])
		attrsOwnerId, _ := strconv.Atoi(attrs["owner_id"])
		attrsFrozen, _ := strconv.ParseBool(attrs["frozen"])
		if err := util.CheckFieldEqualAndSet("Domain", workspace.Domain, domain); err != nil {
			return err
		}
		if err := util.CheckIntFieldEqualAndSet("workspace_id", attrsWorkspaceId, workspace.Id); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("html_url", attrs["html_url"], workspace.HtmlUrl); err != nil {
			return err
		}
		if err := util.CheckIntFieldEqualAndSet("owner_id", attrsOwnerId, workspace.OwnerId); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], workspace.Name); err != nil {
			return err
		}
		if name != "" {
			if err := util.CheckFieldEqualAndSet("Name", workspace.Name, name); err != nil {
				return err
			}
		}
		if err := util.CheckBoolFieldEqual("frozen", attrsFrozen, workspace.Frozen); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("create_date", attrs["create_date"], workspace.CreateDate); err != nil {
			return err
		}
		return nil
	}
}

func testAccWorkspaceGet(n string, workspace *buddy.Workspace) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		w, _, err := acc.ApiClient.WorkspaceService.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		*workspace = *w
		return nil
	}
}

func testAccWorkspaceConfig(domain string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}`, domain)
}

func testAccWorkspaceConfigFull(domain string, salt string, name string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
	name = "%s"
	encryption_salt = "%s"
}`, domain, name, salt)
}
