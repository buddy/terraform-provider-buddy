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
	"partner_token",
	"google_config",
	"google_project",
	"audience",
}

func TestAccIntegration_amazon_trusted(t *testing.T) {
	var integration buddy.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	scope := buddy.IntegrationScopeAdmin
	newScope := buddy.IntegrationScopeWorkspace
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationAmazonTrusted(domain, name, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeAmazon, buddy.IntegrationAuthTypeTrusted, scope, false, false, ""),
				),
			},
			// update integration
			{
				Config: testAccIntegrationAmazonTrusted(domain, newName, newScope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeAmazon, buddy.IntegrationAuthTypeTrusted, newScope, false, false, ""),
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

func TestAccIntegration_amazon_oidc(t *testing.T) {
	var integration buddy.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	scope := buddy.IntegrationScopeAdmin
	newScope := buddy.IntegrationScopeWorkspace
	audience := util.RandString(10)
	newAudience := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationAmazonOidc(domain, name, scope, audience),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeAmazon, buddy.IntegrationAuthTypeOidc, scope, false, false, ""),
				),
			},
			// update integration
			{
				Config: testAccIntegrationAmazonOidc(domain, newName, newScope, newAudience),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeAmazon, buddy.IntegrationAuthTypeOidc, newScope, false, false, ""),
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

func TestAccIntegration_amazon_default(t *testing.T) {
	var integration buddy.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	scope := buddy.IntegrationScopeAdmin
	newScope := buddy.IntegrationScopeWorkspace
	identifier := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationAmazonDefault(domain, name, scope, identifier),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeAmazon, buddy.IntegrationAuthTypeDefault, scope, false, false, identifier),
				),
			},
			// update integration
			{
				Config: testAccIntegrationAmazonDefault(domain, newName, newScope, identifier),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeAmazon, buddy.IntegrationAuthTypeDefault, newScope, false, false, identifier),
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

func TestAccIntegration_amazon_recreate(t *testing.T) {
	var integration buddy.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	scope := buddy.IntegrationScopeAdmin
	newScope := buddy.IntegrationScopeWorkspace
	identifier := util.RandString(10)
	newIdentifier := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationAmazonDefault(domain, name, scope, identifier),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeAmazon, buddy.IntegrationAuthTypeDefault, scope, false, false, identifier),
				),
			},
			// update integration
			{
				Config: testAccIntegrationAmazonDefault(domain, newName, newScope, newIdentifier),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeAmazon, buddy.IntegrationAuthTypeDefault, newScope, false, false, newIdentifier),
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

func TestAccIntegration_amazon(t *testing.T) {
	var integration buddy.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	scope := buddy.IntegrationScopeAdmin
	newScope := buddy.IntegrationScopeWorkspace
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationAmazon(domain, name, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeAmazon, buddy.IntegrationAuthTypeDefault, scope, false, false, ""),
				),
			},
			// update integration
			{
				Config: testAccIntegrationAmazon(domain, newName, newScope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeAmazon, buddy.IntegrationAuthTypeDefault, newScope, false, false, ""),
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
	var integration buddy.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	groupNameA := util.RandString(10)
	groupNameB := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationDigitalOcean(domain, name, groupNameA, groupNameB, groupNameA),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeDigitalOcean, "", buddy.IntegrationScopeGroup, false, true, ""),
				),
			},
			// update integration
			{
				Config: testAccIntegrationDigitalOcean(domain, newName, groupNameA, groupNameB, groupNameB),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeDigitalOcean, "", buddy.IntegrationScopeGroup, false, true, ""),
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
	var integration buddy.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	projectNameA := util.RandString(10)
	projectNameB := util.RandString(10)
	scope := buddy.IntegrationScopeProject
	newScope := buddy.IntegrationScopePrivateInProject
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationShopify(domain, name, projectNameA, projectNameB, scope, projectNameA),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeShopify, buddy.IntegrationAuthTypeToken, scope, true, false, ""),
				),
			},
			// update integration
			{
				Config: testAccIntegrationShopify(domain, newName, projectNameA, projectNameB, newScope, projectNameB),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeShopify, buddy.IntegrationAuthTypeToken, newScope, true, false, ""),
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

func TestAccIntegration_shopify_partner(t *testing.T) {
	var integration buddy.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	projectNameA := util.RandString(10)
	projectNameB := util.RandString(10)
	scope := buddy.IntegrationScopeProject
	newScope := buddy.IntegrationScopePrivateInProject
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationShopifyPartner(domain, name, projectNameA, projectNameB, scope, projectNameA),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeShopify, buddy.IntegrationAuthTypeTokenAppExtension, scope, true, false, ""),
				),
			},
			// update integration
			{
				Config: testAccIntegrationShopifyPartner(domain, newName, projectNameA, projectNameB, newScope, projectNameB),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeShopify, buddy.IntegrationAuthTypeTokenAppExtension, newScope, true, false, ""),
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

func TestAccIntegration_gitlab(t *testing.T) {
	var integration buddy.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	scope := buddy.IntegrationScopeWorkspace
	newScope := buddy.IntegrationScopeWorkspace
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationGitLab(domain, name, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeGitLab, "", scope, false, false, ""),
				),
			},
			// update integration
			{
				Config: testAccIntegrationGitLab(domain, newName, newScope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeGitLab, "", newScope, false, false, ""),
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

func TestAccIntegration_github(t *testing.T) {
	var integration buddy.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	scope := buddy.IntegrationScopeWorkspace
	newScope := buddy.IntegrationScopeWorkspace
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationGitHub(domain, name, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeGitHub, "", scope, false, false, ""),
				),
			},
			// update integration
			{
				Config: testAccIntegrationGitHub(domain, newName, newScope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeGitHub, "", newScope, false, false, ""),
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
	var integration buddy.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	scope := buddy.IntegrationScopeWorkspace
	newScope := buddy.IntegrationScopeAdmin
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationRackspace(domain, name, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeRackspace, "", scope, false, false, ""),
				),
			},
			// update integration
			{
				Config: testAccIntegrationRackspace(domain, newName, newScope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeRackspace, "", newScope, false, false, ""),
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
	var integration buddy.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	scope := buddy.IntegrationScopeWorkspace
	newScope := buddy.IntegrationScopeAdmin
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationCloudflare(domain, name, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeCloudflare, "", scope, false, false, ""),
				),
			},
			// update integration
			{
				Config: testAccIntegrationCloudflare(domain, newName, newScope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeCloudflare, "", newScope, false, false, ""),
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
	var integration buddy.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	scope := buddy.IntegrationScopeWorkspace
	newScope := buddy.IntegrationScopeAdmin
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationUpcloud(domain, name, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeUpcloud, "", scope, false, false, ""),
				),
			},
			// update integration
			{
				Config: testAccIntegrationUpcloud(domain, newName, newScope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeUpcloud, "", newScope, false, false, ""),
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

func TestAccIntegration_stackHawk(t *testing.T) {
	var integration buddy.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	scope := buddy.IntegrationScopeWorkspace
	newScope := buddy.IntegrationScopeAdmin
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationStackHawk(domain, name, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeStackHawk, "", scope, false, false, ""),
				),
			},
			// update integration
			{
				Config: testAccIntegrationStackHawk(domain, newName, newScope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeStackHawk, "", newScope, false, false, ""),
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

func TestAccIntegration_google_oidc(t *testing.T) {
	var integration buddy.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	scope := buddy.IntegrationScopeWorkspace
	newScope := buddy.IntegrationScopeAdmin
	audience := util.RandString(10)
	newAudience := util.RandString(10)
	googleProject := util.RandString(10)
	newGoogleProject := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationAzureGoogleOidc(domain, name, scope, audience, googleProject),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeGoogleServiceAccount, buddy.IntegrationAuthTypeOidc, scope, false, false, ""),
				),
			},
			// update integration
			{
				Config: testAccIntegrationAzureGoogleOidc(domain, newName, newScope, newAudience, newGoogleProject),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeGoogleServiceAccount, buddy.IntegrationAuthTypeOidc, newScope, false, false, ""),
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

func TestAccIntegration_azurecloud_oidc(t *testing.T) {
	var integration buddy.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	scope := buddy.IntegrationScopeWorkspace
	newScope := buddy.IntegrationScopeAdmin
	audience := util.RandString(10)
	newAudience := util.RandString(10)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationAzureCloudOidc(domain, name, scope, audience),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeAzureCloud, buddy.IntegrationAuthTypeOidc, scope, false, false, ""),
				),
			},
			// update integration
			{
				Config: testAccIntegrationAzureCloudOidc(domain, newName, newScope, newAudience),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeAzureCloud, buddy.IntegrationAuthTypeOidc, newScope, false, false, ""),
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
	var integration buddy.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	newName := util.RandString(10)
	scope := buddy.IntegrationScopeWorkspace
	newScope := buddy.IntegrationScopeAdmin
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationAzureCloud(domain, name, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeAzureCloud, buddy.IntegrationAuthTypeDefault, scope, false, false, ""),
				),
			},
			// update integration
			{
				Config: testAccIntegrationAzureCloud(domain, newName, newScope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeAzureCloud, buddy.IntegrationAuthTypeDefault, newScope, false, false, ""),
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

func testAccIntegrationAttributes(n string, integration *buddy.Integration, name string, typ string, authType string, scope string, testScopeProject bool, testScopeGroup bool, identifier string) resource.TestCheckFunc {
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
		if authType != "" {
			if err := util.CheckFieldEqualAndSet("AuthType", integration.AuthType, authType); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("auth_type", attrs["auth_type"], authType); err != nil {
				return err
			}
		}
		if identifier != "" {
			if err := util.CheckFieldEqualAndSet("Identifier", integration.Identifier, identifier); err != nil {
				return err
			}
		} else {
			if err := util.CheckFieldSet("Identifier", integration.Identifier); err != nil {
				return err
			}
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
		if err := util.CheckFieldSet("html_url", attrs["html_url"]); err != nil {
			return err
		}
		return nil
	}
}

func testAccIntegrationGet(n string, integration *buddy.Integration) resource.TestCheckFunc {
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

func testAccIntegrationAmazonOidc(domain string, name string, scope string, audience string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_integration" "bar" {
   domain = "${buddy_workspace.foo.domain}"
   name = "%s"
   type = "%s"
   scope = "%s"
   auth_type = "OIDC"

   audience = "%s"

   role_assumption {
       arn = "arn1"
   }
}
`, domain, name, buddy.IntegrationTypeAmazon, scope, audience)
}

func testAccIntegrationAmazonTrusted(domain string, name string, scope string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_integration" "bar" {
   domain = "${buddy_workspace.foo.domain}"
   name = "%s"
   type = "%s"
   scope = "%s"

   role_assumption {
       arn = "arn1"
			 external_id = "123"
   }
}
`, domain, name, buddy.IntegrationTypeAmazon, scope)
}

func testAccIntegrationAmazonDefault(domain string, name string, scope string, identifier string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_integration" "bar" {
   domain = "${buddy_workspace.foo.domain}"
   name = "%s"
   type = "%s"
   scope = "%s"
   identifier = "%s"
   access_key = "ABC1234567890"
   secret_key = "ABC1234567890"

   auth_type = "DEFAULT"

   role_assumption {
       arn = "arn1"
   }

   role_assumption {
       arn = "arn2"
       external_id = "3"
       duration = 100
   }
}
`, domain, name, buddy.IntegrationTypeAmazon, scope, identifier)
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
`, domain, name, buddy.IntegrationTypeAmazon, scope)
}

func testAccIntegrationGitHub(domain string, name string, scope string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_integration" "bar" {
   domain = "${buddy_workspace.foo.domain}"
   name = "%s"
   type = "%s"
   scope = "%s"
   token = "ABC1234567890"
}
`, domain, name, buddy.IntegrationTypeGitHub, scope)
}

func testAccIntegrationGitLab(domain string, name string, scope string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_integration" "bar" {
   domain = "${buddy_workspace.foo.domain}"
   name = "%s"
   type = "%s"
   scope = "%s"
   token = "ABC1234567890"
}
`, domain, name, buddy.IntegrationTypeGitLab, scope)
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
`, domain, name, buddy.IntegrationTypeRackspace, scope)
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
`, domain, name, buddy.IntegrationTypeCloudflare, scope)
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
`, domain, name, buddy.IntegrationTypeUpcloud, scope)
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
`, domain, name, buddy.IntegrationTypeAzureCloud, scope)
}

func testAccIntegrationAzureGoogleOidc(domain string, name string, scope string, audience string, googleProject string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_integration" "bar" {
   domain = "${buddy_workspace.foo.domain}"
   name = "%s"
   type = "%s"
   scope = "%s"
	 auth_type = "OIDC"
	 audience = "%s"
   google_config = "{}"
   google_project = "%s"
}
`, domain, name, buddy.IntegrationTypeGoogleServiceAccount, scope, audience, googleProject)
}

func testAccIntegrationAzureCloudOidc(domain string, name string, scope string, audience string) string {
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
   audience = "%s"
	 auth_type = "OIDC"
}
`, domain, name, buddy.IntegrationTypeAzureCloud, scope, audience)
}

func testAccIntegrationStackHawk(domain string, name string, scope string) string {
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
}
`, domain, name, buddy.IntegrationTypeStackHawk, scope)
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
`, domain, groupNameA, groupNameA, groupNameB, groupNameB, name, buddy.IntegrationTypeDigitalOcean, buddy.IntegrationScopeGroup, scopeGroupName)
}

func testAccIntegrationShopify(domain string, name string, projectNameA string, projectNameB string, scope string, scopeProjectName string) string {
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
`, domain, projectNameA, projectNameA, projectNameB, projectNameB, name, buddy.IntegrationTypeShopify, scope, scopeProjectName)
}

func testAccIntegrationShopifyPartner(domain string, name string, projectNameA string, projectNameB string, scope string, scopeProjectName string) string {
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
   token = "ABC"
	partner_token = "EFG"
}
`, domain, projectNameA, projectNameA, projectNameB, projectNameB, name, buddy.IntegrationTypeShopify, scope, scopeProjectName)
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
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}
