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

func TestAccProfileEmail(t *testing.T) {
	var pe api.ProfileEmail
	email := util.RandEmail()
	newEmail := util.RandEmail()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccProfileEmailCheckDestroy,
		Steps: []resource.TestStep{
			// create email
			{
				Config: testAccProfileEmailConfig(email),
				Check: resource.ComposeTestCheckFunc(
					testAccProfileEmailGet("buddy_profile_email.foo", &pe),
					testAccProfileEmailAttributes("buddy_profile_email.foo", &pe, email),
				),
			},
			// update email
			{
				Config: testAccProfileEmailConfig(newEmail),
				Check: resource.ComposeTestCheckFunc(
					testAccProfileEmailGet("buddy_profile_email.foo", &pe),
					testAccProfileEmailAttributes("buddy_profile_email.foo", &pe, newEmail),
				),
			},
			// import email
			{
				ResourceName:      "buddy_profile_email.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProfileEmailAttributes(n string, pe *api.ProfileEmail, email string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsConfirmed, _ := strconv.ParseBool(attrs["confirmed"])
		if err := util.CheckFieldEqualAndSet("Email", pe.Email, email); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("Confirmed", pe.Confirmed, false); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("email", attrs["email"], email); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("confirmed", attrsConfirmed, pe.Confirmed); err != nil {
			return err
		}
		return nil
	}
}

func testAccProfileEmailGet(n string, pe *api.ProfileEmail) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		list, _, err := acc.ApiClient.ProfileEmailService.GetList()
		if err != nil {
			return err
		}
		email := rs.Primary.ID
		for _, v := range list.Emails {
			if v.Email == email {
				*pe = *v
				return nil
			}
		}
		return fmt.Errorf("profile email not found")
	}
}

func testAccProfileEmailConfig(email string) string {
	return fmt.Sprintf(`
resource "buddy_profile_email" "foo" {
    email = "%s"
}
`, email)
}

func testAccProfileEmailCheckDestroy(s *terraform.State) error {
	list, _, err := acc.ApiClient.ProfileEmailService.GetList()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_profile_email" {
			continue
		}
		email := rs.Primary.ID
		for _, v := range list.Emails {
			if v.Email == email {
				return util.ErrorResourceExists()
			}
		}
	}
	return nil
}
