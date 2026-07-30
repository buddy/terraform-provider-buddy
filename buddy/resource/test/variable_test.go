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
	note := util.RandString(10)
	newNote := util.RandString(10)
	legacyDescription := util.RandString(10)
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
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, "", key, val, "", false, false),
				),
			},
			// update variable value
			{
				Config: testAccVariableWorkspaceSimpleConfig(domain, key, newValue),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, "", key, newValue, "", false, false),
				),
			},
			// update variable key
			{
				Config: testAccVariableWorkspaceComplexConfig(domain, newKey, newValue, false, true, "note", note),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, "", newKey, newValue, note, false, true),
				),
			},
			// update options
			{
				Config: testAccVariableWorkspaceComplexConfig(domain, newKey, newValue, true, true, "note", newNote),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, "", newKey, newValue, newNote, true, true),
				),
			},
			// deprecated description feeds note
			{
				Config: testAccVariableWorkspaceComplexConfig(domain, newKey, newValue, true, true, "description", legacyDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, "", newKey, newValue, legacyDescription, true, true),
				),
			},
			// migrate from deprecated description to note
			{
				Config: testAccVariableWorkspaceComplexConfig(domain, newKey, newValue, true, true, "note", newNote),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, "", newKey, newValue, newNote, true, true),
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
	note := util.RandString(10)
	newNote := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccVariableCheckDestroy,
		Steps: []resource.TestStep{
			// create variable
			{
				Config: testAccVariableProjectComplexConfig(domain, projectName, key, val, true, true, note),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, projectName, key, val, note, true, true),
				),
			},
			// update variable
			{
				Config: testAccVariableProjectComplexConfig(domain, projectName, key, newValue, false, false, newNote),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, projectName, key, newValue, newNote, false, false),
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

func TestAccVariable_environment(t *testing.T) {
	var variable buddy.Variable
	domain := util.UniqueString()
	key := util.UniqueString()
	val := util.RandString(10)
	newValue := util.RandString(10)
	note := util.RandString(10)
	newNote := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccVariableCheckDestroy,
		Steps: []resource.TestStep{
			// create variable
			{
				Config: testAccVariableEnvironmentComplexConfig(domain, key, val, true, true, note),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, "", key, val, note, true, true),
				),
			},
			// update variable
			{
				Config: testAccVariableEnvironmentComplexConfig(domain, key, newValue, false, false, newNote),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, "", key, newValue, newNote, false, false),
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

func testAccVariableAttributes(n string, variable *buddy.Variable, domain string, projectName string, key string, val string, note string, encrypted bool, settable bool) resource.TestCheckFunc {
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
			if !strings.HasPrefix(variable.Value, "!encrypted") {
				return util.ErrorFieldFormatted("Value", variable.Value, "!encrypted")
			}
		}
		if projectName != "" {
			if err := util.CheckFieldEqualAndSet("Project.Name", variable.Project.Name, projectName); err != nil {
				return err
			}
		}
		if err := util.CheckBoolFieldEqual("Encrypted", variable.Encrypted, encrypted); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("Settable", variable.Settable, settable); err != nil {
			return err
		}
		if err := util.CheckFieldEqual("Note", variable.Note, note); err != nil {
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
			if !strings.HasPrefix(attrs["value_processed"], "!encrypted") {
				return util.ErrorFieldFormatted("value_processed", attrs["value_processed"], "!encrypted")
			}
		}
		if projectName != "" {
			if err := util.CheckFieldEqualAndSet("project_name", attrs["project_name"], projectName); err != nil {
				return err
			}
		}
		if err := util.CheckBoolFieldEqual("encrypted", attrsEncrypted, encrypted); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("settable", attrsSettable, settable); err != nil {
			return err
		}
		if err := util.CheckFieldEqual("note", attrs["note"], note); err != nil {
			return err
		}
		// deprecated description mirrors note
		if err := util.CheckFieldEqual("description", attrs["description"], note); err != nil {
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

func testAccVariableProjectComplexConfig(domain string, projectName string, key string, val string, encrypted bool, settable bool, note string) string {
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
	note = "%s"
}
`, domain, projectName, key, val, encrypted, settable, note)
}

func testAccVariableEnvironmentComplexConfig(domain string, key string, val string, encrypted bool, settable bool, note string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_environment" "aha" {
	domain = "${buddy_workspace.foo.domain}"
	name = "abc"
	identifier = "abc"
}

resource "buddy_variable" "bar" {
	domain = "${buddy_workspace.foo.domain}"
	environment_id = "${buddy_environment.aha.environment_id}"
	key = "%s"
  value = "%s"
	encrypted = %t
	settable = %t
	note = "%s"
}
`, domain, key, val, encrypted, settable, note)
}

func testAccVariableWorkspaceComplexConfig(domain string, key string, val string, encrypted bool, settable bool, noteField string, note string) string {
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
	%s = "%s"
}
`, domain, key, val, encrypted, settable, noteField, note)
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
