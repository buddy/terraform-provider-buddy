package test

import (
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"os"
	"strconv"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"
)

func TestAccProject_github_upgrade(t *testing.T) {
	ghToken := os.Getenv("BUDDY_GH_TOKEN")
	ghPoject := os.Getenv("BUDDY_GH_PROJECT")
	if ghToken == "" || ghPoject == "" {
		return
	}
	var project buddy.Project
	domain := util.UniqueString()
	displayName := util.RandString(10)
	config := testAccProjectGitHubConfig(domain, displayName, ghToken, ghPoject, false, true)
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"buddy": {
						VersionConstraint: "1.12.0",
						Source:            "buddy/buddy",
					},
				},
				Config: config,
			},
			{
				ProtoV6ProviderFactories: acc.ProviderFactories,
				Config:                   config,
				Check: resource.ComposeTestCheckFunc(
					util.TestSleep(10000),
					testAccProjectGet("buddy_project.gh", &project),
					testAccProjectAttributes("buddy_project.gh", &project, displayName, false, buddy.ProjectAccessPrivate, true, false, ""),
				),
			},
		},
	})
}

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
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccProjectCheckDestroy,
		Steps: []resource.TestStep{
			// create project
			{
				Config: testAccProjectGitHubConfig(domain, displayName, ghToken, ghPoject, false, true),
				Check: resource.ComposeTestCheckFunc(
					util.TestSleep(10000),
					testAccProjectGet("buddy_project.gh", &project),
					testAccProjectAttributes("buddy_project.gh", &project, displayName, false, buddy.ProjectAccessPrivate, true, false, ""),
				),
			},
			// update project
			{
				Config: testAccProjectGitHubConfig(domain, newDisplayName, ghToken, ghPoject, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectGet("buddy_project.gh", &project),
					testAccProjectAttributes("buddy_project.gh", &project, newDisplayName, true, buddy.ProjectAccessPrivate, false, false, ""),
				),
			},
			// import project
			{
				ResourceName:            "buddy_project.gh",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"external_project_id", "integration_id"},
			},
		},
	})
}

func TestAccProject_withoutRepository(t *testing.T) {
	var project buddy.Project
	domain := util.UniqueString()
	displayName := util.RandString(10)
	newDisplayName := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccProjectCheckDestroy,
		Steps: []resource.TestStep{
			// create project
			{
				Config: testAccProjectWithoutRepositoryConfig(domain, displayName),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectGet("buddy_project.bar", &project),
					testAccProjectAttributes("buddy_project.bar", &project, displayName, true, buddy.ProjectAccessPrivate, false, false, ""),
				),
			},
			// update project
			{
				Config: testAccProjectWithoutRepositoryConfig(domain, newDisplayName),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectGet("buddy_project.bar", &project),
					testAccProjectAttributes("buddy_project.bar", &project, newDisplayName, true, buddy.ProjectAccessPrivate, false, false, ""),
				),
			},
			// import project
			{
				ResourceName:            "buddy_project.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"without_repository"},
			},
		},
	})
}

func TestAccProject_buddy(t *testing.T) {
	var project buddy.Project
	domain := util.UniqueString()
	displayName := util.RandString(10)
	newDisplayName := util.RandString(10)
	access := buddy.ProjectAccessPublic
	newAccess := buddy.ProjectAccessPrivate
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccProjectCheckDestroy,
		Steps: []resource.TestStep{
			// create project
			{
				Config: testAccProjectBuddyConfig(domain, displayName, access),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectGet("buddy_project.bar", &project),
					testAccProjectAttributes("buddy_project.bar", &project, displayName, true, access, false, false, ""),
				),
			},
			// update project
			{
				Config: testAccProjectBuddyConfig(domain, newDisplayName, newAccess),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectGet("buddy_project.bar", &project),
					testAccProjectAttributes("buddy_project.bar", &project, newDisplayName, true, newAccess, false, false, ""),
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
	fetchSubmodulesEnv := "id_workspace"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccProjectCheckDestroy,
		Steps: []resource.TestStep{
			// create project
			{
				Config: testAccProjectCustomConfig(domain, repoUrl, displayName, true, fetchSubmodulesEnv),
				Check: resource.ComposeTestCheckFunc(
					util.TestSleep(10000),
					testAccProjectGet("buddy_project.bar", &project),
					testAccProjectAttributes("buddy_project.bar", &project, displayName, true, buddy.ProjectAccessPrivate, false, true, fetchSubmodulesEnv),
				),
			},
			// update project
			{
				Config: testAccProjectCustomConfig(domain, repoUrl, newDisplayName, false, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectGet("buddy_project.bar", &project),
					testAccProjectAttributes("buddy_project.bar", &project, newDisplayName, true, buddy.ProjectAccessPrivate, false, false, ""),
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

func testAccProjectAttributes(n string, project *buddy.Project, displayName string, updateDefaultBranch bool, access string, allowPullRequests bool, fetchSubmodules bool, fetchSubmodulesEnv string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsCreatedByMemberId, _ := strconv.Atoi(attrs["created_by.0.member_id"])
		attrsUpdateDefaultBranch, _ := strconv.ParseBool(attrs["update_default_branch_from_external"])
		attrsAllowPullRequests, _ := strconv.ParseBool(attrs["allow_pull_requests"])
		attrsFetchSubmodules, _ := strconv.ParseBool(attrs["fetch_submodules"])
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
		if err := util.CheckBoolFieldEqual("allow_pull_requests", attrsAllowPullRequests, allowPullRequests); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("AllowPullRequests", project.AllowPullRequests, allowPullRequests); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("fetch_submodules", attrsFetchSubmodules, fetchSubmodules); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("FetchSubmodules", project.FetchSubmodules, fetchSubmodules); err != nil {
			return err
		}
		if fetchSubmodules {
			if err := util.CheckFieldEqualAndSet("fetch_submodules_env_key", attrs["fetch_submodules_env_key"], fetchSubmodulesEnv); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("FetchSubmodulesEnvKey", project.FetchSubmodulesEnvKey, fetchSubmodulesEnv); err != nil {
				return err
			}
		}
		if err := util.CheckFieldEqualAndSet("access", attrs["access"], access); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("Access", project.Access, access); err != nil {
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
		if err := util.CheckFieldSet("default_branch", attrs["default_branch"]); err != nil {
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

func testAccProjectCustomConfig(domain string, repoUrl string, displayName string, fetchSubmodules bool, fetchSubmodulesEnv string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_project" "bar" {
   domain = "${buddy_workspace.foo.domain}"
   display_name = "%s"
	custom_repo_url = "%s"
	fetch_submodules = "%t"
   fetch_submodules_env_key = "%s"
}
`, domain, displayName, repoUrl, fetchSubmodules, fetchSubmodulesEnv)
}

func testAccProjectBuddyConfig(domain string, name string, access string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_project" "bar" {
   domain = "${buddy_workspace.foo.domain}"
   display_name = "%s"
	access = "%s"
}
`, domain, name, access)
}

func testAccProjectWithoutRepositoryConfig(domain string, name string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_project" "bar" {
   domain = "${buddy_workspace.foo.domain}"
   display_name = "%s"
	without_repository = true
}
`, domain, name)
}

func testAccProjectGitHubConfig(domain string, name string, ghToken string, ghProject string, updateBranch bool, allowPullRequests bool) string {
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

resource "buddy_project" "gh" {
   domain = "${buddy_workspace.foo.domain}"
   display_name = "%s"
   integration_id = "${buddy_integration.gh.integration_id}"
   external_project_id = "%s"
   update_default_branch_from_external = "%t"
	allow_pull_requests = "%t"
}
`, domain, ghToken, name, ghProject, updateBranch, allowPullRequests)
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
