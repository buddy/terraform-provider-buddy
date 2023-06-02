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

func TestAccSourceProjectGroup_upgrade(t *testing.T) {
	domain := util.UniqueString()
	groupName := util.RandString(10)
	projectName := util.UniqueString()
	permissionName := util.RandString(10)
	pipelineAccessLevel := buddy.PermissionAccessLevelRunOnly
	repoAccessLevel := buddy.PermissionAccessLevelReadWrite
	sandboxAccessLevel := buddy.PermissionAccessLevelReadWrite
	config := testAccSourceProjectGroupConfig(domain, groupName, projectName, permissionName, pipelineAccessLevel, repoAccessLevel, sandboxAccessLevel)
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
					testAccSourceProjectGroupAttributes("data.buddy_project_group.bar", groupName, permissionName, pipelineAccessLevel, repoAccessLevel, sandboxAccessLevel),
				),
			},
		},
	})
}

func TestAccSourceProjectGroup(t *testing.T) {
	domain := util.UniqueString()
	groupName := util.RandString(10)
	projectName := util.UniqueString()
	permissionName := util.RandString(10)
	pipelineAccessLevel := buddy.PermissionAccessLevelRunOnly
	repoAccessLevel := buddy.PermissionAccessLevelReadWrite
	sandboxAccessLevel := buddy.PermissionAccessLevelReadWrite
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceProjectGroupConfig(domain, groupName, projectName, permissionName, pipelineAccessLevel, repoAccessLevel, sandboxAccessLevel),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceProjectGroupAttributes("data.buddy_project_group.bar", groupName, permissionName, pipelineAccessLevel, repoAccessLevel, sandboxAccessLevel),
				),
			},
		},
	})
}

func testAccSourceProjectGroupAttributes(n string, groupName string, permissionName string, pipelineAccessLevel string, repoAccessLevel string, sandboxAccessLevel string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsPermissionPermissionId, _ := strconv.Atoi(attrs["permission.0.permission_id"])
		if err := util.CheckFieldSet("html_url", attrs["html_url"]); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], groupName); err != nil {
			return err
		}
		if err := util.CheckFieldSet("permission.0.html_url", attrs["permission.0.html_url"]); err != nil {
			return err
		}
		if err := util.CheckIntFieldSet("permission.0.permission_id", attrsPermissionPermissionId); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("permission.0.name", attrs["permission.0.name"], permissionName); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("permission.0.type", attrs["permission.0.type"], "CUSTOM"); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("permission.0.pipeline_access_level", attrs["permission.0.pipeline_access_level"], pipelineAccessLevel); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("permission.0.repository_access_level", attrs["permission.0.repository_access_level"], repoAccessLevel); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("permission.0.sandbox_access_level", attrs["permission.0.sandbox_access_level"], sandboxAccessLevel); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourceProjectGroupConfig(domain string, groupName string, projectName string, permissionName string, pipelineAccessLevel string, repoAccessLevel string, sandboxAccessLevel string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_project" "proj" {
   domain = "${buddy_workspace.foo.domain}"
   display_name = "%s"
}

resource "buddy_group" "gr" {
   domain = "${buddy_workspace.foo.domain}"
   name = "%s"
}

resource "buddy_permission" "perm" {
   domain = "${buddy_workspace.foo.domain}"
   name = "%s"
   pipeline_access_level = "%s"
   repository_access_level = "%s"
	sandbox_access_level = "%s"
}

resource "buddy_project_group" "bpg" {
   domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
	group_id = "${buddy_group.gr.group_id}"
	permission_id = "${buddy_permission.perm.permission_id}"
}

data "buddy_project_group" "bar" {
   domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.proj.name}"
	group_id = "${buddy_group.gr.group_id}"
   depends_on = [buddy_project_group.bpg]
}
`, domain, projectName, groupName, permissionName, pipelineAccessLevel, repoAccessLevel, sandboxAccessLevel)
}
