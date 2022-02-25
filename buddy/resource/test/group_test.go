package test

import (
	"buddy-terraform/buddy/acc"
	"buddy-terraform/buddy/api"
	"buddy-terraform/buddy/util"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"strconv"
	"testing"
)

func TestAccGroup(t *testing.T) {
	var group api.Group
	domain := util.UniqueString()
	name := util.RandString(5)
	newName := util.RandString(5)
	newDescription := util.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acc.PreCheck(t) },
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccGroupCheckDestroy,
		Steps: []resource.TestStep{
			// create group
			{
				Config: testAccGroupConfig(domain, name),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupGet("buddy_group.bar", &group),
					testAccGroupAttributes("buddy_group.bar", &group, name, ""),
				),
			},
			// update group
			{
				Config: testAccGroupUpdateConfig(domain, newName, newDescription),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupGet("buddy_group.bar", &group),
					testAccGroupAttributes("buddy_group.bar", &group, newName, newDescription),
				),
			},
			// null desc
			{
				Config: testAccGroupConfig(domain, newName),
				Check: resource.ComposeTestCheckFunc(
					testAccGroupGet("buddy_group.bar", &group),
					testAccGroupAttributes("buddy_group.bar", &group, newName, ""),
				),
			},
			// import group
			{
				ResourceName:      "buddy_group.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccGroupAttributes(n string, group *api.Group, name string, description string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		if err := util.CheckFieldEqualAndSet("Name", group.Name, name); err != nil {
			return err
		}
		if err := util.CheckFieldEqual("Description", group.Description, description); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("group_id", attrs["group_id"], strconv.Itoa(group.Id)); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("html_url", attrs["html_url"], group.HtmlUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqual("description", attrs["description"], group.Description); err != nil {
			return err
		}
		return nil
	}
}

func testAccGroupGet(n string, group *api.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		domain, gid, err := util.DecomposeDoubleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		groupId, err := strconv.Atoi(gid)
		if err != nil {
			return err
		}
		g, _, err := acc.ApiClient.GroupService.Get(domain, groupId)
		if err != nil {
			return err
		}
		*group = *g
		return nil
	}
}

func testAccGroupUpdateConfig(domain string, name string, description string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_group" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    name = "%s"
    description = "%s"
}
`, domain, name, description)
}

func testAccGroupConfig(domain string, name string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_group" "bar" {
   domain = "${buddy_workspace.foo.domain}"
   name = "%s"
}
`, domain, name)
}

func testAccGroupCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_group" {
			continue
		}
		domain, gid, err := util.DecomposeDoubleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		groupId, err := strconv.Atoi(gid)
		if err != nil {
			return err
		}
		group, resp, err := acc.ApiClient.GroupService.Get(domain, groupId)
		if err == nil && group != nil {
			return util.ErrorResourceExists()
		}
		if resp.StatusCode != 404 {
			return err
		}
	}
	return nil
}
