package test

import (
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"
)

func TestAccSourceTarget_byId(t *testing.T) {
	domain := util.UniqueString()
	name := util.RandString(10)
	identifier := util.UniqueString()
	host := "1.1.1.1"
	port := "44"
	path := util.RandString(10)
	username := util.RandString(10)
	password := util.RandString(10)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceTargetByIdConfig(domain, name, identifier, host, port, path, username, password),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.buddy_target.test", "name", name),
					resource.TestCheckResourceAttr("data.buddy_target.test", "identifier", identifier),
					resource.TestCheckResourceAttr("data.buddy_target.test", "type", buddy.TargetTypeSsh),
					resource.TestCheckResourceAttr("data.buddy_target.test", "host", host),
					resource.TestCheckResourceAttr("data.buddy_target.test", "port", port),
					resource.TestCheckResourceAttr("data.buddy_target.test", "path", path),
					resource.TestCheckResourceAttr("data.buddy_target.test", "secure", "false"),
					resource.TestCheckResourceAttr("data.buddy_target.test", "disabled", "false"),
					resource.TestCheckResourceAttrSet("data.buddy_target.test", "target_id"),
					resource.TestCheckResourceAttrSet("data.buddy_target.test", "html_url"),
				),
			},
		},
	})
}

func TestAccSourceTarget_byIdentifier(t *testing.T) {
	domain := util.UniqueString()
	name := util.RandString(10)
	identifier := util.UniqueString()
	repository := "https://a" + util.UniqueString() + ".com"
	username := util.RandString(10)
	password := util.RandString(10)
	tag1 := "a"
	tag2 := "b"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceTargetByIdentifierConfig(domain, name, identifier, repository, username, password, tag1, tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.buddy_target.test", "name", name),
					resource.TestCheckResourceAttr("data.buddy_target.test", "identifier", identifier),
					resource.TestCheckResourceAttr("data.buddy_target.test", "type", buddy.TargetTypeGit),
					resource.TestCheckResourceAttr("data.buddy_target.test", "repository", repository),
					resource.TestCheckResourceAttr("data.buddy_target.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("data.buddy_target.test", "tags.0", tag1),
					resource.TestCheckResourceAttr("data.buddy_target.test", "tags.1", tag2),
					resource.TestCheckResourceAttrSet("data.buddy_target.test", "target_id"),
					resource.TestCheckResourceAttrSet("data.buddy_target.test", "html_url"),
				),
			},
		},
	})
}

func testAccSourceTargetByIdConfig(domain string, name string, identifier string, host string, port string, path string, username string, password string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_target" "test" {
    domain     = buddy_workspace.test.domain
    name       = "%s"
    identifier = "%s"
    type       = "%s"
    host       = "%s"
    port       = "%s"
    path       = "%s"
    auth {
        method   = "%s"
        username = "%s"
        password = "%s"
    }
}

data "buddy_target" "test" {
    domain    = buddy_workspace.test.domain
    target_id = buddy_target.test.target_id
}`, domain, name, identifier, buddy.TargetTypeSsh, host, port, path, buddy.TargetAuthMethodPassword, username, password)
}

func testAccSourceTargetByIdentifierConfig(domain string, name string, identifier string, repository string, username string, password string, tag1 string, tag2 string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
    domain = "%s"
}

resource "buddy_target" "test" {
    domain     = buddy_workspace.test.domain
    name       = "%s"
    identifier = "%s"
    type       = "%s"
    repository = "%s"
    tags       = ["%s", "%s"]
    auth {
        method   = "%s"
        username = "%s"
        password = "%s"
    }
}

data "buddy_target" "test" {
    domain     = buddy_workspace.test.domain
    identifier = buddy_target.test.identifier
}`, domain, name, identifier, buddy.TargetTypeGit, repository, tag1, tag2, buddy.TargetAuthMethodHttp, username, password)
}
