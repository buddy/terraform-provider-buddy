package test

import (
	"buddy-terraform/buddy/acc"
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strconv"
	"testing"
)

func TestAccProfile(t *testing.T) {
	var profile api.Profile
	r := util.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acc.PreCheck(t) },
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      acc.DummyCheckDestroy,
		Steps: []resource.TestStep{
			// update
			{
				Config: testAccProfileUpdateConfig(r),
				Check: resource.ComposeTestCheckFunc(
					testAccProfileGet(&profile),
					testAccProfileAttributes("buddy_profile.me", &profile, &testAccProfileExpectedAttributes{
						Name: fmt.Sprintf("aaaa %d", r),
					}),
				),
			},
			// import
			{
				ResourceName:      "buddy_profile.me",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

type testAccProfileExpectedAttributes struct {
	Name string
}

func testAccProfileAttributes(n string, profile *api.Profile, want *testAccProfileExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsUserId, _ := strconv.Atoi(attrs["member_id"])
		if err := util.CheckFieldEqualAndSet("Name", profile.Name, want.Name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], want.Name); err != nil {
			return err
		}
		if err := util.CheckIntFieldEqualAndSet("member_id", attrsUserId, profile.Id); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("avatar_url", attrs["avatar_url"], profile.AvatarUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("html_url", attrs["html_url"], profile.HtmlUrl); err != nil {
			return err
		}
		return nil
	}
}

func testAccProfileGet(profile *api.Profile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		p, _, err := acc.ApiClient.ProfileService.Get()
		if err != nil {
			return err
		}
		*profile = *p
		return nil
	}
}

func testAccProfileUpdateConfig(r int) string {
	return fmt.Sprintf(`
resource "buddy_profile" "me" {
    name	= "aaaa %d"
}`, r)
}
