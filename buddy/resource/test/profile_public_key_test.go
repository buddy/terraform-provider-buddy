package test

import (
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"
)

func TestAccProfilePublicKey_upgrade(t *testing.T) {
	var key buddy.PublicKey
	err, publicKey, _ := util.GenerateRsaKeyPair()
	if err != nil {
		t.Fatal(err.Error())
	}
	content := publicKey
	title := util.RandString(10)
	config := testAccProfilePublicKeyConfig(content, title)
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
					testAccProfilePublicKeyGet("buddy_profile_public_key.foo", &key),
					testAccProfilePublicKeyAttributes("buddy_profile_public_key.foo", &key, content, title),
				),
			},
		},
	})
}

func TestAccProfilePublicKey(t *testing.T) {
	var key buddy.PublicKey
	err, publicKey, _ := util.GenerateRsaKeyPair()
	if err != nil {
		t.Fatal(err.Error())
	}
	err, publicKey2, _ := util.GenerateRsaKeyPair()
	if err != nil {
		t.Fatal(err.Error())
	}
	content := publicKey
	title := util.RandString(10)
	newContent := publicKey2
	newTitle := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccProfilePublicKeyCheckDestroy,
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
