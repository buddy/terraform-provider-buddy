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

var ignoreImportVerify = []string{
	"username",
	"shop",
	"token",
	"access_key",
	"secret_key",
	"app_id",
	"tenant_id",
	"password",
	"api_key",
	"email",
	"role_assumption",
}

func TestAccIntegration_amazon(t *testing.T) {
	var integration api.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	scope := api.IntegrationScopeAdmin
	newScope := api.IntegrationScopeWorkspace
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationAmazon(domain, name, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, api.IntegrationTypeAmazon, scope, false, false),
				),
			},
			// update integration
			{
				Config: testAccIntegrationAmazon(domain, newName, newScope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, api.IntegrationTypeAmazon, newScope, false, false),
				),
			},
			// import integration
			{
				ResourceName:            "buddy_integration.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoreImportVerify,
			},
		},
	})
}

func TestAccIntegration_digitalocean(t *testing.T) {
	var integration api.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	groupNameA := util.RandString(10)
	groupNameB := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationDigitalOcean(domain, name, groupNameA, groupNameB, groupNameA),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, api.IntegrationTypeDigitalOcean, api.IntegrationScopeGroup, false, true),
				),
			},
			// update integration
			{
				Config: testAccIntegrationDigitalOcean(domain, newName, groupNameA, groupNameB, groupNameB),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, api.IntegrationTypeDigitalOcean, api.IntegrationScopeGroup, false, true),
				),
			},
			// import integration
			{
				ResourceName:            "buddy_integration.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoreImportVerify,
			},
		},
	})
}

func TestAccIntegration_shopify(t *testing.T) {
	var integration api.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	projectNameA := util.RandString(10)
	projectNameB := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationShopify(domain, name, projectNameA, projectNameB, projectNameA),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, api.IntegrationTypeShopify, api.IntegrationScopeProject, true, false),
				),
			},
			// update integration
			{
				Config: testAccIntegrationShopify(domain, newName, projectNameA, projectNameB, projectNameB),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, api.IntegrationTypeShopify, api.IntegrationScopeProject, true, false),
				),
			},
			// import integration
			{
				ResourceName:            "buddy_integration.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoreImportVerify,
			},
		},
	})
}

func TestAccIntegration_rackspace(t *testing.T) {
	var integration api.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	scope := api.IntegrationScopeWorkspace
	newScope := api.IntegrationScopeAdmin
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationRackspace(domain, name, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, api.IntegrationTypeRackspace, scope, false, false),
				),
			},
			// update integration
			{
				Config: testAccIntegrationRackspace(domain, newName, newScope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, api.IntegrationTypeRackspace, newScope, false, false),
				),
			},
			// import integration
			{
				ResourceName:            "buddy_integration.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoreImportVerify,
			},
		},
	})
}

func TestAccIntegration_cloudflare(t *testing.T) {
	var integration api.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	scope := api.IntegrationScopeWorkspace
	newScope := api.IntegrationScopeAdmin
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationCloudflare(domain, name, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, api.IntegrationTypeCloudflare, scope, false, false),
				),
			},
			// update integration
			{
				Config: testAccIntegrationCloudflare(domain, newName, newScope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, api.IntegrationTypeCloudflare, newScope, false, false),
				),
			},
			// import integration
			{
				ResourceName:            "buddy_integration.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoreImportVerify,
			},
		},
	})
}

func TestAccIntegration_upcloud(t *testing.T) {
	var integration api.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	scope := api.IntegrationScopeWorkspace
	newScope := api.IntegrationScopeAdmin
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationUpcloud(domain, name, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, api.IntegrationTypeUpcloud, scope, false, false),
				),
			},
			// update integration
			{
				Config: testAccIntegrationUpcloud(domain, newName, newScope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, api.IntegrationTypeUpcloud, newScope, false, false),
				),
			},
			// import integration
			{
				ResourceName:            "buddy_integration.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoreImportVerify,
			},
		},
	})
}

func TestAccIntegration_azurecloud(t *testing.T) {
	var integration api.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	scope := api.IntegrationScopeWorkspace
	newScope := api.IntegrationScopeAdmin
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProviderFactories: acc.ProviderFactories,
		CheckDestroy:      testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationAzureCloud(domain, name, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, api.IntegrationTypeAzureCloud, scope, false, false),
				),
			},
			// update integration
			{
				Config: testAccIntegrationAzureCloud(domain, newName, newScope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, api.IntegrationTypeAzureCloud, newScope, false, false),
				),
			},
			// import integration
			{
				ResourceName:            "buddy_integration.bar",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: ignoreImportVerify,
			},
		},
	})
}

func testAccIntegrationAttributes(n string, integration *api.Integration, name string, typ string, scope string, testScopeProject bool, testScopeGroup bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		if err := util.CheckFieldEqualAndSet("Name", integration.Name, name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("Type", integration.Type, typ); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("Scope", integration.Scope, scope); err != nil {
			return err
		}
		if err := util.CheckFieldSet("ProjectName", integration.ProjectName); testScopeProject && err != nil {
			return err
		}
		if err := util.CheckIntFieldSet("GroupId", integration.GroupId); testScopeGroup && err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("name", attrs["name"], name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("type", attrs["type"], typ); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("scope", attrs["scope"], scope); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("project_name", attrs["project_name"], integration.ProjectName); testScopeProject && err != nil {
			return err
		}
		attrsGroupId, _ := strconv.Atoi(attrs["group_id"])
		if err := util.CheckIntFieldEqualAndSet("group_id", attrsGroupId, integration.GroupId); testScopeGroup && err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("integration_id", attrs["integration_id"], integration.HashId); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("html_url", attrs["html_url"], integration.HtmlUrl); err != nil {
			return err
		}
		return nil
	}
}

func testAccIntegrationGet(n string, integration *api.Integration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		domain, hashId, err := util.DecomposeDoubleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		i, _, err := acc.ApiClient.IntegrationService.Get(domain, hashId)
		if err != nil {
			return err
		}
		*integration = *i
		return nil
	}
}

func testAccIntegrationAmazon(domain string, name string, scope string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_integration" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    name = "%s"
    type = "%s"
    scope = "%s"
    access_key = "ABC1234567890"
    secret_key = "ABC1234567890"

    role_assumption {
        arn = "arn1"
    }

	role_assumption {
        arn = "arn2"
        external_id = "3"
        duration = 100
    }
}
`, domain, name, api.IntegrationTypeAmazon, scope)
}

func testAccIntegrationRackspace(domain string, name string, scope string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_integration" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    name = "%s"
    type = "%s"
    scope = "%s"
    username = "ABC1234567890"
    token = "ABC1234567890"
}
`, domain, name, api.IntegrationTypeRackspace, scope)
}

func testAccIntegrationCloudflare(domain string, name string, scope string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_integration" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    name = "%s"
    type = "%s"
    scope = "%s"
    api_key = "ABC1234567890"
    token = "ABC1234567890"
    email = "test@test.pl"
}
`, domain, name, api.IntegrationTypeCloudflare, scope)
}

func testAccIntegrationUpcloud(domain string, name string, scope string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_integration" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    name = "%s"
    type = "%s"
    scope = "%s"
    username = "ABC1234567890"
    password = "ABC1234567890"
}
`, domain, name, api.IntegrationTypeUpcloud, scope)
}

func testAccIntegrationAzureCloud(domain string, name string, scope string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_integration" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    name = "%s"
    type = "%s"
    scope = "%s"
    app_id = "ABC1234567890"
    tenant_id = "test@test.pl"
    password = "ABC1234567890"
}
`, domain, name, api.IntegrationTypeAzureCloud, scope)
}

func testAccIntegrationDigitalOcean(domain string, name string, groupNameA string, groupNameB string, scopeGroupName string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_group" "%s" {
    domain = "${buddy_workspace.foo.domain}"
    name = "%s"
}

resource "buddy_group" "%s" {
    domain = "${buddy_workspace.foo.domain}"
    name = "%s"
}

resource "buddy_integration" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    name = "%s"
    type = "%s"
    scope = "%s"
    group_id = "${buddy_group.%s.group_id}"
    token = "ABC"
}
`, domain, groupNameA, groupNameA, groupNameB, groupNameB, name, api.IntegrationTypeDigitalOcean, api.IntegrationScopeGroup, scopeGroupName)
}

func testAccIntegrationShopify(domain string, name string, projectNameA string, projectNameB string, scopeProjectName string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
    domain = "%s"
}

resource "buddy_project" "%s" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_project" "%s" {
    domain = "${buddy_workspace.foo.domain}"
    display_name = "%s"
}

resource "buddy_integration" "bar" {
    domain = "${buddy_workspace.foo.domain}"
    name = "%s"
    type = "%s"
    scope = "%s"
    project_name = "${buddy_project.%s.name}"
    shop = "ABC"
    token = "ABC"
}
`, domain, projectNameA, projectNameA, projectNameB, projectNameB, name, api.IntegrationTypeShopify, api.IntegrationScopeProject, scopeProjectName)
}

func testAccIntegrationCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_integration" {
			continue
		}
		domain, hashId, err := util.DecomposeDoubleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		integration, resp, err := acc.ApiClient.IntegrationService.Get(domain, hashId)
		if err == nil && integration != nil {
			return util.ErrorResourceExists()
		}
		if resp.StatusCode != 404 {
			return err
		}
	}
	return nil
}
