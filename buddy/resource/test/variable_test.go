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
	agentNote := util.RandString(10)
	newAgentNote := util.RandString(10)
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
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, "", key, val, "", "", false, false),
				),
			},
			// update variable value
			{
				Config: testAccVariableWorkspaceSimpleConfig(domain, key, newValue),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, "", key, newValue, "", "", false, false),
				),
			},
			// update variable key
			{
				Config: testAccVariableWorkspaceComplexConfig(domain, newKey, newValue, false, true, "note", note, agentNote),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, "", newKey, newValue, note, agentNote, false, true),
				),
			},
			// update options
			{
				Config: testAccVariableWorkspaceComplexConfig(domain, newKey, newValue, true, true, "note", newNote, newAgentNote),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, "", newKey, newValue, newNote, newAgentNote, true, true),
				),
			},
			// deprecated description feeds note
			{
				Config: testAccVariableWorkspaceComplexConfig(domain, newKey, newValue, true, true, "description", legacyDescription, newAgentNote),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, "", newKey, newValue, legacyDescription, newAgentNote, true, true),
				),
			},
			// migrate from deprecated description to note
			{
				Config: testAccVariableWorkspaceComplexConfig(domain, newKey, newValue, true, true, "note", newNote, newAgentNote),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, "", newKey, newValue, newNote, newAgentNote, true, true),
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
	agentNote := util.RandString(10)
	newAgentNote := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccVariableCheckDestroy,
		Steps: []resource.TestStep{
			// create variable
			{
				Config: testAccVariableProjectComplexConfig(domain, projectName, key, val, true, true, note, agentNote),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, projectName, key, val, note, agentNote, true, true),
				),
			},
			// update variable
			{
				Config: testAccVariableProjectComplexConfig(domain, projectName, key, newValue, false, false, newNote, newAgentNote),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, projectName, key, newValue, newNote, newAgentNote, false, false),
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
	agentNote := util.RandString(10)
	newAgentNote := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccVariableCheckDestroy,
		Steps: []resource.TestStep{
			// create variable
			{
				Config: testAccVariableEnvironmentComplexConfig(domain, key, val, true, true, note, agentNote),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, "", key, val, note, agentNote, true, true),
				),
			},
			// update variable
			{
				Config: testAccVariableEnvironmentComplexConfig(domain, key, newValue, false, false, newNote, newAgentNote),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAttributes("buddy_variable.bar", &variable, domain, "", key, newValue, newNote, newAgentNote, false, false),
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

func testAccVariableAttributes(n string, variable *buddy.Variable, domain string, projectName string, key string, val string, note string, agentNote string, encrypted bool, settable bool) resource.TestCheckFunc {
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
		if err := util.CheckFieldEqual("AgentNote", variable.AgentNote, agentNote); err != nil {
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
		if err := util.CheckFieldEqual("agent_note", attrs["agent_note"], agentNote); err != nil {
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

func testAccVariableProjectComplexConfig(domain string, projectName string, key string, val string, encrypted bool, settable bool, note string, agentNote string) string {
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
	agent_note = "%s"
}
`, domain, projectName, key, val, encrypted, settable, note, agentNote)
}

func testAccVariableEnvironmentComplexConfig(domain string, key string, val string, encrypted bool, settable bool, note string, agentNote string) string {
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
	agent_note = "%s"
}
`, domain, key, val, encrypted, settable, note, agentNote)
}

func testAccVariableWorkspaceComplexConfig(domain string, key string, val string, encrypted bool, settable bool, noteField string, note string, agentNote string) string {
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
	agent_note = "%s"
}
`, domain, key, val, encrypted, settable, noteField, note, agentNote)
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

func testAccVariableAccessRulesAttributes(name string, variable *buddy.Variable, runOnlySettable bool, disabled bool, pipelinesAccessLevel string, sandboxesAccessLevel string, allowedPipelines int, allowedAction string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		attrs := rs.Primary.Attributes
		attrsRunOnlySettable, _ := strconv.ParseBool(attrs["run_only_settable"])
		attrsDisabled, _ := strconv.ParseBool(attrs["disabled"])
		if err := util.CheckBoolFieldEqual("RunOnlySettable", variable.RunOnlySettable, runOnlySettable); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("run_only_settable", attrsRunOnlySettable, runOnlySettable); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("Disabled", variable.Disabled, disabled); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("disabled", attrsDisabled, disabled); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("PipelinesAccessLevel", variable.PipelinesAccessLevel, pipelinesAccessLevel); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("pipelines_access_level", attrs["pipelines_access_level"], pipelinesAccessLevel); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("SandboxesAccessLevel", variable.SandboxesAccessLevel, sandboxesAccessLevel); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("sandboxes_access_level", attrs["sandboxes_access_level"], sandboxesAccessLevel); err != nil {
			return err
		}
		if err := util.CheckIntFieldEqual("len(AllowedPipelines)", len(variable.AllowedPipelines), allowedPipelines); err != nil {
			return err
		}
		if allowedPipelines > 0 {
			if err := util.CheckFieldEqualAndSet("AllowedPipelines[0].AccessLevel", variable.AllowedPipelines[0].AccessLevel, buddy.VariableAccessLevelDenied); err != nil {
				return err
			}
			if err := util.CheckFieldEqual("AllowedPipelines[0].Action", variable.AllowedPipelines[0].Action, allowedAction); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccVariableAccessRulesConfig(domain string, projectName string, pipelineName string, key string, val string, accessLevel string, action string) string {
	actionAttr := ""
	if action != "" {
		actionAttr = fmt.Sprintf("action = \"%s\"", action)
	}
	allowed := ""
	if accessLevel != "" {
		allowed = fmt.Sprintf(`
   allowed_pipeline {
      project = "${buddy_project.proj.name}"
      pipeline = "${buddy_pipeline.pipe.name}"
      access_level = "%s"
      %s
   }
`, accessLevel, actionAttr)
	}
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_project" "proj" {
   domain = "${buddy_workspace.foo.domain}"
   display_name = "%s"
}

resource "buddy_pipeline" "pipe" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   name = "%s"
   event {
      type = "PUSH"
      refs = ["refs/heads/master"]
   }
}

resource "buddy_variable" "bar" {
   domain = "${buddy_workspace.foo.domain}"
   key = "%s"
   value = "%s"
   settable = true
   run_only_settable = true
   disabled = true
   pipelines_access_level = "USE_ONLY"
   sandboxes_access_level = "DENIED"
%s
}
`, domain, projectName, pipelineName, key, val, allowed)
}

func TestAccVariable_accessRules(t *testing.T) {
	var variable buddy.Variable
	domain := util.UniqueString()
	projectName := util.UniqueString()
	pipelineName := util.UniqueString()
	key := util.UniqueString()
	val := util.RandString(10)
	action := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccVariableCheckDestroy,
		Steps: []resource.TestStep{
			// whole pipeline rule
			{
				Config: testAccVariableAccessRulesConfig(domain, projectName, pipelineName, key, val, buddy.VariableAccessLevelDenied, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAccessRulesAttributes("buddy_variable.bar", &variable, true, true, buddy.VariableAccessLevelUseOnly, buddy.VariableAccessLevelDenied, 1, ""),
				),
			},
			// swap it for a single action rule
			{
				Config: testAccVariableAccessRulesConfig(domain, projectName, pipelineName, key, val, buddy.VariableAccessLevelDenied, action),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAccessRulesAttributes("buddy_variable.bar", &variable, true, true, buddy.VariableAccessLevelUseOnly, buddy.VariableAccessLevelDenied, 1, action),
				),
			},
			// drop every rule
			{
				Config: testAccVariableAccessRulesConfig(domain, projectName, pipelineName, key, val, "", ""),
				Check: resource.ComposeTestCheckFunc(
					testAccVariableGet("buddy_variable.bar", &variable),
					testAccVariableAccessRulesAttributes("buddy_variable.bar", &variable, true, true, buddy.VariableAccessLevelUseOnly, buddy.VariableAccessLevelDenied, 0, ""),
				),
			},
		},
	})
}
