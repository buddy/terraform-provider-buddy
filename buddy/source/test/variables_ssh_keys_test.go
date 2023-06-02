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

func TestAccSourceVariablesSshKeys(t *testing.T) {
	err, _, privateKey := util.GenerateRsaKeyPair()
	if err != nil {
		t.Fatal(err.Error())
	}
	err, _, privateKey2 := util.GenerateRsaKeyPair()
	if err != nil {
		t.Fatal(err.Error())
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceVariablesSshKeysConfig(privateKey, privateKey2),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceVariablesSshKeysAttributes("data.buddy_variables_ssh_keys.all", 3),
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
		if err := util.CheckIntFieldEqualAndSet("variables.#", attrsVariablesCount, count); err != nil {
			return err
		}
		for i := 0; i < count; i += 1 {
			index := "variables." + strconv.Itoa(i)
			attrsVariableId, _ := strconv.Atoi(attrs[index+".variable_id"])
			attrsSettable, _ := strconv.ParseBool(attrs[index+".settable"])
			attrsEncrypted, _ := strconv.ParseBool(attrs[index+".encrypted"])
			if err := util.CheckIntFieldSet(index+".variable_id", attrsVariableId); err != nil {
				return err
			}
			if err := util.CheckFieldSet(index+".key", attrs[index+".key"]); err != nil {
				return err
			}
			if err := util.CheckFieldSet(index+".value", attrs[index+".value"]); err != nil {
				return err
			}
			if err := util.CheckBoolFieldEqual(index+".settable", attrsSettable, false); err != nil {
				return err
			}
			if err := util.CheckBoolFieldEqual(index+".encrypted", attrsEncrypted, true); err != nil {
				return err
			}
			if err := util.CheckFieldSet(index+".file_place", attrs[index+".file_place"]); err != nil {
				return err
			}
			if err := util.CheckFieldSet(index+".file_path", attrs[index+".file_path"]); err != nil {
				return err
			}
			if err := util.CheckFieldSet(index+".file_chmod", attrs[index+".file_chmod"]); err != nil {
				return err
			}
			if err := util.CheckFieldSet(index+".checksum", attrs[index+".checksum"]); err != nil {
				return err
			}
			if err := util.CheckFieldSet(index+".key_fingerprint", attrs[index+".key_fingerprint"]); err != nil {
				return err
			}
			if err := util.CheckFieldSet(index+".public_value", attrs[index+".public_value"]); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccSourceVariablesSshKeysConfig(privateKey string, privateKey2 string) string {
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

`, util.UniqueString(), privateKey, privateKey2, privateKey)
}
