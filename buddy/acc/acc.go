package acc

import (
	"buddy-terraform/buddy/provider"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"os"
	"testing"
)

var ApiClient *buddy.Client
var ProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	ApiClient, _ = buddy.NewClient(os.Getenv("BUDDY_TOKEN"), os.Getenv("BUDDY_BASE_URL"), os.Getenv("BUDDY_INSECURE") == "true")
	ProviderFactories = map[string]func() (*schema.Provider, error){
		"buddy": func() (*schema.Provider, error) { return provider.Provider(), nil },
	}
}

func PreCheck(t *testing.T) {
	if token := os.Getenv("BUDDY_TOKEN"); token == "" {
		t.Fatal("BUDDY_TOKEN must be set for acceptance tests")
	}
	if baseUrl := os.Getenv("BUDDY_BASE_URL"); baseUrl == "" {
		t.Fatal("BUDDY_BASE_URL must be set for acceptace tests")
	}
}

func DummyCheckDestroy(_ *terraform.State) error {
	return nil
}
