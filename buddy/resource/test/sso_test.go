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

// todo upgrade sso test

func TestAccSso(t *testing.T) {
	var sso buddy.Sso
	domain := util.UniqueString()
	ssoUrl := "https://login.microsoftonline.com/" + util.UniqueString() + "/saml2"
	issuer := "https://sts.windows.net/" + util.UniqueString()
	signature := buddy.SignatureMethodSha256
	digest := buddy.DigestMethodSha256
	err, cert := util.GenerateCertificate()
	if err != nil {
		t.Fatal(err.Error())
	}
	newSsoUrl := "https://login.microsoftonline.com/" + util.UniqueString() + "/saml2"
	newIssuer := "https://sts.windows.net/" + util.UniqueString()
	newSignature := buddy.SignatureMethodSha512
	newDigest := buddy.DigestMethodSha512
	err, newCert := util.GenerateCertificate()
	if err != nil {
		t.Fatal(err.Error())
	}
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acc.PreCheck(t)
		},
		ProtoV6ProviderFactories: acc.ProviderFactories,
		CheckDestroy:             testAccSsoCheckDestroy,
		Steps: []resource.TestStep{
			// create
			{
				Config: testAccSsoConfig(domain, ssoUrl, issuer, cert, signature, digest),
				Check: resource.ComposeTestCheckFunc(
					testAccSsoGet("buddy_sso.bar", &sso),
					testAccSsoAttributes("buddy_sso.bar", &sso, domain, ssoUrl, issuer, cert, signature, digest, false),
				),
			},
			// update
			{
				Config: testAccSsoConfig(domain, newSsoUrl, newIssuer, newCert, newSignature, newDigest),
				Check: resource.ComposeTestCheckFunc(
					testAccSsoGet("buddy_sso.bar", &sso),
					testAccSsoAttributes("buddy_sso.bar", &sso, domain, newSsoUrl, newIssuer, newCert, newSignature, newDigest, false),
				),
			},
			// require 4 all
			{
				Config: testAccSsoConfigRequireForAll(domain, newSsoUrl, newIssuer, newCert, newSignature, newDigest),
				Check: resource.ComposeTestCheckFunc(
					testAccSsoGet("buddy_sso.bar", &sso),
					testAccSsoAttributes("buddy_sso.bar", &sso, domain, newSsoUrl, newIssuer, newCert, newSignature, newDigest, true),
				),
			},
			{
				ResourceName:      "buddy_sso.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccSsoCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "buddy_sso" {
			continue
		}
		sso, resp, err := acc.ApiClient.SsoService.Get(rs.Primary.ID)
		if err == nil && sso != nil {
			return util.ErrorResourceExists()
		}
		if !util.IsResourceNotFound(resp, err) {
			return err
		}
	}
	return nil
}

func testAccSsoAttributes(n string, sso *buddy.Sso, domain string, ssoUrl string, issuer string, cert string, signature string, digest string, requireForAll bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		attrs := rs.Primary.Attributes
		attrsRequireForAll, _ := strconv.ParseBool(attrs["require_for_all"])
		if err := util.CheckFieldEqualAndSet("html_url", attrs["html_url"], sso.HtmlUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("domain", attrs["domain"], domain); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("sso_url", attrs["sso_url"], ssoUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("SsoUrl", sso.SsoUrl, ssoUrl); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("issuer", attrs["issuer"], issuer); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("Issuer", sso.Issuer, issuer); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("certificate", attrs["certificate"], cert); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("Certificate", sso.Certificate, cert); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("signature", attrs["signature"], signature); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("SignatureMethod", sso.SignatureMethod, signature); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("digest", attrs["digest"], digest); err != nil {
			return err
		}
		if err := util.CheckFieldEqualAndSet("DigestMethod", sso.DigestMethod, digest); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("require_for_all", attrsRequireForAll, requireForAll); err != nil {
			return err
		}
		if err := util.CheckBoolFieldEqual("RequireSsoForAllMembers", sso.RequireSsoForAllMembers, requireForAll); err != nil {
			return err
		}
		return nil
	}
}

func testAccSsoGet(n string, sso *buddy.Sso) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}
		ss, _, err := acc.ApiClient.SsoService.Get(rs.Primary.ID)
		if err != nil {
			return err
		}
		*sso = *ss
		return nil
	}
}

func testAccSsoConfig(domain string, ssoUrl string, issuer string, certificate string, signature string, digest string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_sso" "bar" {
   domain = "${buddy_workspace.foo.domain}"
   sso_url = "%s"
   issuer = "%s"
   certificate = <<EOT
%sEOT
   signature = "%s"
   digest = "%s"
}`, domain, ssoUrl, issuer, certificate, signature, digest)
}

func testAccSsoConfigRequireForAll(domain string, ssoUrl string, issuer string, certificate string, signature string, digest string) string {
	return fmt.Sprintf(`
resource "buddy_workspace" "foo" {
   domain = "%s"
}

resource "buddy_sso" "bar" {
   domain = "${buddy_workspace.foo.domain}"
   sso_url = "%s"
   issuer = "%s"
   certificate = <<EOT
%sEOT
   signature = "%s"
   digest = "%s"
   require_for_all = true
}`, domain, ssoUrl, issuer, certificate, signature, digest)
}
