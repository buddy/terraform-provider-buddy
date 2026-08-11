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
	envIdentifier := util.UniqueString()
	envName := util.RandString(10)
	name1 := "aaaa" + util.RandString(10)
	name2 := util.RandString(10)
	name3 := util.RandString(10)
	name4 := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceSandboxesConfig(domain, projectName, envIdentifier, envName, name1, name2, name3, name4),
				Check: resource.ComposeTestCheckFunc(
					// with neither filter set the API lists only the sandboxes sitting
					// directly in the workspace - it is not an "everything" listing, the
					// project and environment ones are excluded
					testAccSourceSandboxesAttributes("data.buddy_sandboxes.workspace", 1, name4),
					testAccSourceSandboxesAttributes("data.buddy_sandboxes.project", 2, ""),
					// the environment sits in the workspace, not in the project, so the
					// project filter does not reach it
					testAccSourceSandboxesAttributes("data.buddy_sandboxes.env", 1, name3),
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

func testAccSourceSandboxesConfig(domain string, projectName string, envIdentifier string, envName string, name1 string, name2 string, name3 string, name4 string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_project" "proj" {
   domain = "${buddy_workspace.foo.domain}"
   display_name = "%s"
}

resource "buddy_environment" "env" {
   domain = "${buddy_workspace.foo.domain}"
   identifier = "%s"
   name = "%s"
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

resource "buddy_sandbox" "c" {
   domain = "${buddy_workspace.foo.domain}"
   environment_id = "${buddy_environment.env.environment_id}"
   name = "%s"
}

resource "buddy_sandbox" "d" {
   domain = "${buddy_workspace.foo.domain}"
   name = "%s"
}

data "buddy_sandboxes" "workspace" {
   domain = "${buddy_workspace.foo.domain}"
   depends_on = [buddy_sandbox.a, buddy_sandbox.b, buddy_sandbox.c, buddy_sandbox.d]
}

data "buddy_sandboxes" "project" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   depends_on = [buddy_sandbox.a, buddy_sandbox.b, buddy_sandbox.c, buddy_sandbox.d]
}

data "buddy_sandboxes" "env" {
   domain = "${buddy_workspace.foo.domain}"
   environment_id = "${buddy_environment.env.environment_id}"
   depends_on = [buddy_sandbox.a, buddy_sandbox.b, buddy_sandbox.c, buddy_sandbox.d]
}

data "buddy_sandboxes" "name" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   name_regex = "^aaaa"
   depends_on = [buddy_sandbox.a, buddy_sandbox.b, buddy_sandbox.c, buddy_sandbox.d]
}
`, domain, projectName, envIdentifier, envName, name1, name2, name3, name4)
}
