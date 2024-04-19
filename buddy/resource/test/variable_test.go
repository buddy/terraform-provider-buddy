package test

import (
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"strings"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"
)

func TestAccVariable_workspace(t *testing.T) {
	var variable buddy.Variable
	domain := util.UniqueString()
	key := util.UniqueString()
	val := util.RandString(10)
	newValue := util.RandString(10)
	newKey := util.RandString(10)
	desc := util.RandString(10)
	newDesc := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccVariableCheckDestroy,
		Steps: []resource.TestStep{
			// create variable
			{
				Config: testAccVariableWorkspaceSimpleConfig(domain, key, val),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, key, val, "", false, false),
				),
			},
			// update variable value
			{
				Config: testAccVariableWorkspaceSimpleConfig(domain, key, newValue),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, key, newValue, "", false, false),
				),
			},
			// update variable key
			{
				Config: testAccVariableWorkspaceComplexConfig(domain, newKey, newValue, false, true, desc),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, newKey, newValue, desc, false, true),
				),
			},
			// update options
			{
				Config: testAccVariableWorkspaceComplexConfig(domain, newKey, newValue, true, true, newDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, newKey, newValue, newDesc, true, true),
				),
			},
			// import
			{
				ResourceName:            "buddy_variable.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value"},
			},
		},
	})
}

func TestAccVariable_project(t *testing.T) {
	var variable buddy.Variable
	domain := util.UniqueString()
	projectName := util.UniqueString()
	key := util.UniqueString()
	val := util.RandString(10)
	newValue := util.RandString(10)
	desc := util.RandString(10)
	newDesc := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccVariableCheckDestroy,
		Steps: []resource.TestStep{
			// create variable
			{
				Config: testAccVariableProjectComplexConfig(domain, projectName, key, val, true, true, desc),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, key, val, desc, true, true),
				),
			},
			// update variable
			{
				Config: testAccVariableProjectComplexConfig(domain, projectName, key, newValue, false, false, newDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, key, newValue, newDesc, false, false),
				),
			},
			// import
			{
				ResourceName:            "buddy_variable.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"value", "project_name"},
			},
		},
	})
}

func testAccVariableAttributes(n string, variable *buddy.Variable, domain string, key string, val string, description string, encrypted bool, settable bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsEncrypted, _ := strconv.ParseBool(attrs["encrypted"])
		attrsSettable, _ := strconv.ParseBool(attrs["settable"])
		attrsVariableId, _ := strconv.Atoi(attrs["variable_id"])
		if err := util.CheckFieldEqualAndSet("Key", variable.Key, key); err != nil {
			return err
		}
		if !encrypted {
			if err := util.CheckFieldEqualAndSet("Value", variable.Value, val); err != nil {
				return err
			}
		} else {
			if !strings.HasPrefix(variable.Value, "secure!") {
				return util.ErrorFieldFormatted("Value", variable.Value, "secure!")
			}
		}
		if err := util.CheckBoolFieldEqual("Encrypted", variable.Encrypted, encrypted); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("Settable", variable.Settable, settable); err != nil {
			return err
		}
		if err := util.CheckFieldEqual("Description", variable.Description, description); err != nil {
			return err
		}
		if err := util.CheckIntFieldSet("VariableId", variable.Id); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("domain", attrs["domain"], domain); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("key", attrs["key"], key); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("value", attrs["value"], val); err != nil {
			return err
		}
		if !encrypted {
			if err := util.CheckFieldEqualAndSet("value_processed", attrs["value_processed"], val); err != nil {
				return err
			}
		} else {
			if !strings.HasPrefix(attrs["value_processed"], "secure!") {
				return util.ErrorFieldFormatted("value_processed", attrs["value_processed"], "secure!")
			}
		}
		if err := util.CheckBoolFieldEqual("encrypted", attrsEncrypted, encrypted); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("settable", attrsSettable, settable); err != nil {
			return err
		}
		if err := util.CheckFieldEqual("description", attrs["description"], description); err != nil {
			return err
		}
		if err := util.CheckIntFieldSet("variable_id", attrsVariableId); err != nil {
			return err
		}
		return nil
	}
}

func testAccVariableGet(n string, variable *buddy.Variable) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		domain, vid, err := util.DecomposeDoubleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		variableId, err := strconv.Atoi(vid)
		if err != nil {
			return err
		}
		v, _, err := acc.ApiClient.VariableService.Get(domain, variableId)
		if err != nil {
			return err
		}
		*variable = *v
		return nil
	}
}

func testAccVariableWorkspaceSimpleConfig(domain string, key string, val string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_variable" "bar" {
   domain = "${buddy_workspace.foo.domain}"
   key = "%s"
   value = "%s"
}
`, domain, key, val)
}

func testAccVariableProjectComplexConfig(domain string, projectName string, key string, val string, encrypted bool, settable bool, description string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_project" "aha" {
	domain = "${buddy_workspace.foo.domain}"
	display_name = "%s"
}

resource "buddy_variable" "bar" {
   domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.aha.name}"
   key = "%s"
   value = "%s"
	encrypted = %t
	settable = %t
	description = "%s"
}
`, domain, projectName, key, val, encrypted, settable, description)
}

func testAccVariableWorkspaceComplexConfig(domain string, key string, val string, encrypted bool, settable bool, description string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_variable" "bar" {
   domain = "${buddy_workspace.foo.domain}"
   key = "%s"
   value = "%s"
	encrypted = %t
	settable = %t
	description = "%s"
}
`, domain, key, val, encrypted, settable, description)
}

func testAccVariableCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_variable" {
			continue
		}
		domain, vid, err := util.DecomposeDoubleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		variableId, err := strconv.Atoi(vid)
		if err != nil {
			return err
		}
		variable, resp, err := acc.ApiClient.VariableService.Get(domain, variableId)
		if err == nil && variable != nil {
			return util.ErrorResourceExists()
		}
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}
