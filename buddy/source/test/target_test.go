package test

import (
	"fmt"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSourceTarget(t *testing.T) {
	domain := util.UniqueString()
	name := util.RandString(10)
	host := "example.com"
	port := "22"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             acc.DummyCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceTargetConfig(domain, name, host, port),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.buddy_target.by_id", "domain", domain),
					resource.TestCheckResourceAttr("data.buddy_target.by_id", "name", name),
					resource.TestCheckResourceAttr("data.buddy_target.by_id", "type", "SSH"),
					resource.TestCheckResourceAttr("data.buddy_target.by_id", "host", host),
					resource.TestCheckResourceAttr("data.buddy_target.by_id", "port", port),
					resource.TestCheckResourceAttrSet("data.buddy_target.by_id", "target_id"),
					resource.TestCheckResourceAttrSet("data.buddy_target.by_id", "html_url"),
					resource.TestCheckResourceAttr("data.buddy_target.by_name", "domain", domain),
					resource.TestCheckResourceAttr("data.buddy_target.by_name", "name", name),
					resource.TestCheckResourceAttr("data.buddy_target.by_name", "type", "SSH"),
					resource.TestCheckResourceAttr("data.buddy_target.by_name", "host", host),
					resource.TestCheckResourceAttr("data.buddy_target.by_name", "port", port),
					resource.TestCheckResourceAttrSet("data.buddy_target.by_name", "target_id"),
					resource.TestCheckResourceAttrSet("data.buddy_target.by_name", "html_url"),
				),
			},
		},
	})
}

func testAccSourceTargetConfig(domain, name, host, port string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
	domain = "%s"
}

resource "buddy_target" "test" {
	domain = buddy_workspace.test.domain
	name = "%s"
	type = "SSH"
	host = "%s"
	port = "%s"
	auth_method = "PASS"
	auth_username = "testuser"
	auth_password = "password123"
}

data "buddy_target" "by_id" {
	domain = buddy_workspace.test.domain
	target_id = buddy_target.test.target_id
}

data "buddy_target" "by_name" {
	domain = buddy_workspace.test.domain
	name = buddy_target.test.name
}
`, domain, name, host, port)
}
