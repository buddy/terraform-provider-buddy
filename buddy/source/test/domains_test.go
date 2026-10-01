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

func TestAccSourceDomains(t *testing.T) {
	workspaceDomain := util.UniqueString()
	prefix := util.UniqueString()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		CheckDestroy:             acc.DummyCheckDestroy,
		ProtoV6ProviderFactories: acc.ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSourceDomainsConfig(workspaceDomain, prefix),
				Check: resource.ComposeTestCheckFunc(
					testAccSourceDomainsAttributes("data.buddy_domains.all", 2),
					testAccSourceDomainsAttributes("data.buddy_domains.pointed", 2),
					testAccSourceDomainsAttributes("data.buddy_domains.private", 0),
					testAccSourceDomainsAttributes("data.buddy_domains.filter", 1),
				),
			},
		},
	})
}

func testAccSourceDomainsAttributes(n string, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsDomainsCount, _ := strconv.Atoi(attrs["domains.#"])
		if err := util.CheckIntFieldEqual("domains.#", attrsDomainsCount, count); err != nil {
			return err
		}
		if count == 0 {
			return nil
		}
		if err := util.CheckFieldSet("domains.0.domain", attrs["domains.0.domain"]); err != nil {
			return err
		}
		if err := util.CheckFieldSet("domains.0.domain_id", attrs["domains.0.domain_id"]); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("domains.0.type", attrs["domains.0.type"], "POINTED"); err != nil {
			return err
		}
		if err := util.CheckFieldSet("domains.0.html_url", attrs["domains.0.html_url"]); err != nil {
			return err
		}
		return nil
	}
}

func testAccSourceDomainsConfig(workspaceDomain string, prefix string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_domain" "a" {
   workspace_domain = "${buddy_workspace.foo.domain}"
   domain = "%s-a.com"
}

resource "buddy_domain" "b" {
   workspace_domain = "${buddy_workspace.foo.domain}"
   domain = "%s-b.com"
}

data "buddy_domains" "all" {
   workspace_domain = "${buddy_workspace.foo.domain}"
   depends_on = [buddy_domain.a, buddy_domain.b]
}

data "buddy_domains" "pointed" {
   workspace_domain = "${buddy_workspace.foo.domain}"
   type = "POINTED"
   depends_on = [buddy_domain.a, buddy_domain.b]
}

data "buddy_domains" "private" {
   workspace_domain = "${buddy_workspace.foo.domain}"
   type = "PRIVATE"
   depends_on = [buddy_domain.a, buddy_domain.b]
}

data "buddy_domains" "filter" {
   workspace_domain = "${buddy_workspace.foo.domain}"
   domain_regex = "-b\\.com$"
   depends_on = [buddy_domain.a, buddy_domain.b]
}
`, workspaceDomain, prefix, prefix)
}
