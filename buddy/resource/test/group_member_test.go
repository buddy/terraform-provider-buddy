package test

import (
	"buddy-terraform/buddy/acc"
	"buddy-terraform/buddy/util"
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"testing"
)

// todo test upgrade group member

func TestAccGroupMember(t *testing.T) {
	var member buddy.Member
	domain := util.UniqueString()
	groupNameA := util.RandString(5)
	groupNameB := util.RandString(5)
	memberEmailA := util.RandEmail()
	memberEmailB := util.RandEmail()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccGroupMemberDestroy,
		Steps: []resource.TestStep{
			// create group member
			{
				Config: testAccGroupMemberConfig(domain, groupNameA, groupNameB, memberEmailA, memberEmailB, "a", "a"),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupMemberGet("buddy_group_member.bar", &member),
					testAccGroupMemberAttributes("buddy_group_member.bar", &member, buddy.GroupMemberStatusMember),
				),
			},
			// update group
			{
				Config: testAccGroupMemberConfig(domain, groupNameA, groupNameB, memberEmailA, memberEmailB, "a", "b"),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupMemberGet("buddy_group_member.bar", &member),
					testAccGroupMemberAttributes("buddy_group_member.bar", &member, buddy.GroupMemberStatusMember),
				),
			},
			// update member
			{
				Config: testAccGroupMemberConfig(domain, groupNameA, groupNameB, memberEmailA, memberEmailB, "b", "b"),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupMemberGet("buddy_group_member.bar", &member),
					testAccGroupMemberAttributes("buddy_group_member.bar", &member, buddy.GroupMemberStatusMember),
				),
			},
			// add manager
			{
				Config: testAccGroupMemberWithStatusConfig(domain, groupNameA, memberEmailA, buddy.GroupMemberStatusManager),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupMemberGet("buddy_group_member.bar", &member),
					testAccGroupMemberAttributes("buddy_group_member.bar", &member, buddy.GroupMemberStatusManager),
				),
			},
			// change to member
			{
				Config: testAccGroupMemberWithStatusConfig(domain, groupNameA, memberEmailA, buddy.GroupMemberStatusMember),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupMemberGet("buddy_group_member.bar", &member),
					testAccGroupMemberAttributes("buddy_group_member.bar", &member, buddy.GroupMemberStatusMember),
				),
			},
			// import group member
			{
				ResourceName:      "buddy_group_member.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGroupMemberAttributes(n string, member *buddy.Member, status string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		if err := util.CheckFieldEqualAndSet("html_url", attrs["html_url"], member.HtmlUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("email", attrs["email"], member.Email); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("avatar_url", attrs["avatar_url"], member.AvatarUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("admin", attrs["admin"], strconv.FormatBool(member.Admin)); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("workspace_owner", attrs["workspace_owner"], strconv.FormatBool(member.WorkspaceOwner)); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("status", attrs["status"], status); err != nil {
			return err
		}
		return nil
	}
}

func testAccGroupMemberGet(n string, member *buddy.Member) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		domain, gid, mid, err := util.DecomposeTripleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		groupId, err := strconv.Atoi(gid)
		if err != nil {
			return err
		}
		memberId, err := strconv.Atoi(mid)
		if err != nil {
			return err
		}
		m, _, err := acc.ApiClient.GroupService.GetGroupMember(domain, groupId, memberId)
		if err != nil {
			return err
		}
		*member = *m
		return nil
	}
}

func testAccGroupMemberWithStatusConfig(domain string, groupName string, memberEmail string, status string) string {
	return fmt.Sprintf(`

	resource "buddy_workspace" "foo" {
	   domain = "%s"
	}

	resource "buddy_group" "g" {
	   domain = "${buddy_workspace.foo.domain}"
	   name = "%s"
	}

	resource "buddy_member" "m" {
	   domain = "${buddy_workspace.foo.domain}"
	   email = "%s"
	}

	resource "buddy_group_member" "bar" {
	   domain = "${buddy_workspace.foo.domain}"
	   group_id = "${buddy_group.g.group_id}"
	   member_id = "${buddy_member.m.member_id}"
		status = "%s"
	}

`, domain, groupName, memberEmail, status)
}

func testAccGroupMemberConfig(domain string, groupNameA string, groupNameB string, memberEmailA string, memberEmailB string, whichMemberJoin string, whichGroupJoin string) string {
	return fmt.Sprintf(`

	resource "buddy_workspace" "foo" {
	   domain = "%s"
	}

	resource "buddy_group" "a" {
	   domain = "${buddy_workspace.foo.domain}"
	   name = "%s"
	}

	resource "buddy_group" "b" {
	   domain = "${buddy_workspace.foo.domain}"
	   name = "%s"
	}

	resource "buddy_member" "a" {
	   domain = "${buddy_workspace.foo.domain}"
	   email = "%s"
	}

	resource "buddy_member" "b" {
	   domain = "${buddy_workspace.foo.domain}"
	   email = "%s"
	}

	resource "buddy_group_member" "bar" {
	   domain = "${buddy_workspace.foo.domain}"
	   group_id = "${buddy_group.%s.group_id}"
	   member_id = "${buddy_member.%s.member_id}"
	}

`, domain, groupNameA, groupNameB, memberEmailA, memberEmailB, whichGroupJoin, whichMemberJoin)
}

func testAccGroupMemberDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_group_member" {
			continue
		}
		domain, gid, mid, err := util.DecomposeTripleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		groupId, err := strconv.Atoi(gid)
		if err != nil {
			return err
		}
		memberId, err := strconv.Atoi(mid)
		if err != nil {
			return err
		}
		member, resp, err := acc.ApiClient.GroupService.GetGroupMember(domain, groupId, memberId)
		if err == nil && member != nil {
			return util.ErrorResourceExists()
		}
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}
