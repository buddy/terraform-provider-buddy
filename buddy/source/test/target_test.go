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
				Config: testAccSourceTargetConfig(domain, name, hostname, port, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.buddy_target.test", "domain", domain),
					resource.TestCheckResourceAttr("data.buddy_target.test", "name", name),
					resource.TestCheckResourceAttr("data.buddy_target.test", "type", "SSH"),
					resource.TestCheckResourceAttr("data.buddy_target.test", "hostname", hostname),
					resource.TestCheckResourceAttr("data.buddy_target.test", "port", fmt.Sprintf("%d", port)),
					resource.TestCheckResourceAttr("data.buddy_target.test", "username", username),
				),
			},
		},
	})
}

func TestAccSourceTarget_project(t *testing.T) {
	domain := util.UniqueString()
	projectName := util.RandString(10)
	name := util.RandString(10)
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
				Config: testAccSourceTargetProjectConfig(domain, projectName, name, hostname, port, username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.buddy_target.test", "domain", domain),
					resource.TestCheckResourceAttr("data.buddy_target.test", "project_name", projectName),
					resource.TestCheckResourceAttr("data.buddy_target.test", "name", name),
					resource.TestCheckResourceAttr("data.buddy_target.test", "type", "SSH"),
					resource.TestCheckResourceAttr("data.buddy_target.test", "hostname", hostname),
					resource.TestCheckResourceAttr("data.buddy_target.test", "port", fmt.Sprintf("%d", port)),
					resource.TestCheckResourceAttr("data.buddy_target.test", "username", username),
				),
			},
		},
	})
}

func testAccSourceTargetConfig(domain string, name string, hostname string, port int, username string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
	domain = "%s"
}

resource "buddy_target" "bar" {
	domain = buddy_workspace.foo.domain
	name = "%s"
	type = "SSH"
	hostname = "%s"
	port = %d
	username = "%s"
	password = "test123"
}

data "buddy_target" "test" {
	domain = buddy_workspace.foo.domain
	name = buddy_target.bar.name
}`, domain, name, hostname, port, username)
}

func testAccSourceTargetProjectConfig(domain string, projectName string, name string, hostname string, port int, username string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
	domain = "%s"
}

resource "buddy_project" "proj" {
	domain = buddy_workspace.foo.domain
	display_name = "%s"
}

resource "buddy_target" "bar" {
	domain = buddy_workspace.foo.domain
	project_name = buddy_project.proj.name
	name = "%s"
	type = "SSH"
	hostname = "%s"
	port = %d
	username = "%s"
	password = "test123"
}

data "buddy_target" "test" {
	domain = buddy_workspace.foo.domain
	project_name = buddy_project.proj.name
	name = buddy_target.bar.name
}`, domain, projectName, name, hostname, port, username)
}