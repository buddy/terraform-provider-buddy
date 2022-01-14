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

func TestAccSourceVariablesSshKeys(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:      acc.DummyCheckDestroy,
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceVariablesSshKeysConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceVariablesSshKeysAttributes("data.buddy_variables_ssh_keys.all", 2),
					testAccSourceVariablesSshKeysAttributes("data.buddy_variables_ssh_keys.key", 1),
					testAccSourceVariablesSshKeysAttributes("data.buddy_variables_ssh_keys.project", 1),
				),
			},
		},
	})
}

func testAccSourceVariablesSshKeysAttributes(n string, count int) resource.TestCheckFunc {
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
		if err := util.CheckBoolFieldEqual("variables.0.settable", attrsSettable, false); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("variables.0.encrypted", attrsEncrypted, true); err != nil {
			return err
		}
		if err := util.CheckFieldSet("variables.0.file_name", attrs["variables.0.file_name"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("variables.0.file_place", attrs["variables.0.file_place"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("variables.0.file_path", attrs["variables.0.file_path"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("variables.0.file_chmod", attrs["variables.0.file_chmod"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("variables.0.checksum", attrs["variables.0.checksum"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("variables.0.key_fingerprint", attrs["variables.0.key_fingerprint"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("variables.0.public_value", attrs["variables.0.public_value"]); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourceVariablesSshKeysConfig() string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_project" "p" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "abcdef"
}

resource "buddy_variable_ssh_key" "a" {
    domain = "${buddy_workspace.foo.domain}"
    key = "abcdef"
    description = "abcdef"
    file_place = "CONTAINER"
    file_name = "abcdef"
    file_path = "~/abcdef"
    file_chmod = "600"
    value = <<EOT
%s
EOT
}

resource "buddy_variable_ssh_key" "b" {
    domain = "${buddy_workspace.foo.domain}"
    key = "test"
    description = "test"
    file_place = "CONTAINER"
    file_name = "test"
    file_path = "~/test"
    file_chmod = "600"
    value = <<EOT
%s
EOT
}

resource "buddy_variable_ssh_key" "c" {
    domain = "${buddy_workspace.foo.domain}"
    project_name = "${buddy_project.p.name}"
    key = "test"
    description = "test"
    file_place = "CONTAINER"
    file_name = "test"
    file_path = "~/test"
    file_chmod = "600"
    value = <<EOT
%s
EOT
}

data "buddy_variables_ssh_keys" "all" {
    domain = "${buddy_workspace.foo.domain}"
    depends_on = [buddy_variable_ssh_key.a, buddy_variable_ssh_key.b, buddy_variable_ssh_key.c]
}

data "buddy_variables_ssh_keys" "key" {
    domain = "${buddy_workspace.foo.domain}"
    depends_on = [buddy_variable_ssh_key.a, buddy_variable_ssh_key.b, buddy_variable_ssh_key.c]
    key_regex = "^abc" 
}

data "buddy_variables_ssh_keys" "project" {
    domain = "${buddy_workspace.foo.domain}"
    depends_on = [buddy_variable_ssh_key.a, buddy_variable_ssh_key.b, buddy_variable_ssh_key.c]
    project_name = "${buddy_project.p.name}"
    key_regex = "^te" 
}

`, util.UniqueString(), util.SshKey, util.SshKey2, util.SshKey)
}
