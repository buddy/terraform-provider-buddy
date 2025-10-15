package test

import (
	"fmt"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"terraform-provider-buddy/buddy/acc"
	"terraform-provider-buddy/buddy/util"
	"testing"
)

func TestAccSandbox_wait(t *testing.T) {
	var sandbox buddy.Sandbox
	var project buddy.Project
	domain := util.UniqueString()
	projectName := util.UniqueString()
	name := util.RandString(10)
	installCommands := "sleep 10"
	runCommand := "while :; do foo; sleep 2; done"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccSandboxCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSandboxConfigWait(domain, projectName, name, installCommands, runCommand),
				Check: resource.ComposeTestCheckFunc(
					testAccSandboxGet("buddy_sandbox.bar", &sandbox),
					testAccProjectGet("buddy_project.proj", &project),
					testAccSandboxAttributes("buddy_sandbox.bar", &sandbox, &testAccSandboxExpectedAttributes{
						Name:              name,
						InstallCommands:   installCommands,
						RunCommand:        runCommand,
						Project:           &project,
						WaitForApp:        true,
						WaitForConfigured: true,
						WaitForRunning:    true,
					}),
				),
			},
		},
	})
}

func TestAccSandbox_main(t *testing.T) {
	var sandbox buddy.Sandbox
	var project buddy.Project
	domain := util.UniqueString()
	projectName := util.UniqueString()
	identifier := util.UniqueString()
	newIdentifier := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	installCommands := "pwd"
	runCommand := "while :; do foo; sleep 2; done"
	appDir := "/"
	appType := buddy.SandboxAppTypeCmd
	os := buddy.SandboxOsUbuntu2204
	newOs := buddy.SandboxOsUbuntu2404
	resources := buddy.SandboxResource2X4
	tag := util.RandString(10)
	newTag := util.RandString(10)
	tcpName := util.UniqueString()
	tcpEndpoint := "22"
	newTcpEndppoint := "123"
	tlsName := util.UniqueString()
	tlsEndpoint := "222"
	tlsTerminateAt := buddy.SandboxEndpointTlsTerminateAtRegion
	httpName := util.UniqueString()
	httpEndpoint := "444"
	httpCompresion := true
	httpXHeader := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccSandboxCheckDestroy,
		Steps: []resource.TestStep{
			{
				// create
				Config: testAccSandboxConfigOneEndpoint(domain, projectName, identifier, name, installCommands, runCommand, appDir, appType, os, resources, tag, tcpName, tcpEndpoint),
				Check: resource.ComposeTestCheckFunc(
					testAccSandboxGet("buddy_sandbox.bar", &sandbox),
					testAccProjectGet("buddy_project.proj", &project),
					testAccSandboxAttributes("buddy_sandbox.bar", &sandbox, &testAccSandboxExpectedAttributes{
						Identifier:      identifier,
						Name:            name,
						InstallCommands: installCommands,
						RunCommand:      runCommand,
						AppDir:          appDir,
						AppType:         appType,
						Os:              os,
						Resources:       resources,
						Tag:             tag,
						TcpName:         tcpName,
						TcpEndpoint:     tcpEndpoint,
						Project:         &project,
					}),
				),
			},
			{
				// update
				Config: testAccSandboxConfig(domain, projectName, newIdentifier, newName, installCommands, runCommand, appDir, appType, os, resources, newTag, tcpName, newTcpEndppoint, tlsName, tlsEndpoint, tlsTerminateAt, httpName, httpEndpoint, httpCompresion, httpXHeader),
				Check: resource.ComposeTestCheckFunc(
					testAccSandboxGet("buddy_sandbox.bar", &sandbox),
					testAccProjectGet("buddy_project.proj", &project),
					testAccSandboxAttributes("buddy_sandbox.bar", &sandbox, &testAccSandboxExpectedAttributes{
						Identifier:      newIdentifier,
						Name:            newName,
						InstallCommands: installCommands,
						RunCommand:      runCommand,
						AppDir:          appDir,
						AppType:         appType,
						Os:              os,
						Resources:       resources,
						Tag:             newTag,
						TcpName:         tcpName,
						TcpEndpoint:     newTcpEndppoint,
						TlsName:         tlsName,
						TlsEndpoint:     tlsEndpoint,
						TlsTerminateAt:  tlsTerminateAt,
						HttpName:        httpName,
						HttpEndpoint:    httpEndpoint,
						HttpCompresion:  httpCompresion,
						HttpXHeader:     httpXHeader,
						Project:         &project,
					}),
				),
			},
			{
				Config: testAccSandboxConfigOneEndpoint(domain, projectName, identifier, name, installCommands, runCommand, appDir, appType, newOs, resources, tag, tcpName, tcpEndpoint),
				Check: resource.ComposeTestCheckFunc(
					testAccSandboxGet("buddy_sandbox.bar", &sandbox),
					testAccProjectGet("buddy_project.proj", &project),
					testAccSandboxAttributes("buddy_sandbox.bar", &sandbox, &testAccSandboxExpectedAttributes{
						Identifier:      identifier,
						Name:            name,
						InstallCommands: installCommands,
						RunCommand:      runCommand,
						AppDir:          appDir,
						AppType:         appType,
						Os:              newOs,
						Resources:       resources,
						Tag:             tag,
						TcpName:         tcpName,
						TcpEndpoint:     tcpEndpoint,
						Project:         &project,
					}),
				),
			},
			{
				// import
				ResourceName:      "buddy_sandbox.bar",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"wait_for_app_timeout",
					"wait_for_configured_timeout",
					"wait_for_running_timeout",
				},
			},
		},
	})
}

type testAccSandboxExpectedAttributes struct {
	Name              string
	Identifier        string
	InstallCommands   string
	RunCommand        string
	AppDir            string
	AppType           string
	Os                string
	Resources         string
	Tag               string
	TcpName           string
	TcpEndpoint       string
	TlsName           string
	TlsEndpoint       string
	TlsTerminateAt    string
	HttpName          string
	HttpEndpoint      string
	HttpCompresion    bool
	HttpXHeader       string
	WaitForRunning    bool
	WaitForConfigured bool
	WaitForApp        bool
	Project           *buddy.Project
}

func testAccSandboxAttributes(n string, sandbox *buddy.Sandbox, want *testAccSandboxExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		if err := util.CheckFieldEqualAndSet("Name", sandbox.Name, want.Name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], want.Name); err != nil {
			return err
		}
		if want.Identifier != "" {
			if err := util.CheckFieldEqualAndSet("Identifier", sandbox.Identifier, want.Identifier); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("identifier", attrs["identifier"], want.Identifier); err != nil {
				return err
			}
		} else {
			if err := util.CheckFieldSet("Identifier", sandbox.Identifier); err != nil {
				return err
			}
			if err := util.CheckFieldSet("identifier", attrs["identifier"]); err != nil {
				return err
			}
		}
		if want.InstallCommands != "" {
			if err := util.CheckFieldEqualAndSet("InstallCommands", sandbox.InstallCommands, want.InstallCommands); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("installCommands", attrs["install_commands"], want.InstallCommands); err != nil {
				return err
			}
		}
		if want.RunCommand != "" {
			if err := util.CheckFieldEqualAndSet("RunCommand", sandbox.RunCommand, want.RunCommand); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("runCommand", attrs["run_command"], want.RunCommand); err != nil {
				return err
			}
		}
		if want.AppDir != "" {
			if err := util.CheckFieldEqualAndSet("AppDir", sandbox.AppDir, want.AppDir); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("appDir", attrs["app_dir"], want.AppDir); err != nil {
				return err
			}
		}
		if want.AppType != "" {
			if err := util.CheckFieldEqualAndSet("AppType", sandbox.AppType, want.AppType); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("appType", attrs["app_type"], want.AppType); err != nil {
				return err
			}
		}
		if want.Os != "" {
			if err := util.CheckFieldEqualAndSet("Os", sandbox.Os, want.Os); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("os", attrs["os"], want.Os); err != nil {
				return err
			}
		} else {
			if err := util.CheckFieldEqualAndSet("Os", sandbox.Os, buddy.SandboxOsUbuntu2404); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("os", attrs["os"], buddy.SandboxOsUbuntu2404); err != nil {
				return err
			}
		}
		if want.Resources != "" {
			if err := util.CheckFieldEqualAndSet("Resources", sandbox.Resources, want.Resources); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("resources", attrs["resources"], want.Resources); err != nil {
				return err
			}
		} else {
			if err := util.CheckFieldSet("Resources", sandbox.Resources); err != nil {
				return err
			}
			if err := util.CheckFieldSet("resources", attrs["resources"]); err != nil {
				return err
			}
		}
		if want.Tag != "" {
			if err := util.CheckFieldEqualAndSet("Tags[0]", sandbox.Tags[0], want.Tag); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("tags.0", attrs["tags.0"], want.Tag); err != nil {
				return err
			}
		}
		if err := util.CheckFieldEqualAndSet("ProjectName", sandbox.Project.Name, want.Project.Name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("project_name", attrs["project_name"], want.Project.Name); err != nil {
			return err
		}
		if want.TcpName != "" {
			endpoint := findEndpointByName(&sandbox.Endpoints, want.TcpName)
			if endpoint == nil {
				return fmt.Errorf("tcp endpoint is null")
			}
			if endpoint.Endpoint == nil {
				return fmt.Errorf("tcp endpoint endpoint is null")
			}
			if endpoint.Type == nil {
				return fmt.Errorf("tcp endpoint type is null")
			}
			if err := util.CheckFieldEqualAndSet("Endpoint.Endpoint", *endpoint.Endpoint, want.TcpEndpoint); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("Endpoint.Type", *endpoint.Type, buddy.SandboxEndpointTypeTcp); err != nil {
				return err
			}
		} else {
			endpoint := findEndpointByType(&sandbox.Endpoints, buddy.SandboxEndpointTypeTcp)
			if endpoint != nil {
				return fmt.Errorf("tcp endpoint is not null")
			}
		}
		if want.TlsName != "" {
			endpoint := findEndpointByName(&sandbox.Endpoints, want.TlsName)
			if endpoint == nil {
				return fmt.Errorf("tls endpoint is null")
			}
			if endpoint.Endpoint == nil {
				return fmt.Errorf("tls endpoint endpoint is null")
			}
			if endpoint.Type == nil {
				return fmt.Errorf("tls endpoint type is null")
			}
			if endpoint.Tls == nil {
				return fmt.Errorf("tls endpoint tls is null")
			}
			if endpoint.Tls.TerminateAt == nil {
				return fmt.Errorf("tls endpoint tls terminate at is null")
			}
			if err := util.CheckFieldEqualAndSet("Endpoint.Endpoint", *endpoint.Endpoint, want.TlsEndpoint); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("Endpoint.Type", *endpoint.Type, buddy.SandboxEndpointTypeTls); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("Endpoint.Tls.TerminateAt", *endpoint.Tls.TerminateAt, want.TlsTerminateAt); err != nil {
				return err
			}
		} else {
			endpoint := findEndpointByType(&sandbox.Endpoints, buddy.SandboxEndpointTypeTls)
			if endpoint != nil {
				return fmt.Errorf("tls endpoint is not null")
			}
		}
		if want.HttpName != "" {
			endpoint := findEndpointByName(&sandbox.Endpoints, want.HttpName)
			if endpoint == nil {
				return fmt.Errorf("http endpoint is null")
			}
			if endpoint.Endpoint == nil {
				return fmt.Errorf("http endpoint endpoint is null")
			}
			if endpoint.Type == nil {
				return fmt.Errorf("http endpoint type is null")
			}
			if endpoint.Http == nil {
				return fmt.Errorf("http endpoint http is null")
			}
			if endpoint.Http.Compression == nil {
				return fmt.Errorf("http endpoint http compression is null")
			}
			if endpoint.Http.RequestHeaders == nil {
				return fmt.Errorf("http endpoint http request headers is null")
			}
			if endpoint.Http.ResponseHeaders == nil {
				return fmt.Errorf("http endpoint http response headers is null")
			}
			if err := util.CheckFieldEqualAndSet("Endpoint.Endpoint", *endpoint.Endpoint, want.HttpEndpoint); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("Endpoint.Type", *endpoint.Type, buddy.SandboxEndpointTypeHttp); err != nil {
				return err
			}
			if err := util.CheckBoolFieldEqual("Endpoint.Tls.TerminateAt", *endpoint.Http.Compression, want.HttpCompresion); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("Endpoint.Http.RequestHeaders", (*endpoint.Http.RequestHeaders)["X-Custom"], want.HttpXHeader); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("Endpoint.Http.ResponseHeaders", (*endpoint.Http.ResponseHeaders)["X-Custom"], want.HttpXHeader); err != nil {
				return err
			}
		} else {
			endpoint := findEndpointByType(&sandbox.Endpoints, buddy.SandboxEndpointTypeHttp)
			if endpoint != nil {
				return fmt.Errorf("http endpoint is not null")
			}
		}
		if want.WaitForRunning {
			if err := util.CheckFieldEqualAndSet("Status", sandbox.Status, buddy.SandboxStatusRunning); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("status", attrs["status"], buddy.SandboxStatusRunning); err != nil {
				return err
			}
		}
		if want.WaitForConfigured {
			if err := util.CheckFieldEqualAndSet("SetupStatus", sandbox.SetupStatus, buddy.SandboxSetupStatusSuccess); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("setup_status", attrs["setup_status"], buddy.SandboxSetupStatusSuccess); err != nil {
				return err
			}
		}
		if want.WaitForApp {
			if err := util.CheckFieldEqualAndSet("AppStatus", sandbox.AppStatus, buddy.SandboxAppStatusRunning); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("app_status", attrs["app_status"], buddy.SandboxAppStatusRunning); err != nil {
				return err
			}
		}
		return nil
	}
}

func findEndpointByType(endpoints *[]buddy.SandboxEndpoint, typ string) *buddy.SandboxEndpoint {
	if endpoints != nil {
		for _, v := range *endpoints {
			if v.Type != nil && *v.Type == typ {
				return &v
			}
		}
	}
	return nil
}

func findEndpointByName(endpoints *[]buddy.SandboxEndpoint, name string) *buddy.SandboxEndpoint {
	if endpoints != nil {
		for _, v := range *endpoints {
			if v.Name != nil && *v.Name == name {
				return &v
			}
		}
	}
	return nil
}

func testAccSandboxGet(n string, sandbox *buddy.Sandbox) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		domain, sandboxId, err := util.DecomposeDoubleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		sb, _, err := acc.ApiClient.SandboxService.Get(domain, sandboxId)
		if err != nil {
			return err
		}
		*sandbox = *sb
		return nil
	}
}

func testAccSandboxConfigWait(domain string, projectName string, name string, installCommands string, runCommand string) string {
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
    install_commands = "%s"
    run_command = "%s"
    wait_for_running = true
    wait_for_configured = true
    wait_for_app = true
}
`, domain, projectName, name, installCommands, runCommand)
}

func testAccSandboxConfigOneEndpoint(domain string, projectName string, identifier string, name string, installCommands string, runCommand string, appDir string, appType string, os string, resources string, tag string, tcpName string, tcpEndpoint string) string {
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
    identifier = "%s"
    name = "%s"
    install_commands = "%s"
    run_command = "%s"
    app_dir = "%s"
    app_type = "%s"
    os = "%s"
    resources = "%s" 
    tags = ["%s"]
    endpoints = {
      "%s" = {
        endpoint = "%s"
        type = "TCP"
      }
    }
}
`, domain, projectName, identifier, name, installCommands, runCommand, appDir, appType, os, resources, tag, tcpName, tcpEndpoint)
}

func testAccSandboxConfig(domain string, projectName string, identifier string, name string, installCommands string, runCommand string, appDir string, appType string, os string, resources string, tag string, tcpName string, tcpEndpoint string, tlsName string, tlsEndpoint string, tlsTerminateAt string, httpName string, httpEndpoint string, httpCompresion bool, httpXHeader string) string {
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
    identifier = "%s"
    name = "%s"
    install_commands = "%s"
    run_command = "%s"
    app_dir = "%s"
    app_type = "%s"
    os = "%s"
    resources = "%s" 
    tags = ["%s"]
    endpoints = {
      "%s" = {
        endpoint = "%s"
        type = "TCP"
      },
      "%s" = {
        endpoint = "%s"
        type = "TLS"
        tls = {
          terminate_at = "%s"
        }
      }, 
      "%s" = {
        endpoint = "%s"
        type = "HTTP"
        http = {
          compression = "%t"
          request_headers = {
            X-Custom = "%s"
          }
          response_headers = {
            X-Custom = "%s"
          }
        } 
      }
    }
}
`, domain, projectName, identifier, name, installCommands, runCommand, appDir, appType, os, resources, tag, tcpName, tcpEndpoint, tlsName, tlsEndpoint, tlsTerminateAt, httpName, httpEndpoint, httpCompresion, httpXHeader, httpXHeader)
}

func testAccSandboxCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_sandbox" {
			continue
		}
		domain, sandboxId, err := util.DecomposeDoubleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		sandbox, resp, err := acc.ApiClient.SandboxService.Get(domain, sandboxId)
		if err == nil && sandbox != nil {
			return util.ErrorResourceExists()
		}
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}
