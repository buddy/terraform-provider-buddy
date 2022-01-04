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

func TestAccSourceMember(t *testing.T) {
	domain := util.UniqueString()
	name := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:      acc.DummyCheckDestroy,
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceMemberConfig(domain, name),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceMemberAttributes("data.buddy_member.id", name),
					testAccSourceMemberAttributes("data.buddy_member.name", name),
				),
			},
		},
	})
}

func testAccSourceMemberAttributes(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsMemberId, _ := strconv.Atoi(attrs["member_id"])
		attrsAdmin, _ := strconv.ParseBool(attrs["admin"])
		attrsOwner, _ := strconv.ParseBool(attrs["workspace_owner"])
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], name); err != nil {
			return err
		}
		if err := util.CheckIntFieldSet("member_id", attrsMemberId); err != nil {
			return err
		}
		if err := util.CheckFieldSet("html_url", attrs["html_url"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("avatar_url", attrs["avatar_url"]); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("workspace_owner", attrsOwner, true); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("admin", attrsAdmin, true); err != nil {
			return err
		}
		if err := util.CheckFieldSet("email", attrs["email"]); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourceMemberConfig(domain string, name string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_profile" "me" {
    name = "%s"
}

data "buddy_member" "id" {
    domain = "${buddy_workspace.foo.domain}"
    member_id = "${buddy_profile.me.member_id}"
}

data "buddy_member" "name" {
    domain = "${buddy_workspace.foo.domain}"
    name = "${buddy_profile.me.name}"
}

`, domain, name)
}
