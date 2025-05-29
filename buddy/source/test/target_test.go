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

func TestAccTarget(t *testing.T) {
	domain := util.UniqueString()
	name := util.UniqueString()
	identifier := util.UniqueString()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acc.PreCheck(t) },
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTargetConfig(domain, name, identifier),
				Check:  testAccTargetAttributes("buddy_target.t", name, identifier),
			},
			{
				ResourceName:      "buddy_target.t",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccTargetConfig(domain, name, identifier string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
  domain = "%s"
}

resource "buddy_target" "t" {
  domain     = "${buddy_workspace.foo.domain}"
  name       = "%s"
  identifier = "%s"
  type       = "%s"
  host       = "1.1.1.1"
  port       = "33"
  secure     = true
}
`, domain, name, identifier, buddy.TargetTypeFtp)
}

func testAccTargetAttributes(n, name, identifier string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		secure, _ := strconv.ParseBool(attrs["secure"])
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("identifier", attrs["identifier"], identifier); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("type", attrs["type"], buddy.TargetTypeFtp); err != nil {
			return err
		}
		if err := util.CheckFieldSet("target_id", attrs["target_id"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("html_url", attrs["html_url"]); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("host", attrs["host"], "1.1.1.1"); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("port", attrs["port"], "33"); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("secure", secure, true); err != nil {
			return err
		}
		return nil
	}
}
