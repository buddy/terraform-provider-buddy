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

func TestAccSourceMember(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceMemberConfig(util.UniqueString()),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceMemberAttributesCheck("data.buddy_member.id"),
					testAccSourceMemberAttributesCheck("data.buddy_member.name"),
				),
			},
		},
	})
}

func testAccSourceMemberAttributesCheck(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsMemberId, _ := strconv.Atoi(attrs["member_id"])
		attrsAdmin, _ := strconv.ParseBool(attrs["admin"])
		attrsOwner, _ := strconv.ParseBool(attrs["workspace_owner"])
		if err := util.CheckFieldSet("name", attrs["name"]); err != nil {
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

func testAccSourceMemberConfig(domain string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

data "buddy_profile" "me" {
}

data "buddy_member" "id" {
   domain = "${buddy_workspace.foo.domain}"
   member_id = "${data.buddy_profile.me.member_id}"
}

data "buddy_member" "name" {
   domain = "${buddy_workspace.foo.domain}"
   name = "${data.buddy_profile.me.name}"
}

`, domain)
}
