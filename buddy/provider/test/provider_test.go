package test

import (
	"buddy-terraform/buddy/acc"
	"buddy-terraform/buddy/util"
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"testing"
)

// todo move to workspace test
func TestAcc_Provider_UpgradeLatestMajor(t *testing.T) {
	config := fmt.Sprintf(`
		resource "buddy_workspace" "test" {
         domain = "%s"
       }
	`, util.UniqueString())
	//lintignore:AT001
	resource.Test(t, resource.TestCase{
		Steps: []resource.TestStep{
			{
				ExternalProviders: map[string]resource.ExternalProvider{
					"buddy": {
						VersionConstraint: "1.12.0",
						Source:            "buddy/buddy",
					},
				},
				Config: config,
			},
			// test migrating
			{
				ProtoV6ProviderFactories: acc.ProviderFactories,
				Config:                   config,
			},
		},
	})
}
