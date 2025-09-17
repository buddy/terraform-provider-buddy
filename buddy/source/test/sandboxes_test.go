package test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"
)

func TestAccSourceSandboxes(t *testing.T) {
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name1 := "aaaa" + util.RandString(10)
	name2 := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceSandboxesConfig(domain, projectName, name1, name2),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceSandboxesAttributes("data.buddy_sandboxes.all", 2, ""),
					testAccSourceSandboxesAttributes("data.buddy_sandboxes.name", 1, name1),
				),
			},
		},
	})
}

func testAccSourceSandboxesAttributes(n string, count int, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsSandboxesCount, _ := strconv.Atoi(attrs["sandboxes.#"])
		if err := util.CheckIntFieldEqual("sandboxes.#", attrsSandboxesCount, count); err != nil {
			return err
		}
		if count > 0 {
			if name != "" {
				if err := util.CheckFieldEqualAndSet("sandboxes.0.name", attrs["sandboxes.0.name"], name); err != nil {
					return err
				}
			} else {
				if err := util.CheckFieldSet("sandboxes.0.name", attrs["sandboxes.0.name"]); err != nil {
					return err
				}
			}
			if err := util.CheckFieldSet("sandboxes.0.html_url", attrs["sandboxes.0.html_url"]); err != nil {
				return err
			}
			if err := util.CheckFieldSet("sandboxes.0.identifier", attrs["sandboxes.0.identifier"]); err != nil {
				return err
			}
			if err := util.CheckFieldSet("sandboxes.0.sandbox_id", attrs["sandboxes.0.sandbox_id"]); err != nil {
				return err
			}
			if err := util.CheckFieldSet("sandboxes.0.status", attrs["sandboxes.0.status"]); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccSourceSandboxesConfig(domain string, projectName string, name1 string, name2 string) string {
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

resource "buddy_sandbox" "b" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   name = "%s"
}

data "buddy_sandboxes" "all" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   depends_on = [buddy_sandbox.a, buddy_sandbox.b]
}

data "buddy_sandboxes" "name" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   name_regex = "^aaaa"
   depends_on = [buddy_sandbox.a, buddy_sandbox.b]
}
`, domain, projectName, name1, name2)
}
