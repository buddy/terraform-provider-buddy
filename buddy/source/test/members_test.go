package test

import (
	"buddy-terraform/buddy/acc"
	"buddy-terraform/buddy/util"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strconv"
	"testing"
)

func TestAccSourceMembers(t *testing.T) {
	domain := util.UniqueString()
	email1 := util.RandEmail()
	email2 := util.RandEmail()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:      acc.DummyCheckDestroy,
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceMembersConfig(domain, email1, email2),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceMembersAttributes("data.buddy_members.m", 3),
					testAccSourceMembersAttributes("data.buddy_members.filter", 1),
				),
			},
		},
	})
}

func testAccSourceMembersAttributes(n string, count int) resource.TestCheckFunc {
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
		if err := util.CheckIntFieldSet("members.0.member_id", attrsMemberId); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourceMembersConfig(domain string, email1 string, email2 string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "w" {
    domain = "%s"
}

resource "buddy_member" "m1" {
    domain = "${buddy_workspace.w.domain}"
    email ="%s"
}

resource "buddy_member" "m2" {
    domain = "${buddy_workspace.w.domain}"
    email ="%s"
}

resource "buddy_profile" "me" {
	name = "abcdef"
}

data "buddy_members" "m" {
    domain = "${buddy_workspace.w.domain}"
    depends_on = [buddy_member.m1, buddy_member.m2, buddy_profile.me]
}

data "buddy_members" "filter" {
    domain = "${buddy_workspace.w.domain}"
    depends_on = [buddy_member.m1, buddy_member.m2, buddy_profile.me]
	name_regex = "^abc"
}
`, domain, email1, email2)
}
