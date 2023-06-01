package test

import (
	"buddy-terraform/buddy/acc"
	"buddy-terraform/buddy/util"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
)

func TestAccSourceVariables(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceVariablesConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceVariablesAttributes("data.buddy_variables.all", 2),
					testAccSourceVariablesAttributes("data.buddy_variables.key", 1),
					testAccSourceVariablesAttributes("data.buddy_variables.project", 1),
				),
			},
		},
	})
}

func testAccSourceVariablesAttributes(n string, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsVariablesCount, _ := strconv.Atoi(attrs["variables.#"])
		attrsVariableId, _ := strconv.Atoi(attrs["variables.0.variable_id"])
		attrsSettable, _ := strconv.ParseBool(attrs["variables.0.settable"])
		attrsEncrypted, _ := strconv.ParseBool(attrs["variables.0.encrypted"])
		if err := util.CheckIntFieldEqualAndSet("variables.#", attrsVariablesCount, count); err != nil {
			return err
		}
		if err := util.CheckIntFieldSet("variables.0.variable_id", attrsVariableId); err != nil {
			return err
		}
		if err := util.CheckFieldSet("variables.0.key", attrs["variables.0.key"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("variables.0.description", attrs["variables.0.description"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("variables.0.value", attrs["variables.0.value"]); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("variables.0.settable", attrsSettable, true); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("variables.0.encrypted", attrsEncrypted, true); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourceVariablesConfig() string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_variable" "a" {
   domain = "${buddy_workspace.foo.domain}"
   key = "abcdef"
   value = "abcdef"
	encrypted = true
	settable = true
	description = "abcdef"
}

resource "buddy_variable" "aa" {
   domain = "${buddy_workspace.foo.domain}"
   key = "ueteryw"
   value = "ueteryw"
	encrypted = true
	settable = true
	description = "ueteryw"
}

resource "buddy_project" "p" {
   domain = "${buddy_workspace.foo.domain}"
   display_name = "abcdef"
}

resource "buddy_variable" "b" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.p.name}"
   key = "test"
   value = "test"
	encrypted = true
	settable = true
	description = "test"
}

data "buddy_variables" "all" {
   domain = "${buddy_workspace.foo.domain}"
   depends_on = [buddy_variable.a, buddy_variable.b, buddy_variable.aa]
}

data "buddy_variables" "key" {
   domain = "${buddy_workspace.foo.domain}"
   depends_on = [buddy_variable.a, buddy_variable.b, buddy_variable.aa]
   key_regex = "^abc"
}

data "buddy_variables" "project" {
   domain = "${buddy_workspace.foo.domain}"
   depends_on = [buddy_variable.a, buddy_variable.b, buddy_variable.aa]
   project_name = "abcdef"
   key_regex = "^te"
}
`, util.UniqueString())
}
