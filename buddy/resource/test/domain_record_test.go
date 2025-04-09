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

func TestAccDomainRecord(t *testing.T) {
	var record buddy.Record
	workspaceDomain := util.UniqueString()
	domain := util.UniqueString() + ".com"
	name := util.UniqueString() + "." + domain
	typ := "A"
	ttl := 60
	value := "1.1.1.1"
	newTtl := 3600
	newValue := "2.2.2.2"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccDomainRecordDestroy,
		Steps: []resource.TestStep{
			// create domain & record
			{
				Config: testAccDomainRecordConfig(workspaceDomain, domain, name, typ, ttl, value),
				Check: resource.ComposeTestCheckFunc(
					testAccDomainRecordGet("buddy_domain_record.foo", &record),
					testAccDomainRecordAttributes("buddy_domain_record.foo", &record),
				),
			},
			// update record
			{
				Config: testAccDomainRecordConfig(workspaceDomain, domain, name, typ, newTtl, newValue),
				Check: resource.ComposeTestCheckFunc(
					testAccDomainRecordGet("buddy_domain_record.foo", &record),
					testAccDomainRecordAttributes("buddy_domain_record.foo", &record),
				),
			},
			// import record
			{
				ResourceName:      "buddy_domain_record.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDomainRecordAttributes(n string, record *buddy.Record) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsTtl, _ := strconv.Atoi(attrs["ttl"])
		if err := util.CheckIntFieldEqualAndSet("ttl", attrsTtl, record.Ttl); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("value.0", attrs["value.0"], record.Values[0]); err != nil {
			return err
		}
		return nil
	}
}

func testAccDomainRecordGet(n string, record *buddy.Record) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		workspaceDomain, domain, typ, err := util.DecomposeTripleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		r, _, err := acc.ApiClient.DomainService.GetRecord(workspaceDomain, domain, typ)
		if err != nil {
			return err
		}
		*record = *r
		return nil
	}
}

func testAccDomainRecordConfig(workspaceDomain string, domain string, name string, typ string, ttl int, value string) string {
	return fmt.Sprintf(`

  resource "buddy_workspace" "foo" {
	   domain = "%s"
	}

  resource "buddy_domain" "foo" {
     workspace_domain = "${buddy_workspace.foo.domain}"
     domain = "%s"
  }

  resource "buddy_domain_record" "foo" {
     depends_on = [buddy_domain.foo]
     workspace_domain = "${buddy_workspace.foo.domain}"
     domain = "%s"
     type = "%s"
     ttl = %d
     value = ["%s"]
  }
`, workspaceDomain, domain, name, typ, ttl, value)
}

func testAccDomainRecordDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_domain_record" {
			continue
		}
		workspaceDomain, name, typ, err := util.DecomposeTripleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		record, resp, err := acc.ApiClient.DomainService.GetRecord(workspaceDomain, name, typ)
		if err == nil && record != nil {
			return util.ErrorResourceExists()
		}
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}
