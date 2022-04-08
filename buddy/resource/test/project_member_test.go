package test

import (
	"buddy-terraform/buddy/acc"
	"buddy-terraform/buddy/util"
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strconv"
	"testing"
)

func TestAccProjectMember(t *testing.T) {
	var member buddy.ProjectMember
	domain := util.UniqueString()
	emailA := util.RandEmail()
	emailB := util.RandEmail()
	projectDisplayNameA := util.RandString(10)
	projectDisplayNameB := util.RandString(10)
	permissionNameA := util.RandString(10)
	permissionNameB := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccProjectMemberCheckDestroy,
		Steps: []resource.TestStep{
			// create
			{
				Config: testAccProjectMemberConfig(domain, emailA, emailB, projectDisplayNameA, projectDisplayNameB, permissionNameA, permissionNameB, "a", "a", "a"),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectMemberGet("buddy_project_member.bar", &member),
					testAccProjectMemberAttributes("buddy_project_member.bar", &member, permissionNameA),
				),
			},
			// update member
			{
				Config: testAccProjectMemberConfig(domain, emailA, emailB, projectDisplayNameA, projectDisplayNameB, permissionNameA, permissionNameB, "a", "b", "a"),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectMemberGet("buddy_project_member.bar", &member),
					testAccProjectMemberAttributes("buddy_project_member.bar", &member, permissionNameA),
				),
			},
			// update project
			{
				Config: testAccProjectMemberConfig(domain, emailA, emailB, projectDisplayNameA, projectDisplayNameB, permissionNameA, permissionNameB, "b", "b", "a"),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectMemberGet("buddy_project_member.bar", &member),
					testAccProjectMemberAttributes("buddy_project_member.bar", &member, permissionNameA),
				),
			},
			// update permission
			{
				Config: testAccProjectMemberConfig(domain, emailA, emailB, projectDisplayNameA, projectDisplayNameB, permissionNameA, permissionNameB, "b", "b", "b"),
				Check: resource.ComposeTestCheckFunc(
					testAccProjectMemberGet("buddy_project_member.bar", &member),
					testAccProjectMemberAttributes("buddy_project_member.bar", &member, permissionNameB),
				),
			},
			// import
			{
				ResourceName:      "buddy_project_member.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProjectMemberAttributes(n string, member *buddy.ProjectMember, permissionName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsMemberId, _ := strconv.Atoi(attrs["member_id"])
		attrsPermissionId, _ := strconv.Atoi(attrs["permission_id"])
		attrsAdmin, _ := strconv.ParseBool(attrs["admin"])
		attrsOwner, _ := strconv.ParseBool(attrs["workspace_owner"])
		attrsPermissionPermissionId, _ := strconv.Atoi(attrs["permission.0.permission_id"])
		if err := util.CheckIntFieldEqualAndSet("member_id", attrsMemberId, member.Id); err != nil {
			return err
		}
		if err := util.CheckIntFieldEqualAndSet("permission_id", attrsPermissionId, member.PermissionSet.Id); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("html_url", attrs["html_url"], member.HtmlUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("email", attrs["email"], member.Email); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("avatar_url", attrs["avatar_url"], member.AvatarUrl); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("admin", attrsAdmin, member.Admin); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("workspace_owner", attrsOwner, member.WorkspaceOwner); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("permission.0.html_url", attrs["permission.0.html_url"], member.PermissionSet.HtmlUrl); err != nil {
			return err
		}
		if err := util.CheckIntFieldEqualAndSet("permission.0.permission_id", attrsPermissionPermissionId, member.PermissionSet.Id); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("permission.0.name", attrs["permission.0.name"], member.PermissionSet.Name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("permission.0.name", attrs["permission.0.name"], permissionName); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("permission.0.type", attrs["permission.0.type"], member.PermissionSet.Type); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("permission.0.pipeline_access_level", attrs["permission.0.pipeline_access_level"], member.PermissionSet.PipelineAccessLevel); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("permission.0.repository_access_level", attrs["permission.0.repository_access_level"], member.PermissionSet.RepositoryAccessLevel); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("permission.0.sandbox_access_level", attrs["permission.0.sandbox_access_level"], member.PermissionSet.SandboxAccessLevel); err != nil {
			return err
		}
		return nil
	}
}

func testAccProjectMemberGet(n string, member *buddy.ProjectMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		domain, projectName, mid, err := util.DecomposeTripleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		memberId, err := strconv.Atoi(mid)
		if err != nil {
			return err
		}
		m, _, err := acc.ApiClient.ProjectMemberService.GetProjectMember(domain, projectName, memberId)
		if err != nil {
			return err
		}
		*member = *m
		return nil
	}
}

func testAccProjectMemberConfig(domain string, emailA string, emailB string, projectDisplayNameA string, projectDisplayNameB string, permissionNameA string, permissionNameB string, whichProject string, whichMember string, whichPermission string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_member" "a" {
    domain = "${buddy_workspace.foo.domain}"
    email = "%s"
    admin = true
}

resource "buddy_member" "b" {
    domain = "${buddy_workspace.foo.domain}"
    email = "%s"
}

resource "buddy_project" "a" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_project" "b" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_permission" "a" {
    domain = "${buddy_workspace.foo.domain}"
    name = "%s"
    pipeline_access_level = "%s"
    repository_access_level = "%s"
	sandbox_access_level = "%s"
}

resource "buddy_permission" "b" {
    domain = "${buddy_workspace.foo.domain}"
    name = "%s"
    pipeline_access_level = "%s"
    repository_access_level = "%s"
	sandbox_access_level = "%s"
}

resource "buddy_project_member" "bar" {
	domain = "${buddy_workspace.foo.domain}"
	project_name = "${buddy_project.%s.name}"
	member_id = "${buddy_member.%s.member_id}"
	permission_id = "${buddy_permission.%s.permission_id}"
}

`,
		domain,
		emailA,
		emailB,
		projectDisplayNameA,
		projectDisplayNameB,
		permissionNameA,
		buddy.PermissionAccessLevelReadOnly,
		buddy.PermissionAccessLevelReadOnly,
		buddy.PermissionAccessLevelReadOnly,
		permissionNameB,
		buddy.PermissionAccessLevelReadWrite,
		buddy.PermissionAccessLevelReadWrite,
		buddy.PermissionAccessLevelReadWrite,
		whichProject,
		whichMember,
		whichPermission,
	)
}

func testAccProjectMemberCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_project_member" {
			continue
		}
		domain, projectName, mid, err := util.DecomposeTripleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		memberId, err := strconv.Atoi(mid)
		if err != nil {
			return err
		}
		member, resp, err := acc.ApiClient.ProjectMemberService.GetProjectMember(domain, projectName, memberId)
		if err == nil && member != nil {
			return util.ErrorResourceExists()
		}
		if resp.StatusCode != 404 {
			return err
		}
	}
	return nil
}
