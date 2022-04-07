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

func TestAccProfilePublicKey(t *testing.T) {
	var key buddy.PublicKey
	content := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQCG0Ug3U8DoJ6+z36D2h2+oc4UoQRihLNGcAO9SHglFXp+dn1aGJrqeoOrmo4bj5AcydjY33Ylm7ixZEe85vD5INCeldMd8JGmZTj57mwzqpKXFrag+/v9F9qmSEPxKZ1cQj7Q/nRi/hJIoJbsxymrxWhdJZnDNeqwdusR78Xkftw== test"
	title := util.RandString(10)
	newContent := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC5h1SgFvq45BGpYIDowIlaWiGe24kZg2DJ8NYqFo003PcAGdk30oJNvBqJfooGaI7GUkVoCzx9w3Oz/CYmC/NKsz45yUafJRwOBAQ4Gtt1o5RfLNIgj4GlfP1WmCFXIe4cuzureUkPCUIx2K+i1oAOdbEVorzfR3zqPIN/0u3Jwq3nLGmLYS8xCTq3odJT7GvyAj1jyOnXo+dpYZRm6LIteAkhtrnIAI+Le87Bp7JivPZwov/DG7HjW4IuStlTCJOQoYSUtTTu/zBWfSIbmZFakNqIBpiw8vmCeOOgBeOA5u4/JdfNMxH3CP0zxspjPplxkl3DiK/bBs1EGL0zvJrf test2"
	newTitle := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccProfilePublicKeyCheckDestroy,
		Steps: []resource.TestStep{
			// create key
			{
				Config: testAccProfilePublicKeyConfig(content, title),
				Check: resource.ComposeTestCheckFunc(
					testAccProfilePublicKeyGet("buddy_profile_public_key.foo", &key),
					testAccProfilePublicKeyAttributes("buddy_profile_public_key.foo", &key, content, title),
				),
			},
			// update key
			{
				Config: testAccProfilePublicKeyConfig(newContent, newTitle),
				Check: resource.ComposeTestCheckFunc(
					testAccProfilePublicKeyGet("buddy_profile_public_key.foo", &key),
					testAccProfilePublicKeyAttributes("buddy_profile_public_key.foo", &key, newContent, newTitle),
				),
			},
			// import key
			{
				ResourceName:      "buddy_profile_public_key.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccProfilePublicKeyAttributes(n string, key *buddy.PublicKey, content string, title string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		if err := util.CheckFieldEqualAndSet("Content", key.Content, content); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("Title", key.Title, title); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("content", attrs["content"], content); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("title", attrs["title"], title); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("html_url", attrs["html_url"], key.HtmlUrl); err != nil {
			return err
		}
		return nil
	}
}

func testAccProfilePublicKeyGet(n string, key *buddy.PublicKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		keyId, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}
		k, _, err := acc.ApiClient.PublicKeyService.Get(keyId)
		if err != nil {
			return err
		}
		*key = *k
		return nil
	}
}

func testAccProfilePublicKeyConfig(content string, title string) string {
	return fmt.Sprintf(`
resource "buddy_profile_public_key" "foo" {
    content = "%s"
    title = "%s"
}
`, content, title)
}

func testAccProfilePublicKeyCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_profile_public_key" {
			continue
		}
		keyId, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}
		key, resp, err := acc.ApiClient.PublicKeyService.Get(keyId)
		if err == nil && key != nil {
			return util.ErrorResourceExists()
		}
		if resp.StatusCode != 404 {
			return err
		}
	}
	return nil
}
