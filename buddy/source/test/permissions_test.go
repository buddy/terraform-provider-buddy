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

func TestAccSourcePermissions_upgrade(t *testing.T) {
	domain := util.UniqueString()
	config := testAccSourcePermissionsConfig(domain)
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"buddy": {
						VersionConstraint: "1.12.0",
						Source:            "buddy/buddy",
					},
				},
				Config: config,
			},
			{
				ProtoV6ProviderFactories: acc.ProviderFactories,
				Config:                   config,
				Check: resource.ComposeTestCheckFunc(
					testAccSourcePermissionsAttributes("data.buddy_permissions.all", 4),
					testAccSourcePermissionsAttributes("data.buddy_permissions.name", 1),
					testAccSourcePermissionsAttributes("data.buddy_permissions.type", 1),
				),
			},
		},
	})
}

func TestAccSourcePermissions(t *testing.T) {
	domain := util.UniqueString()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourcePermissionsConfig(domain),
				Check: resource.ComposeTestCheckFunc(
					testAccSourcePermissionsAttributes("data.buddy_permissions.all", 4),
					testAccSourcePermissionsAttributes("data.buddy_permissions.name", 1),
					testAccSourcePermissionsAttributes("data.buddy_permissions.type", 1),
				),
			},
		},
	})
}

func testAccSourcePermissionsAttributes(n string, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsPermissionsCount, _ := strconv.Atoi(attrs["permissions.#"])
		attrsPermissionId, _ := strconv.Atoi(attrs["permissions.0.permission_id"])
		if err := util.CheckIntFieldEqualAndSet("permissions.#", attrsPermissionsCount, count); err != nil {
			return err
		}
		if err := util.CheckFieldSet("permissions.0.name", attrs["permissions.0.name"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("permissions.0.pipeline_access_level", attrs["permissions.0.pipeline_access_level"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("permissions.0.sandbox_access_level", attrs["permissions.0.sandbox_access_level"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("permissions.0.repository_access_level", attrs["permissions.0.repository_access_level"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("permissions.0.html_url", attrs["permissions.0.html_url"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("permissions.0.type", attrs["permissions.0.type"]); err != nil {
			return err
		}
		if err := util.CheckIntFieldSet("permissions.0.permission_id", attrsPermissionId); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourcePermissionsConfig(domain string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_permission" "perm" {
   domain = "${buddy_workspace.foo.domain}"
   name = "abcdef"
   pipeline_access_level = "%s"
   repository_access_level = "%s"
	sandbox_access_level = "%s"
	project_team_access_level = "%s"
}

data "buddy_permissions" "all" {
   domain = "${buddy_workspace.foo.domain}"
   depends_on = [buddy_permission.perm]
}

data "buddy_permissions" "name" {
   domain = "${buddy_workspace.foo.domain}"
   depends_on = [buddy_permission.perm]
   name_regex = "^abc"
}

data "buddy_permissions" "type" {
   domain = "${buddy_workspace.foo.domain}"
   depends_on = [buddy_permission.perm]
   type = "DEVELOPER"
}
`, domain, buddy.PermissionAccessLevelReadWrite, buddy.PermissionAccessLevelManage, buddy.PermissionAccessLevelReadWrite, buddy.PermissionAccessLevelManage)
}
