package acc

import (
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"os"
	"terraform-provider-buddy/buddy/provider"
	"testing"
)

var ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)
var ApiClient *buddy.Client

func init() {
	ApiClient, _ = buddy.NewClient(os.Getenv("BUDDY_TOKEN"), os.Getenv("BUDDY_BASE_URL"), os.Getenv("BUDDY_INSECURE") == "true")
	ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"buddy": providerserver.NewProtocol6WithError(provider.New("test")()),
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
