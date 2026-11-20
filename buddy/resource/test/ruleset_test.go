package test

import (
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"regexp"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"
)

func TestAccRuleset(t *testing.T) {
	var ruleset buddy.Ruleset
	domain := util.UniqueString()
	projectName := util.UniqueString()
	email := util.RandEmail()
	groupName := util.RandString(10)
	pipelineIdentifier := util.UniqueString()
	identifier := util.UniqueString()
	newIdentifier := util.UniqueString()
	name := util.RandString(10)
	note := util.RandString(10)
	agentNote := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccRulesetCheckDestroy,
		Steps: []resource.TestStep{
			// create with every field
			{
				Config: testAccRulesetFullConfig(domain, projectName, email, groupName, pipelineIdentifier, identifier, name, note, agentNote),
				Check: resource.ComposeTestCheckFunc(
					testAccRulesetGet("buddy_ruleset.bar", &ruleset),
					resource.TestCheckResourceAttrPtr("buddy_ruleset.bar", "ruleset_id", &ruleset.Id),
					resource.TestCheckResourceAttrPtr("buddy_ruleset.bar", "html_url", &ruleset.HtmlUrl),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "identifier", identifier),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "name", name),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "note", note),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "agent_note", agentNote),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "disabled", "false"),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "refs.#", "2"),
					resource.TestCheckTypeSetElemAttr("buddy_ruleset.bar", "refs.*", buddy.RulesetRefDefault),
					resource.TestCheckTypeSetElemAttr("buddy_ruleset.bar", "refs.*", "refs/heads/release/*"),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "bypass_admins", "true"),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "bypass_others", buddy.RulesetBypassAccessDenied),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "bypass_user.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("buddy_ruleset.bar", "bypass_user.*", map[string]string{"access": buddy.RulesetBypassAccessAllowed}),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "bypass_group.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("buddy_ruleset.bar", "bypass_group.*", map[string]string{"access": buddy.RulesetBypassAccessDenied}),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "required_pipelines.#", "1"),
					resource.TestCheckTypeSetElemAttr("buddy_ruleset.bar", "required_pipelines.*", pipelineIdentifier),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "require_up_to_date", "true"),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "require_pull_request", "true"),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "block_force_push", "true"),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "restrict_deletion", "true"),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "restrict_creation", "false"),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "restrict_update", "false"),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "missing_sources.#", "0"),
					testAccRulesetApi(&ruleset, identifier, name, note, false, true, buddy.RulesetBypassAccessDenied, 1, 1, []string{pipelineIdentifier}, true),
				),
			},
			// change everything, fields left out are cleared
			{
				Config: testAccRulesetMinimalConfig(domain, projectName, email, groupName, pipelineIdentifier, newIdentifier),
				Check: resource.ComposeTestCheckFunc(
					testAccRulesetGet("buddy_ruleset.bar", &ruleset),
					resource.TestCheckResourceAttrPtr("buddy_ruleset.bar", "ruleset_id", &ruleset.Id),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "identifier", newIdentifier),
					resource.TestCheckNoResourceAttr("buddy_ruleset.bar", "name"),
					resource.TestCheckNoResourceAttr("buddy_ruleset.bar", "note"),
					resource.TestCheckNoResourceAttr("buddy_ruleset.bar", "agent_note"),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "disabled", "true"),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "refs.#", "1"),
					resource.TestCheckTypeSetElemAttr("buddy_ruleset.bar", "refs.*", "refs/tags/v*"),
					resource.TestCheckNoResourceAttr("buddy_ruleset.bar", "bypass_admins"),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "bypass_others", buddy.RulesetBypassAccessAllowed),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "bypass_user.#", "0"),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "bypass_group.#", "0"),
					resource.TestCheckNoResourceAttr("buddy_ruleset.bar", "required_pipelines"),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "require_up_to_date", "false"),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "require_pull_request", "false"),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "block_force_push", "false"),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "restrict_deletion", "false"),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "restrict_creation", "true"),
					resource.TestCheckResourceAttr("buddy_ruleset.bar", "restrict_update", "true"),
					testAccRulesetApi(&ruleset, newIdentifier, "", "", true, false, buddy.RulesetBypassAccessAllowed, 0, 0, nil, false),
				),
			},
			// import
			{
				ResourceName:      "buddy_ruleset.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// back to every field
			{
				Config: testAccRulesetFullConfig(domain, projectName, email, groupName, pipelineIdentifier, identifier, name, note, agentNote),
				Check: resource.ComposeTestCheckFunc(
					testAccRulesetGet("buddy_ruleset.bar", &ruleset),
					testAccRulesetApi(&ruleset, identifier, name, note, false, true, buddy.RulesetBypassAccessDenied, 1, 1, []string{pipelineIdentifier}, true),
				),
			},
		},
	})
}

// A ruleset without bypass entries imported into the state must plan no changes against a configuration
// without bypass_user / bypass_group blocks
func TestAccRulesetImportPlan(t *testing.T) {
	domain := util.UniqueString()
	projectName := util.UniqueString()
	identifier := util.UniqueString()
	var rulesetId string
	var project string
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccRulesetCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRulesetImportBaseConfig(domain, projectName),
				Check:  resource.TestCheckResourceAttrWith("buddy_project.proj", "name", func(v string) error { project = v; return nil }),
			},
			{
				PreConfig: func() {
					refs := []string{"refs/heads/main"}
					admins := true
					r, _, err := acc.ApiClient.RulesetService.Create(domain, project, &buddy.RulesetOps{
						Identifier: &identifier,
						Refs:       &refs,
						Bypass:     &buddy.RulesetBypass{Admins: &admins},
						Rules:      &buddy.RulesetRules{BlockForcePush: true},
					})
					if err != nil {
						t.Fatal(err)
					}
					rulesetId = r.Id
				},
				Config:             testAccRulesetImportConfig(domain, projectName, identifier),
				ResourceName:       "buddy_ruleset.bar",
				ImportState:        true,
				ImportStatePersist: true,
				ImportStateIdFunc: func(*terraform.State) (string, error) {
					return util.ComposeTripleId(domain, project, rulesetId), nil
				},
			},
			{
				Config:   testAccRulesetImportConfig(domain, projectName, identifier),
				PlanOnly: true,
			},
		},
	})
}

func testAccRulesetImportBaseConfig(domain string, projectName string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
  domain = "%s"
}

resource "buddy_project" "proj" {
  domain       = "${buddy_workspace.foo.domain}"
  display_name = "%s"
}
`, domain, projectName)
}

func testAccRulesetImportConfig(domain string, projectName string, identifier string) string {
	return testAccRulesetImportBaseConfig(domain, projectName) + fmt.Sprintf(`
resource "buddy_ruleset" "bar" {
  domain           = "${buddy_workspace.foo.domain}"
  project_name     = "${buddy_project.proj.name}"
  identifier       = "%s"
  refs             = ["refs/heads/main"]
  bypass_admins    = true
  block_force_push = true
}
`, identifier)
}

func TestAccRulesetPipelineNotMarkingCommits(t *testing.T) {
	domain := util.UniqueString()
	projectName := util.UniqueString()
	pipelineIdentifier := util.UniqueString()
	identifier := util.UniqueString()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccRulesetCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccRulesetPipelineNotMarkingConfig(domain, projectName, pipelineIdentifier, identifier),
				ExpectError: regexp.MustCompile(`does not mark commits`),
			},
		},
	})
}

func testAccRulesetApi(ruleset *buddy.Ruleset, identifier string, name string, note string, disabled bool, admins bool, others string, users int, groups int, pipelines []string, upToDate bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if err := util.CheckFieldEqualAndSet("Identifier", ruleset.Identifier, identifier); err != nil {
			return err
		}
		if err := util.CheckFieldEqual("Name", ruleset.Name, name); err != nil {
			return err
		}
		if err := util.CheckFieldEqual("Note", ruleset.Note, note); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("Disabled", ruleset.Disabled, disabled); err != nil {
			return err
		}
		gotAdmins := ruleset.Bypass.Admins != nil && *ruleset.Bypass.Admins
		if err := util.CheckBoolFieldEqual("Bypass.Admins", gotAdmins, admins); err != nil {
			return err
		}
		if err := util.CheckFieldEqual("Bypass.Others", ruleset.Bypass.Others, others); err != nil {
			return err
		}
		if err := util.CheckIntFieldEqual("len(Bypass.Users)", len(ruleset.Bypass.Users), users); err != nil {
			return err
		}
		if err := util.CheckIntFieldEqual("len(Bypass.Groups)", len(ruleset.Bypass.Groups), groups); err != nil {
			return err
		}
		if err := util.CheckIntFieldEqual("len(RequiredChecks.Pipelines)", len(ruleset.Rules.RequiredChecks.Pipelines), len(pipelines)); err != nil {
			return err
		}
		for i, p := range pipelines {
			if err := util.CheckFieldEqual("RequiredChecks.Pipelines", ruleset.Rules.RequiredChecks.Pipelines[i], p); err != nil {
				return err
			}
		}
		if err := util.CheckBoolFieldEqual("RequiredChecks.RequireUpToDate", ruleset.Rules.RequiredChecks.RequireUpToDate, upToDate); err != nil {
			return err
		}
		return nil
	}
}

func testAccRulesetGet(n string, ruleset *buddy.Ruleset) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		domain, projectName, rulesetId, err := util.DecomposeTripleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		r, _, err := acc.ApiClient.RulesetService.Get(domain, projectName, rulesetId)
		if err != nil {
			return err
		}
		*ruleset = *r
		return nil
	}
}

func testAccRulesetCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_ruleset" {
			continue
		}
		domain, projectName, rulesetId, err := util.DecomposeTripleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		ruleset, resp, err := acc.ApiClient.RulesetService.Get(domain, projectName, rulesetId)
		if err == nil && ruleset != nil {
			return util.ErrorResourceExists()
		}
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}

func testAccRulesetBaseConfig(domain string, projectName string, email string, groupName string, pipelineIdentifier string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
  domain = "%s"
}

resource "buddy_project" "proj" {
  domain       = "${buddy_workspace.foo.domain}"
  display_name = "%s"
}

resource "buddy_member" "a" {
  domain = "${buddy_workspace.foo.domain}"
  email  = "%s"
}

resource "buddy_group" "g" {
  domain = "${buddy_workspace.foo.domain}"
  name   = "%s"
}

resource "buddy_permission" "a" {
  domain                  = "${buddy_workspace.foo.domain}"
  name                    = "perm"
  pipeline_access_level   = "READ_WRITE"
  repository_access_level = "READ_WRITE"
  sandbox_access_level    = "READ_ONLY"
}

resource "buddy_project_member" "m" {
  domain        = "${buddy_workspace.foo.domain}"
  project_name  = "${buddy_project.proj.name}"
  member_id     = "${buddy_member.a.member_id}"
  permission_id = "${buddy_permission.a.permission_id}"
}

resource "buddy_project_group" "g" {
  domain        = "${buddy_workspace.foo.domain}"
  project_name  = "${buddy_project.proj.name}"
  group_id      = "${buddy_group.g.group_id}"
  permission_id = "${buddy_permission.a.permission_id}"
}

resource "buddy_pipeline" "pip" {
  domain       = "${buddy_workspace.foo.domain}"
  project_name = "${buddy_project.proj.name}"
  name         = "build"
  identifier   = "%s"
}
`, domain, projectName, email, groupName, pipelineIdentifier)
}

func testAccRulesetFullConfig(domain string, projectName string, email string, groupName string, pipelineIdentifier string, identifier string, name string, note string, agentNote string) string {
	return testAccRulesetBaseConfig(domain, projectName, email, groupName, pipelineIdentifier) + fmt.Sprintf(`
resource "buddy_ruleset" "bar" {
  domain               = "${buddy_workspace.foo.domain}"
  project_name         = "${buddy_project.proj.name}"
  identifier           = "%s"
  name                 = "%s"
  note                 = "%s"
  agent_note           = "%s"
  refs                 = [":default", "refs/heads/release/*"]
  bypass_admins        = true
  bypass_others        = "DENIED"
  required_pipelines   = ["${buddy_pipeline.pip.identifier}"]
  require_up_to_date   = true
  require_pull_request = true
  block_force_push     = true
  restrict_deletion    = true

  bypass_user {
    id     = "${buddy_project_member.m.member_id}"
    access = "ALLOWED"
  }

  bypass_group {
    id     = "${buddy_project_group.g.group_id}"
    access = "DENIED"
  }
}
`, identifier, name, note, agentNote)
}

func testAccRulesetMinimalConfig(domain string, projectName string, email string, groupName string, pipelineIdentifier string, identifier string) string {
	return testAccRulesetBaseConfig(domain, projectName, email, groupName, pipelineIdentifier) + fmt.Sprintf(`
resource "buddy_ruleset" "bar" {
  domain            = "${buddy_workspace.foo.domain}"
  project_name      = "${buddy_project.proj.name}"
  identifier        = "%s"
  disabled          = true
  refs              = ["refs/tags/v*"]
  bypass_others     = "ALLOWED"
  restrict_creation = true
  restrict_update   = true
}
`, identifier)
}

func testAccRulesetPipelineNotMarkingConfig(domain string, projectName string, pipelineIdentifier string, identifier string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
  domain = "%s"
}

resource "buddy_project" "proj" {
  domain       = "${buddy_workspace.foo.domain}"
  display_name = "%s"
}

resource "buddy_pipeline" "pip" {
  domain                      = "${buddy_workspace.foo.domain}"
  project_name                = "${buddy_project.proj.name}"
  name                        = "nostatus"
  identifier                  = "%s"
  do_not_create_commit_status = true
}

resource "buddy_ruleset" "bar" {
  domain             = "${buddy_workspace.foo.domain}"
  project_name       = "${buddy_project.proj.name}"
  identifier         = "%s"
  refs               = [":default"]
  required_pipelines = ["${buddy_pipeline.pip.identifier}"]
}
`, domain, projectName, pipelineIdentifier, identifier)
}
