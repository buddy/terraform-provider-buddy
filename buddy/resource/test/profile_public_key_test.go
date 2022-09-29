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
	content := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQDNhGKQNVsHFpx5/doxM74M0KsZqEA3TRN9x31jpESUrJsamBD9WGSNywlexfszqWb0OcDxqAA/ZthtxOrwvO8BPxioCiXwQ9s4ASE3375ueilxyp4W6MQ8v3vNB1ossTyU6LQGoIkKb4Km3s962e14RB5IUO9XsF5Vt9u+f8Y0J2+gHQyXqBClrTl0J7xJYzdhW9U00p4OFQmp4S5JM4zxO9M/XTF8vXdKUMmcRHmkZiJS4BNlTX5oFoC7nKF7nqSczmBYnjO6b6YhdpueHDY76iqdi+R1uwVt/kA/z5n7hHGxJYB6Vx8TNTg0QsmrhmkDd3DcdiHSfwVXUc0CDXzkOu1WbwvmNLv9XGvH2coJ3QIFPLgLUe5uP7jr3TwlJyERWxe4VJxqTVbPPV5Eik9zEfroXyhne2tWATu7aaCge0nbSBJVpmKTFbHAzeP0+sMPJS7Rv2jkKRar3LmJhM5/E5ditqD40dQfYtvRK2/utD3eShiOyUbjZYTZazjVaREdeqcZUjywlVSyHa37TNtaH5IiDQMA0ClmEnIaQCPYMjtmXwWbpdP5cqyNUpMD3Rk+cw7VXf5QoV+uij7CzlhLrSjgmJyDeAh+d9uvW3TJeWJsHdr0fd6VXtdnkZQqpfEuL95V4f5tiDF6nKSP7wTgxZWsOGZ3xKAnEJ7P5Rk6LQ== test"
	title := util.RandString(10)
	newContent := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCsn4Z8puW9BDiJm9yE9rCDQxiL1yT6JcNf988is3zk8n0g9tTqSVJo6RsfHpJIkMo8EdYslsr6VCyELjWAZqrJ4lCdpgiVHRpgleDOBrq+7SObSqE5xn5Rn6NqI04vpZDoyaQW65H+5S9kwWAcDJujS/0EO1turAb+jZDjVPHUn4aHllmSjWbhvI7KkMh34Y1SFNVNI2bSH6tFcVgr5UI4/HkoCzsl+yzhYVDUALwO+3QbJsH+r+uKvI6MgGAecUb1lJJoKQBdQlC/aWTlJ+kjbGzo4hV27H36LCGVxrZbNOfjMvInu5Tpie4qWdeJVOWZl5jVaDG0jdifTUqod9NGgLZdV5LpB0EMBnd1RNCNYTTlqe7/WWtjadbIME47SNyeteulF51UcON6ZUKjvgiEfc6++TPfi9AkHvhpbpJHNBnl419oqS8Sqvz/OUj1G56L9AMQG5WBrwyC9v6QRHygaUtvGKrBbI3/OFifnCkjslYe5MZZ2AimDRFSo9sYSONl/pKUdnwokpO7j6gJlpxnTAc3XhKNGGTMR0nsFl/t1kcxut/XFOe6NItBLMUOJNHqrcjlDvyWoNtwOllv/EoY422xmy4Ue1oT8gFaHM+k/9x/mSuc38p1Lo7VziL549CPR0Rv+aDAJSR1TQqjBvt6/AP7Yov3cG9x1jpDKq0zfw== test2"
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
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}
