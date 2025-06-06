package test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"strconv"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"
)

func TestAccSourceEnvironments(t *testing.T) {
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name1 := "aaaa" + util.UniqueString()
	name2 := util.UniqueString()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceEnvironmentsConfig(domain, projectName, name1, name2),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceEnvironmentsAttributes("data.buddy_environments.all", 2, ""),
					testAccSourceEnvironmentsAttributes("data.buddy_environments.name", 1, name1),
				),
			},
		},
	})
}

func testAccSourceEnvironmentsAttributes(n string, count int, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsEnvironmentsCount, _ := strconv.Atoi(attrs["environments.#"])
		if err := util.CheckIntFieldEqual("environments.#", attrsEnvironmentsCount, count); err != nil {
			return err
		}
		if count > 0 {
			if name != "" {
				if err := util.CheckFieldEqualAndSet("environments.0.name", attrs["environments.0.name"], name); err != nil {
					return err
				}
				if err := util.CheckFieldEqualAndSet("environments.0.identifier", attrs["environments.0.identifier"], name); err != nil {
					return err
				}
			} else {
				if err := util.CheckFieldSet("environments.0.name", attrs["environments.0.name"]); err != nil {
					return err
				}
				if err := util.CheckFieldSet("environments.0.identifier", attrs["environments.0.identifier"]); err != nil {
					return err
				}
			}
			if err := util.CheckFieldSet("environments.0.html_url", attrs["environments.0.html_url"]); err != nil {
				return err
			}
			if err := util.CheckFieldSet("environments.0.environment_id", attrs["environments.0.environment_id"]); err != nil {
				return err
			}
			if err := util.CheckFieldSet("environments.0.public_url", attrs["environments.0.public_url"]); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("environments.0.tags.0", attrs["environments.0.tags.0"], "a"); err != nil {
				return err
			}
		}
		return nil
	}
}

func testAccSourceEnvironmentsConfig(domain string, projectName string, name1 string, name2 string) string {
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

resource "buddy_environment" "b" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   name = "%s"
   identifier = "%s"
   public_url = "https://b.com"
   tags = ["a"]
}

data "buddy_environments" "all" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   depends_on = [buddy_environment.a, buddy_environment.b]
}

data "buddy_environments" "name" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.proj.name}"
   name_regex = "^aaaa"
   depends_on = [buddy_environment.a, buddy_environment.b]
}
`, domain, projectName, name1, name1, name2, name2)
}
