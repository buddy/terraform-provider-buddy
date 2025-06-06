package test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"
)

func TestAccSourceEnvironment(t *testing.T) {
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name := util.UniqueString()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceEnvironmentConfig(domain, projectName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceEnvironmentAttributes("data.buddy_environment.id", name),
					testAccSourceEnvironmentAttributes("data.buddy_environment.name", name),
				),
			},
		},
	})
}

func testAccSourceEnvironmentAttributes(n string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("identifier", attrs["identifier"], name); err != nil {
			return err
		}
		if err := util.CheckFieldSet("html_url", attrs["html_url"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("environment_id", attrs["environment_id"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("public_url", attrs["public_url"]); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("tags.0", attrs["tags.0"], "a"); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourceEnvironmentConfig(domain string, projectName string, name string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_project" "proj" {
   domain = "${buddy_workspace.foo.domain}"
   display_name = "%s"
}

resource "buddy_environment" "a" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   name = "%s"
   identifier = "%s"
	 public_url = "https://a.com"
   tags = ["a"]
}

data "buddy_environment" "id" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   environment_id = "${buddy_environment.a.environment_id}"
}

data "buddy_environment" "name" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   name = "${buddy_environment.a.name}"
}
`, domain, projectName, name, name)
}
