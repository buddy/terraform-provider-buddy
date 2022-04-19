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

func TestAccMember(t *testing.T) {
	var member buddy.Member
	var permission buddy.Permission
	domain := util.UniqueString()
	email := util.RandEmail()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccMemberCheckDestroy,
		Steps: []resource.TestStep{
			// create member
			{
				Config: testAccMemberConfig(domain, email),
				Check: resource.ComposeTestCheckFunc(
					testAccMemberGet("buddy_member.bar", &member),
					testAccMemberAttributes("buddy_member.bar", &member, false, email, nil),
				),
			},
			// update member
			{
				Config: testAccMemberUpdateConfig(domain, email),
				Check: resource.ComposeTestCheckFunc(
					testAccMemberGet("buddy_member.bar", &member),
					testAccMemberAttributes("buddy_member.bar", &member, true, email, nil),
				),
			},
			// update auto assign
			{
				Config: testAccMemberUpdateAutoAssignConfig(domain, email),
				Check: resource.ComposeTestCheckFunc(
					testAccMemberGet("buddy_member.bar", &member),
					testAccPermissionGet("buddy_permission.perm", &permission),
					testAccMemberAttributes("buddy_member.bar", &member, true, email, &permission),
				),
			},
			// update member
			{
				Config: testAccMemberConfig(domain, email),
				Check: resource.ComposeTestCheckFunc(
					testAccMemberGet("buddy_member.bar", &member),
					testAccMemberAttributes("buddy_member.bar", &member, false, email, nil),
				),
			},
			// import member
			{
				ResourceName:      "buddy_member.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccMemberAttributes(n string, member *buddy.Member, admin bool, email string, permission *buddy.Permission) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsAdmin, _ := strconv.ParseBool(attrs["admin"])
		attrsOwner, _ := strconv.ParseBool(attrs["workspace_owner"])
		attrsMemberId, _ := strconv.Atoi(attrs["member_id"])
		attrsAutoAssignToProjects, _ := strconv.ParseBool(attrs["auto_assign_to_new_projects"])
		attrsAutoAssignToProjectsPermissionId, _ := strconv.Atoi(attrs["auto_assign_permission_set_id"])
		if err := util.CheckBoolFieldEqual("Admin", member.Admin, admin); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("Email", member.Email, email); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("admin", attrsAdmin, admin); err != nil {
			return err
		}
		if err := util.CheckIntFieldEqualAndSet("member_id", attrsMemberId, member.Id); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("html_url", attrs["html_url"], member.HtmlUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("avatar_url", attrs["avatar_url"], member.AvatarUrl); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("workspace_owner", attrsOwner, member.WorkspaceOwner); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("email", attrs["email"], email); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("auto_assign_to_new_projects", attrsAutoAssignToProjects, member.AutoAssignToNewProjects); err != nil {
			return err
		}
		if err := util.CheckIntFieldEqual("auto_assign_permission_set_id", attrsAutoAssignToProjectsPermissionId, member.AutoAssignPermissionSetId); err != nil {
			return err
		}
		if permission != nil {
			if err := util.CheckBoolFieldEqual("AutoAssignToNewProjects", member.AutoAssignToNewProjects, true); err != nil {
				return err
			}
			if err := util.CheckIntFieldEqual("AutoAssignPermissionSetId", member.AutoAssignPermissionSetId, permission.Id); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccMemberGet(n string, member *buddy.Member) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		domain, mid, err := util.DecomposeDoubleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		memberId, err := strconv.Atoi(mid)
		if err != nil {
			return err
		}
		m, _, err := acc.ApiClient.MemberService.Get(domain, memberId)
		if err != nil {
			return err
		}
		*member = *m
		return nil
	}
}

func testAccMemberUpdateConfig(domain string, email string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_permission" "perm" {
    domain = "${buddy_workspace.foo.domain}"
    name = "test"
    pipeline_access_level = "READ_ONLY"
    repository_access_level = "READ_ONLY"
	sandbox_access_level = "READ_ONLY"
}

resource "buddy_member" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    email = "%s"
    admin = true
}
`, domain, email)
}

func testAccMemberUpdateAutoAssignConfig(domain string, email string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_permission" "perm" {
    domain = "${buddy_workspace.foo.domain}"
    name = "test"
    pipeline_access_level = "READ_ONLY"
    repository_access_level = "READ_ONLY"
	sandbox_access_level = "READ_ONLY"
}

resource "buddy_member" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    email = "%s"
    admin = true
	auto_assign_to_new_projects = true
	auto_assign_permission_set_id = "${buddy_permission.perm.permission_id}"
}
`, domain, email)
}

func testAccMemberConfig(domain string, email string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_permission" "perm" {
    domain = "${buddy_workspace.foo.domain}"
    name = "test"
    pipeline_access_level = "READ_ONLY"
    repository_access_level = "READ_ONLY"
	sandbox_access_level = "READ_ONLY"
}

resource "buddy_member" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    email = "%s"
}
`, domain, email)
}

func testAccMemberCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_member" {
			continue
		}
		domain, mid, err := util.DecomposeDoubleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		memberId, err := strconv.Atoi(mid)
		if err != nil {
			return err
		}
		member, resp, err := acc.ApiClient.MemberService.Get(domain, memberId)
		if err == nil && member != nil {
			return util.ErrorResourceExists()
		}
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}
