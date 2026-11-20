package test

import (
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"
)

func TestAccSourceRuleset(t *testing.T) {
	domain := util.UniqueString()
	projectName := util.UniqueString()
	identifier := util.UniqueString()
	name := util.RandString(10)
	note := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceRulesetConfig(domain, projectName, identifier, name, note),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceRulesetChecks("data.buddy_ruleset.id", identifier, name, note),
					testAccSourceRulesetChecks("data.buddy_ruleset.identifier", identifier, name, note),
					resource.TestCheckResourceAttrPair("data.buddy_ruleset.identifier", "ruleset_id", "buddy_ruleset.bar", "ruleset_id"),
				),
			},
		},
	})
}

func TestAccSourceRulesets(t *testing.T) {
	domain := util.UniqueString()
	projectName := util.UniqueString()
	identifier1 := "aaaa" + util.UniqueString()
	identifier2 := util.UniqueString()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceRulesetsConfig(domain, projectName, identifier1, identifier2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.buddy_rulesets.all", "rulesets.#", "2"),
					resource.TestCheckResourceAttr("data.buddy_rulesets.filtered", "rulesets.#", "1"),
					resource.TestCheckResourceAttr("data.buddy_rulesets.filtered", "rulesets.0.identifier", identifier1),
					resource.TestCheckNoResourceAttr("data.buddy_rulesets.filtered", "rulesets.0.name"),
					resource.TestCheckResourceAttr("data.buddy_rulesets.filtered", "rulesets.0.disabled", "true"),
					resource.TestCheckResourceAttr("data.buddy_rulesets.filtered", "rulesets.0.refs.#", "1"),
					resource.TestCheckResourceAttrSet("data.buddy_rulesets.filtered", "rulesets.0.ruleset_id"),
					resource.TestCheckResourceAttrSet("data.buddy_rulesets.filtered", "rulesets.0.html_url"),
				),
			},
		},
	})
}

func testAccSourceRulesetChecks(n string, identifier string, name string, note string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(n, "ruleset_id"),
		resource.TestCheckResourceAttrSet(n, "html_url"),
		resource.TestCheckResourceAttr(n, "identifier", identifier),
		resource.TestCheckResourceAttr(n, "name", name),
		resource.TestCheckResourceAttr(n, "note", note),
		resource.TestCheckNoResourceAttr(n, "agent_note"),
		resource.TestCheckResourceAttr(n, "disabled", "false"),
		resource.TestCheckResourceAttr(n, "refs.#", "1"),
		resource.TestCheckTypeSetElemAttr(n, "refs.*", buddy.RulesetRefDefault),
		resource.TestCheckResourceAttr(n, "bypass_admins", "true"),
		resource.TestCheckNoResourceAttr(n, "bypass_others"),
		resource.TestCheckResourceAttr(n, "bypass_user.#", "0"),
		resource.TestCheckResourceAttr(n, "bypass_group.#", "0"),
		resource.TestCheckResourceAttr(n, "required_pipelines.#", "1"),
		resource.TestCheckNoResourceAttr(n, "required_unit_tests_suites"),
		resource.TestCheckNoResourceAttr(n, "required_visual_tests_suites"),
		resource.TestCheckResourceAttr(n, "require_up_to_date", "false"),
		resource.TestCheckResourceAttr(n, "block_force_push", "true"),
		resource.TestCheckResourceAttr(n, "restrict_deletion", "false"),
		resource.TestCheckResourceAttr(n, "missing_sources.#", "0"),
	)
}

func testAccSourceRulesetConfig(domain string, projectName string, identifier string, name string, note string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
  domain = "%s"
}

resource "buddy_project" "proj" {
  domain       = "${buddy_workspace.foo.domain}"
  display_name = "%s"
}

resource "buddy_pipeline" "pip" {
  domain       = "${buddy_workspace.foo.domain}"
  project_name = "${buddy_project.proj.name}"
  name         = "build"
}

resource "buddy_ruleset" "bar" {
  domain             = "${buddy_workspace.foo.domain}"
  project_name       = "${buddy_project.proj.name}"
  identifier         = "%s"
  name               = "%s"
  note               = "%s"
  refs               = [":default"]
  bypass_admins      = true
  required_pipelines = ["${buddy_pipeline.pip.identifier}"]
  block_force_push   = true
}

data "buddy_ruleset" "id" {
  domain       = "${buddy_workspace.foo.domain}"
  project_name = "${buddy_project.proj.name}"
  ruleset_id   = "${buddy_ruleset.bar.ruleset_id}"
}

data "buddy_ruleset" "identifier" {
  domain       = "${buddy_workspace.foo.domain}"
  project_name = "${buddy_project.proj.name}"
  identifier   = upper("${buddy_ruleset.bar.identifier}")
}
`, domain, projectName, identifier, name, note)
}

func testAccSourceRulesetsConfig(domain string, projectName string, identifier1 string, identifier2 string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
  domain = "%s"
}

resource "buddy_project" "proj" {
  domain       = "${buddy_workspace.foo.domain}"
  display_name = "%s"
}

resource "buddy_ruleset" "a" {
  domain       = "${buddy_workspace.foo.domain}"
  project_name = "${buddy_project.proj.name}"
  identifier   = "%s"
  disabled     = true
  refs         = [":default"]
}

resource "buddy_ruleset" "b" {
  domain       = "${buddy_workspace.foo.domain}"
  project_name = "${buddy_project.proj.name}"
  identifier   = "%s"
  refs         = ["refs/tags/*"]
}

data "buddy_rulesets" "all" {
  domain       = "${buddy_workspace.foo.domain}"
  project_name = "${buddy_project.proj.name}"
  depends_on   = [buddy_ruleset.a, buddy_ruleset.b]
}

data "buddy_rulesets" "filtered" {
  domain           = "${buddy_workspace.foo.domain}"
  project_name     = "${buddy_project.proj.name}"
  identifier_regex = "^aaaa"
  depends_on       = [buddy_ruleset.a, buddy_ruleset.b]
}
`, domain, projectName, identifier1, identifier2)
}
