package test

import (
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"
)

// Changing anything but the scope of a variable must update it in place (same variable_id),
// changing the scope must still replace it.

// stateValueCollector is what statecheck.CompareValue returns (its type is unexported).
type stateValueCollector interface {
	AddStateValue(resourceAddress string, attributePath tfjsonpath.Path) statecheck.StateCheck
}

func TestAccVariable_scopeInPlaceUpdate(t *testing.T) {
	testAccVariableScopeInPlaceUpdate(t, "buddy_variable", fmt.Sprintf("value = \"%s\"\n   encrypted = true", util.RandString(10)), testAccVariableCheckDestroy)
}

func testAccVariableScopeInPlaceUpdate(t *testing.T, resourceType string, body string, checkDestroy resource.TestCheckFunc) {
	domain := util.UniqueString()
	projectName := util.UniqueString()
	pipelineName := util.UniqueString()
	key := util.UniqueString()
	name := resourceType + ".bar"
	sameWorkspaceVar := statecheck.CompareValue(compare.ValuesSame())
	sameProjectVar := statecheck.CompareValue(compare.ValuesSame())
	config := func(project bool, rule bool) string {
		return testAccVariableScopeConfig(resourceType, domain, projectName, pipelineName, key, body, project, rule)
	}
	step := func(project bool, rule bool, action plancheck.ResourceActionType, same stateValueCollector) resource.TestStep {
		s := resource.TestStep{
			Config: config(project, rule),
			ConfigPlanChecks: resource.ConfigPlanChecks{
				PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(name, action)},
			},
		}
		if same != nil {
			s.ConfigStateChecks = []statecheck.StateCheck{same.AddStateValue(name, tfjsonpath.New("variable_id"))}
		}
		return s
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             checkDestroy,
		Steps: []resource.TestStep{
			// workspace scope
			step(false, false, plancheck.ResourceActionCreate, sameWorkspaceVar),
			step(false, true, plancheck.ResourceActionUpdate, sameWorkspaceVar),
			// project scope replaces
			step(true, true, plancheck.ResourceActionDestroyBeforeCreate, sameProjectVar),
			step(true, false, plancheck.ResourceActionUpdate, sameProjectVar),
			// back to workspace scope (project_name removed) replaces
			step(false, false, plancheck.ResourceActionDestroyBeforeCreate, nil),
		},
	})
}

func testAccVariableScopeConfig(resourceType string, domain string, projectName string, pipelineName string, key string, body string, project bool, rule bool) string {
	scope := ""
	if project {
		scope = `project_name = "${buddy_project.proj.name}"`
	}
	allowed := ""
	if rule {
		allowed = fmt.Sprintf(`
   allowed_pipeline {
      project = "${buddy_project.proj.name}"
      pipeline = "${buddy_pipeline.pipe.name}"
      access_level = "%s"
   }
`, buddy.VariableAccessLevelDenied)
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

resource "%s" "bar" {
   domain = "${buddy_workspace.foo.domain}"
   key = "%s"
   %s
   %s
%s
}
`, domain, projectName, pipelineName, resourceType, key, scope, body, allowed)
}
