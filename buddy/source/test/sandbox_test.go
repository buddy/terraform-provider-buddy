package test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"
)

func TestAccSourceSandbox(t *testing.T) {
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceSandboxConfig(domain, projectName, name),
				Check:  testAccSourceSandboxAttributes("data.buddy_sandbox.a", name),
			},
		},
	})
}

func testAccSourceSandboxAttributes(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], name); err != nil {
			return err
		}
		if err := util.CheckFieldSet("html_url", attrs["html_url"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("identifier", attrs["identifier"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("sandbox_id", attrs["sandbox_id"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("status", attrs["status"]); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourceSandboxConfig(domain string, projectName string, name string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_project" "proj" {
   domain = "${buddy_workspace.foo.domain}"
   display_name = "%s"
}

resource "buddy_sandbox" "a" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   name = "%s"
}

data "buddy_sandbox" "a" {
   domain = "${buddy_workspace.foo.domain}"
   sandbox_id = "${buddy_sandbox.a.sandbox_id}"
}
`, domain, projectName, name)
}
