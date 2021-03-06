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

func TestAccSourceProfile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acc.PreCheck(t) },
		CheckDestroy:      acc.DummyCheckDestroy,
		ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceProfileConfig(),
				Check:  resource.ComposeTestCheckFunc(testAccSourceProfile("data.buddy_profile.me")),
			},
		},
	})
}

func testAccSourceProfileConfig() string {
	return `
data "buddy_profile" "me" {
     
}
`
}

func testAccSourceProfile(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsUserId, _ := strconv.Atoi(attrs["member_id"])
		if err := util.CheckFieldSet("name", attrs["name"]); err != nil {
			return err
		}
		if err := util.CheckIntFieldSet("member_id", attrsUserId); err != nil {
			return err
		}
		if err := util.CheckFieldSet("avatar_url", attrs["avatar_url"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("html_url", attrs["html_url"]); err != nil {
			return err
		}
		return nil
	}
}
