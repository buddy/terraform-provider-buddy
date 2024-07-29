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
	"permissions",
	"audience",
}

func TestAccIntegration_amazon_trusted(t *testing.T) {
	var integration buddy.Integration
	domain := util.UniqueString()
	name := util.RandString(10)
	others := buddy.IntegrationPermissionManage
	admins := buddy.IntegrationPermissionManage
	newName := util.RandString(10)
	scope := buddy.IntegrationScopeWorkspace
	newOthers := buddy.IntegrationPermissionUseOnly
	perms := buddy.IntegrationPermissions{
		Others: others,
		Admins: admins,
	}
	newPerms := buddy.IntegrationPermissions{
		Others: newOthers,
		Admins: admins,
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationAmazonTrusted(domain, name, scope, others, admins),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeAmazon, buddy.IntegrationAuthTypeTrusted, scope, false, "", &perms, true, false),
				),
			},
			// update integration
			{
				Config: testAccIntegrationAmazonTrusted(domain, newName, scope, newOthers, admins),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeAmazon, buddy.IntegrationAuthTypeTrusted, scope, false, "", &newPerms, true, false),
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
	scope := buddy.IntegrationScopeWorkspace
	audience := util.RandString(10)
	newAudience := util.RandString(10)
	others := buddy.IntegrationPermissionUseOnly
	admins := buddy.IntegrationPermissionManage
	email := util.RandEmail()
	userAccessLevel := buddy.IntegrationPermissionDenied
	groupName := util.RandString(10)
	groupAccessLevel := buddy.IntegrationPermissionUseOnly
	userAccess := buddy.IntegrationResourcePermission{
		AccessLevel: userAccessLevel,
	}
	groupAccess := buddy.IntegrationResourcePermission{
		AccessLevel: groupAccessLevel,
	}
	perms := buddy.IntegrationPermissions{
		Others: others,
		Admins: admins,
		Users:  []*buddy.IntegrationResourcePermission{&userAccess},
		Groups: []*buddy.IntegrationResourcePermission{&groupAccess},
	}
	newOthers := buddy.IntegrationPermissionManage
	newUserAccessLevel := buddy.IntegrationPermissionUseOnly
	newUserAccess := buddy.IntegrationResourcePermission{
		AccessLevel: newUserAccessLevel,
	}
	newPerms := buddy.IntegrationPermissions{
		Others: newOthers,
		Admins: admins,
		Users:  []*buddy.IntegrationResourcePermission{&newUserAccess},
		Groups: []*buddy.IntegrationResourcePermission{&groupAccess},
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationAmazonOidc(domain, name, scope, audience, others, admins, email, userAccessLevel, groupName, groupAccessLevel),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeAmazon, buddy.IntegrationAuthTypeOidc, scope, false, "", &perms, true, false),
				),
			},
			// update integration
			{
				Config: testAccIntegrationAmazonOidc(domain, newName, scope, newAudience, newOthers, admins, email, newUserAccessLevel, groupName, groupAccessLevel),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeAmazon, buddy.IntegrationAuthTypeOidc, scope, false, "", &newPerms, true, false),
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
	scope := buddy.IntegrationScopeWorkspace
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
				Config: testAccIntegrationAmazonDefault(domain, name, scope, identifier, true, "a"),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeAmazon, buddy.IntegrationAuthTypeDefault, scope, false, identifier, nil, true, true),
				),
			},
			// update integration
			{
				Config: testAccIntegrationAmazonDefault(domain, newName, scope, identifier, false, "b"),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeAmazon, buddy.IntegrationAuthTypeDefault, scope, false, identifier, nil, false, true),
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
	scope := buddy.IntegrationScopeWorkspace
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
				Config: testAccIntegrationAmazonDefault(domain, name, scope, identifier, false, "b"),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeAmazon, buddy.IntegrationAuthTypeDefault, scope, false, identifier, nil, false, true),
				),
			},
			// update integration
			{
				Config: testAccIntegrationAmazonDefault(domain, newName, scope, newIdentifier, true, "a"),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeAmazon, buddy.IntegrationAuthTypeDefault, scope, false, newIdentifier, nil, true, true),
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
	scope := buddy.IntegrationScopeWorkspace
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
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeAmazon, buddy.IntegrationAuthTypeDefault, scope, false, "", nil, true, false),
				),
			},
			// update integration
			{
				Config: testAccIntegrationAmazon(domain, newName, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeAmazon, buddy.IntegrationAuthTypeDefault, scope, false, "", nil, true, false),
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
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccIntegrationCheckDestroy,
		Steps: []resource.TestStep{
			// create integration
			{
				Config: testAccIntegrationDigitalOcean(domain, name),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeDigitalOcean, "", buddy.IntegrationScopeWorkspace, false, "", nil, true, false),
				),
			},
			// update integration
			{
				Config: testAccIntegrationDigitalOcean(domain, newName),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeDigitalOcean, "", buddy.IntegrationScopeWorkspace, false, "", nil, true, false),
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
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeShopify, buddy.IntegrationAuthTypeToken, scope, true, "", nil, true, false),
				),
			},
			// update integration
			{
				Config: testAccIntegrationShopify(domain, newName, projectNameA, projectNameB, scope, projectNameB),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeShopify, buddy.IntegrationAuthTypeToken, scope, true, "", nil, true, false),
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
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeShopify, buddy.IntegrationAuthTypeTokenAppExtension, scope, true, "", nil, true, false),
				),
			},
			// update integration
			{
				Config: testAccIntegrationShopifyPartner(domain, newName, projectNameA, projectNameB, scope, projectNameB),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeShopify, buddy.IntegrationAuthTypeTokenAppExtension, scope, true, "", nil, true, false),
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
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeGitLab, "", scope, false, "", nil, true, false),
				),
			},
			// update integration
			{
				Config: testAccIntegrationGitLab(domain, newName, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeGitLab, "", scope, false, "", nil, true, false),
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
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeGitHub, "", scope, false, "", nil, true, false),
				),
			},
			// update integration
			{
				Config: testAccIntegrationGitHub(domain, newName, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeGitHub, "", scope, false, "", nil, true, false),
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
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeRackspace, "", scope, false, "", nil, true, false),
				),
			},
			// update integration
			{
				Config: testAccIntegrationRackspace(domain, newName, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeRackspace, "", scope, false, "", nil, true, false),
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
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeCloudflare, "", scope, false, "", nil, true, false),
				),
			},
			// update integration
			{
				Config: testAccIntegrationCloudflare(domain, newName, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeCloudflare, "", scope, false, "", nil, true, false),
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
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeUpcloud, "", scope, false, "", nil, true, false),
				),
			},
			// update integration
			{
				Config: testAccIntegrationUpcloud(domain, newName, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeUpcloud, "", scope, false, "", nil, true, false),
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
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeStackHawk, "", scope, false, "", nil, true, false),
				),
			},
			// update integration
			{
				Config: testAccIntegrationStackHawk(domain, newName, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeStackHawk, "", scope, false, "", nil, true, false),
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
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeGoogleServiceAccount, buddy.IntegrationAuthTypeOidc, scope, false, "", nil, true, false),
				),
			},
			// update integration
			{
				Config: testAccIntegrationAzureGoogleOidc(domain, newName, scope, newAudience, newGoogleProject),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeGoogleServiceAccount, buddy.IntegrationAuthTypeOidc, scope, false, "", nil, true, false),
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
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeAzureCloud, buddy.IntegrationAuthTypeOidc, scope, false, "", nil, true, false),
				),
			},
			// update integration
			{
				Config: testAccIntegrationAzureCloudOidc(domain, newName, scope, newAudience),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeAzureCloud, buddy.IntegrationAuthTypeOidc, scope, false, "", nil, true, false),
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
					testAccIntegrationAttributes("buddy_integration.bar", &integration, name, buddy.IntegrationTypeAzureCloud, buddy.IntegrationAuthTypeDefault, scope, false, "", nil, true, false),
				),
			},
			// update integration
			{
				Config: testAccIntegrationAzureCloud(domain, newName, scope),
				Check: resource.ComposeTestCheckFunc(
					testAccIntegrationGet("buddy_integration.bar", &integration),
					testAccIntegrationAttributes("buddy_integration.bar", &integration, newName, buddy.IntegrationTypeAzureCloud, buddy.IntegrationAuthTypeDefault, scope, false, "", nil, true, false),
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

func testAccIntegrationAttributes(n string, integration *buddy.Integration, name string, typ string, authType string, scope string, testScopeProject bool, identifier string, permissions *buddy.IntegrationPermissions, allPipelinesAllowed bool, allowOnePipeline bool) resource.TestCheckFunc {
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
		if err := util.CheckBoolFieldEqual("AllPipelinesAllowed", integration.AllPipelinesAllowed, allPipelinesAllowed); err != nil {
			return err
		}
		attrsAllPipelinesAllowed, _ := strconv.ParseBool(attrs["all_pipelines_allowed"])
		if err := util.CheckBoolFieldEqual("all_pipelines_allowed", attrsAllPipelinesAllowed, allPipelinesAllowed); err != nil {
			return err
		}
		if allowOnePipeline {
			if err := util.CheckBoolFieldEqual("AllowedPipelines[0].Id", integration.AllowedPipelines[0].Id > 0, true); err != nil {
				return err
			}
			attrsAllowedPipelineId, _ := strconv.Atoi(attrs["allowed_pipelines.0"])
			if err := util.CheckIntFieldEqualAndSet("allowed_pipelines.0", attrsAllowedPipelineId, integration.AllowedPipelines[0].Id); err != nil {
				return err
			}
		} else {
			if err := util.CheckIntFieldEqual("AllowedPipelines", len(integration.AllowedPipelines), 0); err != nil {
				return err
			}
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
		if err := util.CheckFieldEqualAndSet("integration_id", attrs["integration_id"], integration.HashId); err != nil {
			return err
		}
		if err := util.CheckFieldSet("html_url", attrs["html_url"]); err != nil {
			return err
		}
		if permissions != nil {
			if err := util.CheckFieldEqualAndSet("Permissions.Others", integration.Permissions.Others, permissions.Others); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("permissions.0.others", attrs["permissions.0.others"], permissions.Others); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("Permissions.Admins", integration.Permissions.Admins, permissions.Admins); err != nil {
				return err
			}
			if err := util.CheckFieldEqualAndSet("permissions.0.admins", attrs["permissions.0.admins"], permissions.Admins); err != nil {
				return err
			}
			if len(permissions.Users) > 0 {
				if err := util.CheckFieldEqualAndSet("Permissions.Users[0].AccessLevel", integration.Permissions.Users[0].AccessLevel, permissions.Users[0].AccessLevel); err != nil {
					return err
				}
				if err := util.CheckBoolFieldEqual("Permissions.Users[0].Id", integration.Permissions.Users[0].Id > 0, true); err != nil {
					return err
				}
				if err := util.CheckFieldEqualAndSet("permissions.0.user.0.access_level", attrs["permissions.0.user.0.access_level"], permissions.Users[0].AccessLevel); err != nil {
					return err
				}
				attrsPermUserId, _ := strconv.Atoi(attrs["permissions.0.user.0.id"])
				if err := util.CheckIntFieldEqualAndSet("permissions.0.user.0.id", attrsPermUserId, integration.Permissions.Users[0].Id); err != nil {
					return err
				}
			} else {
				if err := util.CheckIntFieldEqual("Permissions.Users", len(integration.Permissions.Users), 0); err != nil {
					return err
				}
			}
			if len(permissions.Groups) > 0 {
				if err := util.CheckFieldEqualAndSet("Permissions.Groups[0].AccessLevel", integration.Permissions.Groups[0].AccessLevel, permissions.Groups[0].AccessLevel); err != nil {
					return err
				}
				if err := util.CheckBoolFieldEqual("Permissions.Groups[0].Id", integration.Permissions.Groups[0].Id > 0, true); err != nil {
					return err
				}
				if err := util.CheckFieldEqualAndSet("permissions.0.group.0.access_level", attrs["permissions.0.group.0.access_level"], permissions.Groups[0].AccessLevel); err != nil {
					return err
				}
				attrsPermUserId, _ := strconv.Atoi(attrs["permissions.0.group.0.id"])
				if err := util.CheckIntFieldEqualAndSet("permissions.0.group.0.id", attrsPermUserId, integration.Permissions.Groups[0].Id); err != nil {
					return err
				}
			} else {
				if err := util.CheckIntFieldEqual("Permissions.Groups", len(integration.Permissions.Groups), 0); err != nil {
					return err
				}
			}
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

func testAccIntegrationAmazonOidc(domain string, name string, scope string, audience string, others string, admins string, email string, userAccessLevel string, groupName string, groupAccessLevel string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_member" "a" {
    domain = "${buddy_workspace.foo.domain}"
    email = "%s"
}

resource "buddy_group" "g" {
	domain = "${buddy_workspace.foo.domain}"
	name = "%s"
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

   permissions {
     others = "%s"
     admins = "%s"
     user {
       id = "${buddy_member.a.member_id}"
       access_level = "%s"
     }
     group {
       id = "${buddy_group.g.group_id}"
       access_level = "%s"
     }
   }
}
`, domain, email, groupName, name, buddy.IntegrationTypeAmazon, scope, audience, others, admins, userAccessLevel, groupAccessLevel)
}

func testAccIntegrationAmazonTrusted(domain string, name string, scope string, others string, admins string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_integration" "bar" {
   domain = "${buddy_workspace.foo.domain}"
   name = "%s"
   type = "%s"
   scope = "%s"
	 auth_type = "TRUSTED"

   role_assumption {
       arn = "arn1"
			 external_id = "123"
   }

   permissions {
     others = "%s"
     admins = "%s"
   }
}
`, domain, name, buddy.IntegrationTypeAmazon, scope, others, admins)
}

func testAccIntegrationAmazonDefault(domain string, name string, scope string, identifier string, allPipelinesAllowed bool, allowedPipeline string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_project" "p" {
   domain = "${buddy_workspace.foo.domain}" 
   display_name = "ppp"
}

resource "buddy_pipeline" "a" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.p.name}"
   name = "a"
   on = "CLICK"
}

resource "buddy_pipeline" "b" {
   domain = "${buddy_workspace.foo.domain}"
   project_name = "${buddy_project.p.name}"
   name = "b"
   on = "CLICK"
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

	 all_pipelines_allowed = %t
   allowed_pipelines = ["${buddy_pipeline.%s.pipeline_id}"]
}
`, domain, name, buddy.IntegrationTypeAmazon, scope, identifier, allPipelinesAllowed, allowedPipeline)
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

func testAccIntegrationDigitalOcean(domain string, name string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_integration" "bar" {
   domain = "${buddy_workspace.foo.domain}"
   name = "%s"
   type = "%s"
   scope = "%s"
   token = "ABC"
}
`, domain, name, buddy.IntegrationTypeDigitalOcean, buddy.IntegrationScopeWorkspace)
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
