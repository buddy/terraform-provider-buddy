package acc

import (
	"fmt"
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
		t.Fatal("BUDDY_BASE_URL must be set for acceptance tests")
	}
}

func DummyCheckDestroy(_ *terraform.State) error {
	return nil
}

// MainWorkspaceDomain returns the workspace created together with the token (the oldest one), only it carries the token's plan
func MainWorkspaceDomain(t *testing.T) string {
	workspaces, _, err := ApiClient.WorkspaceService.GetList()
	if err != nil {
		t.Fatal(err)
	}
	var main *buddy.Workspace
	for _, w := range workspaces.Workspaces {
		if main == nil || w.Id < main.Id {
			main = w
		}
	}
	if main == nil {
		t.Fatal(fmt.Errorf("no workspace found"))
	}
	return main.Domain
}
