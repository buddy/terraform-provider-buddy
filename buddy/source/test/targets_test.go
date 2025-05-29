package test

import (
	"fmt"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSourceTargets(t *testing.T) {
	domain := util.UniqueString()
	name1 := util.RandString(10)
	name2 := util.RandString(10)
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
				Config: testAccSourceTargetsConfig(domain, name1, name2, host, port),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.buddy_targets.test", "domain", domain),
					resource.TestCheckResourceAttr("data.buddy_targets.test", "targets.#", "2"),
				),
			},
		},
	})
}

func TestAccSourceTargets_regex(t *testing.T) {
	domain := util.UniqueString()
	prefix := util.RandString(5)
	name1 := fmt.Sprintf("%s-ssh", prefix)
	name2 := fmt.Sprintf("%s-ftp", prefix)
	name3 := util.RandString(10)
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
				Config: testAccSourceTargetsRegexConfig(domain, name1, name2, name3, host, port, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.buddy_targets.test", "domain", domain),
					resource.TestCheckResourceAttr("data.buddy_targets.test", "name_regex", fmt.Sprintf("^%s-", prefix)),
					resource.TestCheckResourceAttr("data.buddy_targets.test", "targets.#", "2"),
				),
			},
		},
	})
}

func testAccSourceTargetsConfig(domain, name1, name2, host, port string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
	domain = "%s"
}

resource "buddy_target" "target1" {
	domain = buddy_workspace.test.domain
	name = "%s"
	type = "SSH"
	host = "%s"
	port = "%s"
	auth_method = "PASS"
	auth_username = "testuser"
	auth_password = "password123"
}

resource "buddy_target" "target2" {
	domain = buddy_workspace.test.domain
	name = "%s"
	type = "FTP"
	host = "%s"
	port = "21"
	auth_username = "ftpuser"
	auth_password = "password123"
}

data "buddy_targets" "test" {
	domain = buddy_workspace.test.domain
	depends_on = [buddy_target.target1, buddy_target.target2]
}
`, domain, name1, host, port, name2, host)
}

func testAccSourceTargetsRegexConfig(domain, name1, name2, name3, host, port, prefix string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "test" {
	domain = "%s"
}

resource "buddy_target" "target1" {
	domain = buddy_workspace.test.domain
	name = "%s"
	type = "SSH"
	host = "%s"
	port = "%s"
	auth_method = "PASS"
	auth_username = "testuser"
	auth_password = "password123"
}

resource "buddy_target" "target2" {
	domain = buddy_workspace.test.domain
	name = "%s"
	type = "FTP"
	host = "%s"
	port = "21"
	auth_username = "ftpuser"
	auth_password = "password123"
}

resource "buddy_target" "target3" {
	domain = buddy_workspace.test.domain
	name = "%s"
	type = "SFTP"
	host = "%s"
	port = "%s"
	auth_method = "PASS"
	auth_username = "sftpuser"
	auth_password = "password123"
}

data "buddy_targets" "test" {
	domain = buddy_workspace.test.domain
	name_regex = "^%s-"
	depends_on = [buddy_target.target1, buddy_target.target2, buddy_target.target3]
}
`, domain, name1, host, port, name2, host, name3, host, port, prefix)
}
