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

func TestAccDomain(t *testing.T) {
	var domain buddy.Domain
	workspaceDomain := util.UniqueString()
	name := util.UniqueString() + ".com"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainConfig(workspaceDomain, name),
				Check: resource.ComposeTestCheckFunc(
					testAccDomainGet("buddy_domain.foo", &domain),
					testAccDomainAttributes("buddy_domain.foo", &domain, name, buddy.DomainTypePointed),
				),
			},
			{
				ResourceName:            "buddy_domain.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"on_owner_behalf"},
			},
		},
	})
}

func TestAccDomainPrivate(t *testing.T) {
	var domain buddy.Domain
	// private zones need a paid plan, which only the main workspace has
	workspaceDomain := acc.MainWorkspaceDomain(t)
	name := util.UniqueString() + ".lan"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccDomainDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainPrivateConfig(workspaceDomain, name),
				Check: resource.ComposeTestCheckFunc(
					testAccDomainGet("buddy_domain.foo", &domain),
					testAccDomainAttributes("buddy_domain.foo", &domain, name, buddy.DomainTypePrivate),
				),
			},
			{
				ResourceName:            "buddy_domain.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"on_owner_behalf"},
			},
		},
	})
}

func testAccDomainGet(n string, domain *buddy.Domain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		workspaceDomain, domainId, err := util.DecomposeDoubleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		d, _, err := acc.ApiClient.DomainService.Get(workspaceDomain, domainId)
		if err != nil {
			return err
		}
		*domain = *d
		return nil
	}
}

func testAccDomainAttributes(n string, domain *buddy.Domain, name string, typ string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		if err := util.CheckFieldEqualAndSet("domain", attrs["domain"], name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("domain.Name", domain.Name, name); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("type", attrs["type"], typ); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("domain.Type", domain.Type, typ); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("domain_id", attrs["domain_id"], domain.Id); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("html_url", attrs["html_url"], domain.HtmlUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqual("auto_renew", attrs["auto_renew"], "false"); err != nil {
			return err
		}
		return nil
	}
}

func testAccDomainConfig(workspaceDomain string, name string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_domain" "foo" {
   workspace_domain = "${buddy_workspace.foo.domain}"
   domain = "%s"
}
`, workspaceDomain, name)
}

func testAccDomainPrivateConfig(workspaceDomain string, name string) string {
	return fmt.Sprintf(`
resource "buddy_domain" "foo" {
   workspace_domain = "%s"
   domain = "%s"
   type = "PRIVATE"
}
`, workspaceDomain, name)
}

func testAccDomainDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_domain" {
			continue
		}
		workspaceDomain, domainId, err := util.DecomposeDoubleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		domain, resp, err := acc.ApiClient.DomainService.Get(workspaceDomain, domainId)
		if err == nil && domain != nil {
			return util.ErrorResourceExists()
		}
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}
