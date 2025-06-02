package test

import (
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"
)

func TestAccSourceTarget(t *testing.T) {
	domain := util.UniqueString()
	name := util.RandString(10)
	identifier := util.UniqueString()
	host := "1.1.1.1"
	port := "22"
	path := "/"
	tag := util.RandString(3)
	username := util.RandString(10)
	password := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceTargetConfig(domain, name, identifier, host, port, path, username, password, tag),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceTargetAttributes("data.buddy_target.test", name, identifier, host, port, path, tag)),
			},
		},
	})
}

func testAccSourceTargetAttributes(n string, name string, identifier string, host string, port string, path string, tag string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("identifier", attrs["identifier"], identifier); err != nil {
			return err
		}
		if err := util.CheckFieldSet("html_url", attrs["html_url"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("target_id", attrs["target_id"]); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("type", attrs["type"], buddy.TargetTypeSsh); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("host", attrs["host"], host); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("port", attrs["port"], port); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("path", attrs["path"], path); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("tags.0", attrs["tags.0"], tag); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourceTargetConfig(domain string, name string, identifier string, host string, port string, path string, username string, password string, tag string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_target" "test" {
    domain     = buddy_workspace.test.domain
    name       = "%s"
    identifier = "%s"
    type       = "SSH"
    host       = "%s"
    port       = "%s"
    path       = "%s"
    tags       = ["%s"]
    auth {
        method   = "PASSWORD"
        username = "%s"
        password = "%s"
    }
}

data "buddy_target" "test" {
    domain    = buddy_workspace.test.domain
    target_id = buddy_target.test.target_id
}`, domain, name, identifier, host, port, path, tag, username, password)
}
