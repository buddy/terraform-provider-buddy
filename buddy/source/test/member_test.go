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

func TestAccSourceMember_upgrade(t *testing.T) {
	domain := util.UniqueString()
	config := testAccSourceMemberConfig(domain)
	p, _, _ := acc.ApiClient.ProfileService.Get()
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
					testAccSourceMemberAttributes("data.buddy_member.id", strconv.Itoa(p.Id)),
					testAccSourceMemberAttributes("data.buddy_member.name", p.Name),
				),
			},
		},
	})
}

func TestAccSourceMember(t *testing.T) {
	domain := util.UniqueString()
	p, _, _ := acc.ApiClient.ProfileService.Get()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceMemberConfig(domain),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceMemberAttributes("data.buddy_member.id", p.Name),
					testAccSourceMemberAttributes("data.buddy_member.name", p.Name),
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
