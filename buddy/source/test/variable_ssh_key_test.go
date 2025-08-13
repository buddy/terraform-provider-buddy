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

func TestAccSourceVariableSshKey(t *testing.T) {
	domain := util.UniqueString()
	key := util.RandString(10)
	desc := util.RandString(10)
	filePlace := buddy.VariableSshKeyFilePlaceContainer
	filePath := "~/.ssh/test2"
	fileChmod := "660"
	_, privateKey, err := util.GenerateRsaKeyPair()
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
				Config: testAccSourceVariableSshKeyConfig(domain, key, desc, privateKey, filePlace, filePath, fileChmod),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceVariableSshKeyAttributes("data.buddy_variable_ssh_key.id", key, desc, filePlace, filePath, fileChmod),
					testAccSourceVariableSshKeyAttributes("data.buddy_variable_ssh_key.key", key, desc, filePlace, filePath, fileChmod),
				),
			},
		},
	})
}

func testAccSourceVariableSshKeyAttributes(n string, key string, desc string, filePlace string, filePath string, fileChmod string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsVariableId, _ := strconv.Atoi(attrs["variable_id"])
		attrsEncrypted, _ := strconv.ParseBool(attrs["encrypted"])
		attrsSettable, _ := strconv.ParseBool(attrs["settable"])
		if err := util.CheckFieldEqualAndSet("key", attrs["key"], key); err != nil {
			return err
		}
		if err := util.CheckFieldSet("value", attrs["value"]); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("file_place", attrs["file_place"], filePlace); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("file_path", attrs["file_path"], filePath); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("file_chmod", attrs["file_chmod"], fileChmod); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("settable", attrsSettable, false); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("description", attrs["description"], desc); err != nil {
			return err
		}
		if err := util.CheckIntFieldSet("variable_id", attrsVariableId); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("encrypted", attrsEncrypted, true); err != nil {
			return err
		}
		if err := util.CheckFieldSet("checksum", attrs["checksum"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("key_fingerprint", attrs["key_fingerprint"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("public_value", attrs["public_value"]); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourceVariableSshKeyConfig(domain string, key string, desc string, val string, filePlace string, filePath string, fileChmod string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_variable_ssh_key" "var" {
   domain = "${buddy_workspace.foo.domain}"
   key = "%s"
   file_place = "%s"
   file_path = "%s"
   file_chmod = "%s"
	description = "%s"
   value = <<EOT
%s
EOT
}

data "buddy_variable_ssh_key" "id" {
   domain = "${buddy_workspace.foo.domain}"
   variable_id = "${buddy_variable_ssh_key.var.variable_id}"
}

data "buddy_variable_ssh_key" "key" {
   domain = "${buddy_workspace.foo.domain}"
   key = "${buddy_variable_ssh_key.var.key}"
}
`, domain, key, filePlace, filePath, fileChmod, desc, val)
}
