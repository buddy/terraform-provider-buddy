package test

import (
	"buddy-terraform/buddy/acc"
	"buddy-terraform/buddy/util"
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"strconv"
	"testing"
)

func TestAccProject_github(t *testing.T) {
	ghToken := os.Getenv("BUDDY_GH_TOKEN")
	ghPoject := os.Getenv("BUDDY_GH_PROJECT")
	if ghToken == "" || ghPoject == "" {
		return
	}
	var project buddy.Project
	domain := util.UniqueString()
	displayName := util.RandString(10)
	newDisplayName := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccProjectCheckDestroy,
		Steps: []resource.TestStep{
			// create project
			{
				Config: testAccProjectGitHubConfig(domain, displayName, ghToken, ghPoject, false),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectGet("buddy_project.bar", &project),
					testAccProjectAttributes("buddy_project.bar", &project, displayName, false),
				),
			},
			// update project
			{
				Config: testAccProjectGitHubConfig(domain, newDisplayName, ghToken, ghPoject, true),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectGet("buddy_project.bar", &project),
					testAccProjectAttributes("buddy_project.bar", &project, newDisplayName, true),
				),
			},
			// import project
			{
				ResourceName:            "buddy_project.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"external_project_id", "integration_id"},
			},
		},
	})
}

func TestAccProject_buddy(t *testing.T) {
	var project buddy.Project
	domain := util.UniqueString()
	displayName := util.RandString(10)
	newDisplayName := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccProjectCheckDestroy,
		Steps: []resource.TestStep{
			// create project
			{
				Config: testAccProjectBuddyConfig(domain, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectGet("buddy_project.bar", &project),
					testAccProjectAttributes("buddy_project.bar", &project, displayName, false),
				),
			},
			// update project
			{
				Config: testAccProjectBuddyConfig(domain, newDisplayName),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectGet("buddy_project.bar", &project),
					testAccProjectAttributes("buddy_project.bar", &project, newDisplayName, false),
				),
			},
			// import project
			{
				ResourceName:      "buddy_project.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccProject_custom(t *testing.T) {
	var project buddy.Project
	repoUrl := "git@github.com:octocat/Hello-World.git"
	domain := util.UniqueString()
	displayName := util.RandString(10)
	newDisplayName := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccProjectCheckDestroy,
		Steps: []resource.TestStep{
			// create project
			{
				Config: testAccProjectCustomConfig(domain, repoUrl, displayName),
				Check: resource.ComposeTestCheckFunc(
					util.TestSleep(10000),
					testAccProjectGet("buddy_project.bar", &project),
					testAccProjectAttributes("buddy_project.bar", &project, displayName, false),
				),
			},
			// update project
			{
				Config: testAccProjectCustomConfig(domain, repoUrl, newDisplayName),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectGet("buddy_project.bar", &project),
					testAccProjectAttributes("buddy_project.bar", &project, newDisplayName, false),
				),
			},
			// import project
			{
				ResourceName:            "buddy_project.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_repo_url"},
			},
		},
	})
}

func testAccProjectAttributes(n string, project *buddy.Project, displayName string, updateDefaultBranch bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsCreatedByMemberId, _ := strconv.Atoi(attrs["created_by.0.member_id"])
		attrsUpdateDefaultBranch, _ := strconv.ParseBool(attrs["update_default_branch_from_external"])
		if err := util.CheckFieldEqualAndSet("DisplayName", project.DisplayName, displayName); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("display_name", attrs["display_name"], displayName); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("update_default_branch_from_external", attrsUpdateDefaultBranch, updateDefaultBranch); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("UpdateDefaultBranchFromExternal", project.UpdateDefaultBranchFromExternal, updateDefaultBranch); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], project.Name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("html_url", attrs["html_url"], project.HtmlUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("status", attrs["status"], project.Status); err != nil {
			return err
		}
		if err := util.CheckDateFieldEqual("create_date", attrs["create_date"], project.CreateDate); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("http_repository", attrs["http_repository"], project.HttpRepository); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("ssh_repository", attrs["ssh_repository"], project.SshRepository); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("default_branch", attrs["default_branch"], project.DefaultBranch); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("created_by.0.html_url", attrs["created_by.0.html_url"], project.CreatedBy.HtmlUrl); err != nil {
			return err
		}
		if err := util.CheckIntFieldEqualAndSet("created_by.0.member_id", attrsCreatedByMemberId, project.CreatedBy.Id); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("created_by.0.name", attrs["created_by.0.name"], project.CreatedBy.Name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("created_by.0.avatar_url", attrs["created_by.0.avatar_url"], project.CreatedBy.AvatarUrl); err != nil {
			return err
		}
		return nil
	}
}

func testAccProjectGet(n string, project *buddy.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		domain, name, err := util.DecomposeDoubleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		p, _, err := acc.ApiClient.ProjectService.Get(domain, name)
		if err != nil {
			return err
		}
		*project = *p
		return nil
	}
}

func testAccProjectCustomConfig(domain string, repoUrl string, displayName string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_project" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
	custom_repo_url = "%s"
}
`, domain, displayName, repoUrl)
}

func testAccProjectBuddyConfig(domain string, name string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_project" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s" 
}
`, domain, name)
}

func testAccProjectGitHubConfig(domain string, name string, ghToken string, ghProject string, updateBranch bool) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_integration" "gh" {
	domain = "${buddy_workspace.foo.domain}"
	name = "gh"
    type = "GIT_HUB"
    scope = "WORKSPACE"
    token = "%s"
}

resource "buddy_project" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
    integration_id = "${buddy_integration.gh.integration_id}"
    external_project_id = "%s"
    update_default_branch_from_external = "%t"
}
`, domain, ghToken, name, ghProject, updateBranch)
}

func testAccProjectCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_project" {
			continue
		}
		domain, name, err := util.DecomposeDoubleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		project, resp, err := acc.ApiClient.ProjectService.Get(domain, name)
		if err == nil && project != nil {
			return util.ErrorResourceExists()
		}
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}
