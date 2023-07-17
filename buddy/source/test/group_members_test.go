package test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"
)

func TestAccSourceGroupMembers_upgrade(t *testing.T) {
	domain := util.UniqueString()
	groupName := util.UniqueString()
	email1 := util.RandEmail()
	email2 := util.RandEmail()
	config := testAccSourceGroupMembersConfig(domain, groupName, email1, email2)
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
					testAccSourceGroupMembersAttributes("data.buddy_group_members.gm", 3),
					testAccSourceGroupMembersAttributes("data.buddy_group_members.filter", 1),
				),
			},
		},
	})
}

func TestAccSourceGroupMembers(t *testing.T) {
	domain := util.UniqueString()
	groupName := util.UniqueString()
	email1 := util.RandEmail()
	email2 := util.RandEmail()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceGroupMembersConfig(domain, groupName, email1, email2),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceGroupMembersAttributes("data.buddy_group_members.gm", 3),
					testAccSourceGroupMembersAttributes("data.buddy_group_members.filter", 1),
				),
			},
		},
	})
}

func testAccSourceGroupMembersAttributes(n string, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsMembersCount, _ := strconv.Atoi(attrs["members.#"])
		attrsMemberId, _ := strconv.Atoi(attrs["members.0.member_id"])
		if err := util.CheckIntFieldEqual("members.#", attrsMembersCount, count); err != nil {
			return err
		}
		if err := util.CheckFieldSet("members.0.html_url", attrs["members.0.html_url"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("members.0.avatar_url", attrs["members.0.avatar_url"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("members.0.email", attrs["members.0.email"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("members.0.status", attrs["members.0.status"]); err != nil {
			return err
		}
		if err := util.CheckIntFieldSet("members.0.member_id", attrsMemberId); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourceGroupMembersConfig(domain string, groupName string, email1 string, email2 string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "w" {
   domain = "%s"
}

resource "buddy_group" "g" {
   domain = "${buddy_workspace.w.domain}"
   name = "%s"
}

resource "buddy_member" "m1" {
   domain = "${buddy_workspace.w.domain}"
   email ="%s"
}

resource "buddy_member" "m2" {
   domain = "${buddy_workspace.w.domain}"
   email ="%s"
}

data "buddy_profile" "me" {
}

resource "buddy_group_member" "gm1" {
   domain = "${buddy_workspace.w.domain}"
   group_id = "${buddy_group.g.group_id}"
   member_id = "${buddy_member.m1.member_id}"
}

resource "buddy_group_member" "gm2" {
   domain = "${buddy_workspace.w.domain}"
   group_id = "${buddy_group.g.group_id}"
   member_id = "${buddy_member.m2.member_id}"
	status = "MEMBER"
}

resource "buddy_group_member" "gm3" {
   domain = "${buddy_workspace.w.domain}"
   group_id = "${buddy_group.g.group_id}"
   member_id = "${data.buddy_profile.me.member_id}"
	status = "MANAGER"
}

data "buddy_group_members" "gm" {
   domain = "${buddy_workspace.w.domain}"
   group_id = "${buddy_group.g.group_id}"
   depends_on = [buddy_group_member.gm1, buddy_group_member.gm2, buddy_group_member.gm3]
}

data "buddy_group_members" "filter" {
   domain = "${buddy_workspace.w.domain}"
   group_id = "${buddy_group.g.group_id}"
	 name_regex = "^${data.buddy_profile.me.name}"
   depends_on = [buddy_group_member.gm1, buddy_group_member.gm2, buddy_group_member.gm3]
}
`, domain, groupName, email1, email2)
}
