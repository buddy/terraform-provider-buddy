package test

import (
	"buddy-terraform/buddy/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"testing"
)

func TestAccProvider(t *testing.T) {
	if err := provider.Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestAccProvider_imp(t *testing.T) {
	var _ *schema.Provider = provider.Provider()
}
