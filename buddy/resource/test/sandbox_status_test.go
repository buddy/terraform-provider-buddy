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

type testAccSandboxStatusExpectedAttributes struct {
	Status string
	Wait   bool
}

func testAccSandboxStatusAttributes(n string, sandbox *buddy.Sandbox, want *testAccSandboxStatusExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsWaitt, _ := strconv.ParseBool(attrs["wait_for_status"])
		if want.Wait {
			if err := util.CheckFieldEqualAndSet("Sandbox.Status", sandbox.Status, want.Status); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("sandbox_status", attrs["sandbox_status"], want.Status); err != nil {
				return err
			}
		}
		if err := util.CheckFieldEqualAndSet("status", attrs["status"], want.Status); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("wait", attrsWaitt, want.Wait); err != nil {
			return err
		}
		return nil
	}
}

func TestAccSandboxStatus(t *testing.T) {
	var sandbox buddy.Sandbox
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             acc.DummyCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSandboxStatusConfig(domain, projectName, name, buddy.SandboxStatusRunning, false),
				Check: resource.ComposeTestCheckFunc(
					testAccSandboxGet("buddy_sandbox.bar", &sandbox),
					testAccSandboxStatusAttributes("buddy_sandbox_status.s", &sandbox, &testAccSandboxStatusExpectedAttributes{
						Wait:   false,
						Status: buddy.SandboxStatusRunning,
					}),
				),
			},
			{
				Config: testAccSandboxStatusConfig(domain, projectName, name, buddy.SandboxStatusStopped, true),
				Check: resource.ComposeTestCheckFunc(
					testAccSandboxGet("buddy_sandbox.bar", &sandbox),
					testAccSandboxStatusAttributes("buddy_sandbox_status.s", &sandbox, &testAccSandboxStatusExpectedAttributes{
						Wait:   true,
						Status: buddy.SandboxStatusStopped,
					}),
				),
			},
			{
				Config: testAccSandboxStatusConfig(domain, projectName, name, buddy.SandboxStatusRunning, true),
				Check: resource.ComposeTestCheckFunc(
					testAccSandboxGet("buddy_sandbox.bar", &sandbox),
					testAccSandboxStatusAttributes("buddy_sandbox_status.s", &sandbox, &testAccSandboxStatusExpectedAttributes{
						Wait:   true,
						Status: buddy.SandboxStatusRunning,
					}),
				),
			},
			{
				Config: testAccSandboxStatusConfig(domain, projectName, name, buddy.SandboxStatusStopped, false),
				Check: resource.ComposeTestCheckFunc(
					testAccSandboxGet("buddy_sandbox.bar", &sandbox),
					testAccSandboxStatusAttributes("buddy_sandbox_status.s", &sandbox, &testAccSandboxStatusExpectedAttributes{
						Wait:   false,
						Status: buddy.SandboxStatusStopped,
					}),
				),
			},
		},
	})
}

func testAccSandboxStatusConfig(domain string, projectName string, name string, status string, wait bool) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_project" "proj" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_sandbox" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    project_name = "${buddy_project.proj.name}"
    name = "%s"
}

resource "buddy_sandbox_status" "s" {
    domain = "${buddy_workspace.foo.domain}"
    sandbox_id = "${buddy_sandbox.bar.sandbox_id}"
    status = "%s"
    wait_for_status = "%t"
}
`, domain, projectName, name, status, wait)
}
