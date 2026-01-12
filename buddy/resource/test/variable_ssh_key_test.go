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

func TestAccVariableSshKey_workspace(t *testing.T) {
  var variable buddy.Variable
  domain := util.UniqueString()
  key := util.UniqueString()
  newKey := util.UniqueString()
  desc := util.RandString(10)
  filePlace := buddy.VariableSshKeyFilePlaceContainer
  filePath := "~/.ssh/test"
  fileChmod := "600"
  _, privateKey, err := util.GenerateRsaKeyPair()
  if err != nil {
    t.Fatal(err.Error())
  }
  _, privateKey2, err := util.GenerateRsaKeyPair()
  if err != nil {
    t.Fatal(err.Error())
  }
  resource.Test(t, resource.TestCase{
    PreCheck: func() {
      acc.PreCheck(t)
    },
    ProtoV6ProviderFactories: acc.ProviderFactories,
    CheckDestroy:             testAccVariableSshKeyCheckDestroy,
    Steps: []resource.TestStep{
      // create variable
      {
        Config: testAccVariableSshKeyWorkspaceSimpleConfig(domain, key, filePlace, filePath, fileChmod, privateKey),
        Check: resource.ComposeTestCheckFunc(
          testAccVariableGet("buddy_variable_ssh_key.bar", &variable),
          testAccVariableSshKeyAttributes("buddy_variable_ssh_key.bar", &variable, domain, "", key, privateKey, filePlace, filePath, fileChmod, ""),
        ),
      },
      // update variable value
      {
        Config: testAccVariableSshKeyWorkspaceSimpleConfig(domain, key, filePlace, filePath, fileChmod, privateKey2),
        Check: resource.ComposeTestCheckFunc(
          testAccVariableGet("buddy_variable_ssh_key.bar", &variable),
          testAccVariableSshKeyAttributes("buddy_variable_ssh_key.bar", &variable, domain, "", key, privateKey2, filePlace, filePath, fileChmod, ""),
        ),
      },
      // update variable key
      {
        Config: testAccVariableSshKeyWorkspaceSimpleConfig(domain, newKey, filePlace, filePath, fileChmod, privateKey2),
        Check: resource.ComposeTestCheckFunc(
          testAccVariableGet("buddy_variable_ssh_key.bar", &variable),
          testAccVariableSshKeyAttributes("buddy_variable_ssh_key.bar", &variable, domain, "", newKey, privateKey2, filePlace, filePath, fileChmod, ""),
        ),
      },
      // update options
      {
        Config: testAccVariableSshKeyWorkspaceComplexConfig(domain, newKey, filePlace, filePath, fileChmod, desc, privateKey2),
        Check: resource.ComposeTestCheckFunc(
          testAccVariableGet("buddy_variable_ssh_key.bar", &variable),
          testAccVariableSshKeyAttributes("buddy_variable_ssh_key.bar", &variable, domain, "", newKey, privateKey2, filePlace, filePath, fileChmod, desc),
        ),
      },
      // import
      {
        ResourceName:            "buddy_variable_ssh_key.bar",
        ImportState:             true,
        ImportStateVerify:       true,
        ImportStateVerifyIgnore: []string{"value"},
      },
    },
  })
}

func TestAccVariableSshKey_project(t *testing.T) {
  var variable buddy.Variable
  domain := util.UniqueString()
  projectName := util.UniqueString()
  key := util.UniqueString()
  desc := util.RandString(10)
  filePlace := buddy.VariableSshKeyFilePlaceContainer
  filePath := "~/.ssh/test2"
  fileChmod := "660"
  _, privateKey, err := util.GenerateRsaKeyPair()
  if err != nil {
    t.Fatal(err.Error())
  }
  _, privateKey2, err := util.GenerateRsaKeyPair()
  if err != nil {
    t.Fatal(err.Error())
  }
  resource.Test(t, resource.TestCase{
    PreCheck: func() {
      acc.PreCheck(t)
    },
    ProtoV6ProviderFactories: acc.ProviderFactories,
    CheckDestroy:             testAccVariableSshKeyCheckDestroy,
    Steps: []resource.TestStep{
      // create variable
      {
        Config: testAccVariableSshKeyProjectComplexConfig(domain, projectName, key, filePlace, filePath, fileChmod, "", privateKey),
        Check: resource.ComposeTestCheckFunc(
          testAccVariableGet("buddy_variable_ssh_key.bar", &variable),
          testAccVariableSshKeyAttributes("buddy_variable_ssh_key.bar", &variable, domain, projectName, key, privateKey, filePlace, filePath, fileChmod, ""),
        ),
      },
      // update variable
      {
        Config: testAccVariableSshKeyProjectComplexConfig(domain, projectName, key, filePlace, filePath, fileChmod, desc, privateKey2),
        Check: resource.ComposeTestCheckFunc(
          testAccVariableGet("buddy_variable_ssh_key.bar", &variable),
          testAccVariableSshKeyAttributes("buddy_variable_ssh_key.bar", &variable, domain, projectName, key, privateKey2, filePlace, filePath, fileChmod, desc),
        ),
      },
      // import
      {
        ResourceName:            "buddy_variable_ssh_key.bar",
        ImportState:             true,
        ImportStateVerify:       true,
        ImportStateVerifyIgnore: []string{"value"},
      },
    },
  })
}

func TestAccVariableSshKey_environment(t *testing.T) {
  var variable buddy.Variable
  domain := util.UniqueString()
  key := util.UniqueString()
  desc := util.RandString(10)
  filePlace := buddy.VariableSshKeyFilePlaceContainer
  filePath := "~/.ssh/test3"
  fileChmod := "666"
  _, privateKey, err := util.GenerateRsaKeyPair()
  if err != nil {
    t.Fatal(err.Error())
  }
  _, privateKey2, err := util.GenerateRsaKeyPair()
  if err != nil {
    t.Fatal(err.Error())
  }
  resource.Test(t, resource.TestCase{
    PreCheck: func() {
      acc.PreCheck(t)
    },
    ProtoV6ProviderFactories: acc.ProviderFactories,
    CheckDestroy:             testAccVariableSshKeyCheckDestroy,
    Steps: []resource.TestStep{
      // create variable
      {
        Config: testAccVariableSshKeyEnvironmentComplexConfig(domain, key, filePlace, filePath, fileChmod, "", privateKey),
        Check: resource.ComposeTestCheckFunc(
          testAccVariableGet("buddy_variable_ssh_key.bar", &variable),
          testAccVariableSshKeyAttributes("buddy_variable_ssh_key.bar", &variable, domain, "", key, privateKey, filePlace, filePath, fileChmod, ""),
        ),
      },
      // update variable
      {
        Config: testAccVariableSshKeyEnvironmentComplexConfig(domain, key, filePlace, filePath, fileChmod, desc, privateKey2),
        Check: resource.ComposeTestCheckFunc(
          testAccVariableGet("buddy_variable_ssh_key.bar", &variable),
          testAccVariableSshKeyAttributes("buddy_variable_ssh_key.bar", &variable, domain, "", key, privateKey2, filePlace, filePath, fileChmod, desc),
        ),
      },
      // import
      {
        ResourceName:            "buddy_variable_ssh_key.bar",
        ImportState:             true,
        ImportStateVerify:       true,
        ImportStateVerifyIgnore: []string{"value"},
      },
    },
  })
}

func testAccVariableSshKeyAttributes(n string, variable *buddy.Variable, domain string, projectName string, key string, val string, filePlace string, filePath string, fileChmod string, description string) resource.TestCheckFunc {
  return func(s *terraform.State) error {
    rs, ok := s.RootModule().Resources[n]
    if !ok {
      return fmt.Errorf("not found: %s", n)
    }
    attrs := rs.Primary.Attributes
    attrsVariableId, _ := strconv.Atoi(attrs["variable_id"])
    attrsEncrypted, _ := strconv.ParseBool(attrs["encrypted"])
    attrsSettable, _ := strconv.ParseBool(attrs["settable"])
    if projectName != "" {
      if err := util.CheckFieldEqualAndSet("Project.Name", variable.Project.Name, projectName); err != nil {
        return err
      }
      if err := util.CheckFieldEqualAndSet("project_name", attrs["project_name"], projectName); err != nil {
        return err
      }
    }
    if err := util.CheckFieldEqualAndSet("Key", variable.Key, key); err != nil {
      return err
    }
    if !strings.HasPrefix(variable.Value, "!encrypted") {
      return util.ErrorFieldFormatted("Value", variable.Value, "!encrypted")
    }
    if err := util.CheckFieldSet("Checksum", variable.Checksum); err != nil {
      return err
    }
    if err := util.CheckFieldSet("KeyFingerprint", variable.KeyFingerprint); err != nil {
      return err
    }
    if err := util.CheckFieldSet("PublicValue", variable.PublicValue); err != nil {
      return err
    }
    if err := util.CheckBoolFieldEqual("Settable", variable.Settable, false); err != nil {
      return err
    }
    if err := util.CheckBoolFieldEqual("Encrypted", variable.Encrypted, true); err != nil {
      return err
    }
    if err := util.CheckIntFieldSet("VariableId", variable.Id); err != nil {
      return err
    }
    if err := util.CheckFieldEqualAndSet("FilePlace", variable.FilePlace, filePlace); err != nil {
      return err
    }
    if err := util.CheckFieldEqualAndSet("FilePath", variable.FilePath, filePath); err != nil {
      return err
    }
    if err := util.CheckFieldEqual("Description", variable.Description, description); err != nil {
      return err
    }
    if err := util.CheckFieldEqualAndSet("FileChmod", variable.FileChmod, fileChmod); err != nil {
      return err
    }
    if err := util.CheckFieldEqualAndSet("domain", attrs["domain"], domain); err != nil {
      return err
    }
    if err := util.CheckFieldEqualAndSet("key", attrs["key"], key); err != nil {
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
    if err := util.CheckFieldEqualAndSet("value", attrs["value"], val+"\n"); err != nil {
      return err
    }
    if !strings.HasPrefix(attrs["value_processed"], "!encrypted") {
      return util.ErrorFieldFormatted("value_processed", attrs["value_processed"], "!encrypted")
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
    if err := util.CheckFieldEqual("description", attrs["description"], description); err != nil {
      return err
    }
    if err := util.CheckIntFieldSet("variable_id", attrsVariableId); err != nil {
      return err
    }
    return nil
  }
}

func testAccVariableSshKeyProjectComplexConfig(domain string, projectName string, key string, filePlace string, filePath string, fileChmod string, description string, val string) string {
  return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_project" "aha" {
	domain = "${buddy_workspace.foo.domain}"
	display_name = "%s"
}

resource "buddy_variable_ssh_key" "bar" {
   domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.aha.name}"
   key = "%s"
   file_place = "%s"
   file_path = "%s"
   file_chmod = "%s"
	description = "%s"
   value = <<EOT
%s
EOT
}
`, domain, projectName, key, filePlace, filePath, fileChmod, description, val)
}

func testAccVariableSshKeyEnvironmentComplexConfig(domain string, key string, filePlace string, filePath string, fileChmod string, description string, val string) string {
  return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_environment" "aha" {
	domain = "${buddy_workspace.foo.domain}"
	name = "efg"
	identifier = "efg" 
}

resource "buddy_variable_ssh_key" "bar" {
   domain = "${buddy_workspace.foo.domain}"
	 environment_id = "${buddy_environment.aha.environment_id}"
   key = "%s"
   file_place = "%s"
   file_path = "%s"
   file_chmod = "%s"
	 description = "%s"
   value = <<EOT
%s
EOT
}
`, domain, key, filePlace, filePath, fileChmod, description, val)
}

func testAccVariableSshKeyWorkspaceComplexConfig(domain string, key string, filePlace string, filePath string, fileChmod string, description string, val string) string {
  return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_variable_ssh_key" "bar" {
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
`, domain, key, filePlace, filePath, fileChmod, description, val)
}

func testAccVariableSshKeyWorkspaceSimpleConfig(domain string, key string, filePlace string, filePath string, fileChmod string, val string) string {
  return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_variable_ssh_key" "bar" {
   domain = "${buddy_workspace.foo.domain}"
   key = "%s"
   file_place = "%s"
   file_path = "%s"
   file_chmod = "%s"
   value = <<EOT
%s
EOT
}
`, domain, key, filePlace, filePath, fileChmod, val)
}

func testAccVariableSshKeyCheckDestroy(s *terraform.State) error {
  for _, rs := range s.RootModule().Resources {
    if rs.Type != "buddy_variable_ssh_key" {
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
