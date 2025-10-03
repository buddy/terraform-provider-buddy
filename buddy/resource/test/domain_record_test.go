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

func TestAccDomainGeoRecord(t *testing.T) {
	var record buddy.Record
	workspaceDomain := util.UniqueString()
	domain := util.UniqueString() + ".com"
	name := util.UniqueString() + "." + domain
	typ := "TXT"
	ttl := 300
	value := "A"
	continentName := buddy.DomainRecordContinentEurope
	continentValue := "B"
	countryName := buddy.DomainRecordCountryItaly
	countryValue := "C"
	newContinentName := buddy.DomainRecordContinentAsia
	newContinentValue := "G"
	newCountryName := buddy.DomainRecordCountryJapan
	newCountryValue := "F"
	newTtl := 600
	newValue := "B"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccDomainRecordDestroy,
		Steps: []resource.TestStep{
			// create domain & record
			{
				Config: testAccDomainGeoRecordConfig(workspaceDomain, domain, name, typ, ttl, value, countryName, countryValue, continentName, continentValue),
				Check: resource.ComposeTestCheckFunc(
					testAccDomainRecordGet("buddy_domain_record.foo", &record),
					testAccDomainRecordAttributes("buddy_domain_record.foo", &record),
				),
			},
			// update record
			{
				Config: testAccDomainGeoRecordConfig(workspaceDomain, domain, name, typ, newTtl, newValue, newCountryName, newCountryValue, newContinentName, newContinentValue),
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
		if err := util.CheckFieldEqualAndSet("routing", attrs["routing"], record.Routing); err != nil {
			return err
		}
		if err := util.CheckIntFieldEqualAndSet("ttl", attrsTtl, record.Ttl); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("value.0", attrs["value.0"], record.Values[0]); err != nil {
			return err
		}
		if record.Routing == buddy.DomainRecordRoutingGeolocation {
			if err := util.CheckFieldEqual("1 country", attrs["country.%"], "1"); err != nil {
				return err
			}
			if err := util.CheckFieldEqual("1 continent", attrs["continent.%"], "1"); err != nil {
				return err
			}
			i := 0
			for k, v := range record.Country {
				if err := util.CheckFieldEqual("country val", attrs[fmt.Sprintf("country.%s.%d", k, i)], v[i]); err != nil {
					return err
				}
				i += 1
			}
			i = 0
			for k, v := range record.Continent {
				if err := util.CheckFieldEqual("continent val", attrs[fmt.Sprintf("continent.%s.%d", k, i)], v[i]); err != nil {
					return err
				}
				i += 1
			}
		} else {
			if err := util.CheckFieldEqual("no country", attrs["country.%"], ""); err != nil {
				return err
			}
			if err := util.CheckFieldEqual("no continent", attrs["continent.%"], ""); err != nil {
				return err
			}
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
		workspaceDomain, domainId, domain, typ, err := util.DecomposeQuadrupleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		r, _, err := acc.ApiClient.DomainService.GetRecord(workspaceDomain, domainId, domain, typ)
		if err != nil {
			return err
		}
		*record = *r
		return nil
	}
}

func testAccDomainGeoRecordConfig(workspaceDomain string, domain string, name string, typ string, ttl int, value string, countryName string, countryValue string, continentName string, continentValue string) string {
	return fmt.Sprintf(`

  resource "buddy_workspace" "foo" {
	   domain = "%s"
	}

  resource "buddy_domain" "foo" {
     workspace_domain = "${buddy_workspace.foo.domain}"
     domain = "%s"
  }

  resource "buddy_domain_record" "foo" {
		 domain_id = "${buddy_domain.foo.domain_id}"
     workspace_domain = "${buddy_workspace.foo.domain}"
     domain = "%s"
     type = "%s"
     ttl = %d
     routing = "Geolocation"
     value = ["%s"]
     country = {
       %s = ["%s"]
		 }
     continent = {
       %s = ["%s"]
     }
  }
`, workspaceDomain, domain, name, typ, ttl, value, countryName, countryValue, continentName, continentValue)
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
     domain_id = "${buddy_domain.foo.domain_id}"
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
		workspaceDomain, domainId, name, typ, err := util.DecomposeQuadrupleId(rs.Primary.ID)
		if err != nil {
			return err
		}
		record, resp, err := acc.ApiClient.DomainService.GetRecord(workspaceDomain, domainId, name, typ)
		if err == nil && record != nil {
			return util.ErrorResourceExists()
		}
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}
