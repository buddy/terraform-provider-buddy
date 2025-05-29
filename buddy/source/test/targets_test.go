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
	hostname := "example.com"
	username := "testuser"
	port := 22

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             acc.DummyCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceTargetsConfig(domain, name1, name2, hostname, port, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.buddy_targets.test", "domain", domain),
					resource.TestCheckResourceAttr("data.buddy_targets.test", "targets.#", "2"),
				),
			},
		},
	})
}

func TestAccSourceTargets_project(t *testing.T) {
	domain := util.UniqueString()
	projectName := util.RandString(10)
	name1 := util.RandString(10)
	name2 := util.RandString(10)
	hostname := "example.com"
	username := "testuser"
	port := 22

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             acc.DummyCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceTargetsProjectConfig(domain, projectName, name1, name2, hostname, port, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.buddy_targets.test", "domain", domain),
					resource.TestCheckResourceAttr("data.buddy_targets.test", "project_name", projectName),
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
	hostname := "example.com"
	username := "testuser"
	port := 22

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             acc.DummyCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceTargetsRegexConfig(domain, name1, name2, name3, hostname, port, username, prefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.buddy_targets.test", "domain", domain),
					resource.TestCheckResourceAttr("data.buddy_targets.test", "name_regex", fmt.Sprintf("^%s-", prefix)),
					resource.TestCheckResourceAttr("data.buddy_targets.test", "targets.#", "2"),
				),
			},
		},
	})
}

func testAccSourceTargetsConfig(domain string, name1 string, name2 string, hostname string, port int, username string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
	domain = "%s"
}

resource "buddy_target" "bar1" {
	domain = buddy_workspace.foo.domain
	name = "%s"
	type = "SSH"
	hostname = "%s"
	port = %d
	username = "%s"
	password = "test123"
}

resource "buddy_target" "bar2" {
	domain = buddy_workspace.foo.domain
	name = "%s"
	type = "FTP"
	hostname = "%s"
	port = 21
	username = "%s"
	password = "test123"
}

data "buddy_targets" "test" {
	domain = buddy_workspace.foo.domain
	depends_on = [buddy_target.bar1, buddy_target.bar2]
}`, domain, name1, hostname, port, username, name2, hostname, username)
}

func testAccSourceTargetsProjectConfig(domain string, projectName string, name1 string, name2 string, hostname string, port int, username string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
	domain = "%s"
}

resource "buddy_project" "proj" {
	domain = buddy_workspace.foo.domain
	display_name = "%s"
}

resource "buddy_target" "bar1" {
	domain = buddy_workspace.foo.domain
	project_name = buddy_project.proj.name
	name = "%s"
	type = "SSH"
	hostname = "%s"
	port = %d
	username = "%s"
	password = "test123"
}

resource "buddy_target" "bar2" {
	domain = buddy_workspace.foo.domain
	project_name = buddy_project.proj.name
	name = "%s"
	type = "FTP"
	hostname = "%s"
	port = 21
	username = "%s"
	password = "test123"
}

data "buddy_targets" "test" {
	domain = buddy_workspace.foo.domain
	project_name = buddy_project.proj.name
	depends_on = [buddy_target.bar1, buddy_target.bar2]
}`, domain, projectName, name1, hostname, port, username, name2, hostname, username)
}

func testAccSourceTargetsRegexConfig(domain string, name1 string, name2 string, name3 string, hostname string, port int, username string, prefix string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
	domain = "%s"
}

resource "buddy_target" "bar1" {
	domain = buddy_workspace.foo.domain
	name = "%s"
	type = "SSH"
	hostname = "%s"
	port = %d
	username = "%s"
	password = "test123"
}

resource "buddy_target" "bar2" {
	domain = buddy_workspace.foo.domain
	name = "%s"
	type = "FTP"
	hostname = "%s"
	port = 21
	username = "%s"
	password = "test123"
}

resource "buddy_target" "bar3" {
	domain = buddy_workspace.foo.domain
	name = "%s"
	type = "SFTP"
	hostname = "%s"
	port = %d
	username = "%s"
	password = "test123"
}

data "buddy_targets" "test" {
	domain = buddy_workspace.foo.domain
	name_regex = "^%s-"
	depends_on = [buddy_target.bar1, buddy_target.bar2, buddy_target.bar3]
}`, domain, name1, hostname, port, username, name2, hostname, username, name3, hostname, port, username, prefix)
}