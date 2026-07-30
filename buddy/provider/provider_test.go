package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"testing"
)

func TestRegionBaseUrl(t *testing.T) {
	tests := []struct {
		region string
		want   string
		wantOk bool
	}{
		{"us", "https://api.buddy.works", true},
		{"as", "https://api.asia.buddy.works", true},
		{"eu", "https://api.eu.buddy.works", true},
		{"US", "https://api.buddy.works", true},
		{"Eu", "https://api.eu.buddy.works", true},
		{" as ", "https://api.asia.buddy.works", true},
		{"asia", "", false},
		{"", "", false},
	}
	for _, test := range tests {
		got, ok := regionBaseUrl(test.region)
		if ok != test.wantOk {
			t.Errorf("regionBaseUrl(%q) ok = %v, want %v", test.region, ok, test.wantOk)
			continue
		}
		if got != test.want {
			t.Errorf("regionBaseUrl(%q) = %q, want %q", test.region, got, test.want)
		}
	}
}

func TestResolveBaseUrl(t *testing.T) {
	tests := []struct {
		name       string
		baseUrl    types.String
		region     types.String
		envBaseUrl string
		envRegion  string
		want       string
		wantErr    bool
	}{
		{
			name:   "nothing set falls back to the client default",
			region: types.StringNull(), baseUrl: types.StringNull(),
			want: "",
		},
		{
			name:   "region sets the base url",
			region: types.StringValue("eu"), baseUrl: types.StringNull(),
			want: "https://api.eu.buddy.works",
		},
		{
			name:   "region is case insensitive",
			region: types.StringValue("AS"), baseUrl: types.StringNull(),
			want: "https://api.asia.buddy.works",
		},
		{
			name:   "base url wins over region",
			region: types.StringValue("eu"), baseUrl: types.StringValue("https://api.onprem.example.com"),
			want: "https://api.onprem.example.com",
		},
		{
			name:   "base url env wins over region attribute",
			region: types.StringValue("eu"), baseUrl: types.StringNull(),
			envBaseUrl: "https://api.onprem.example.com",
			want:       "https://api.onprem.example.com",
		},
		{
			name:   "region attribute wins over region env",
			region: types.StringValue("us"), baseUrl: types.StringNull(),
			envRegion: "eu",
			want:      "https://api.buddy.works",
		},
		{
			name:   "region env is used when the attribute is absent",
			region: types.StringNull(), baseUrl: types.StringNull(),
			envRegion: "as",
			want:      "https://api.asia.buddy.works",
		},
		{
			name:   "unknown region is an error",
			region: types.StringValue("antarctica"), baseUrl: types.StringNull(),
			wantErr: true,
		},
		{
			name:   "unknown region env is an error",
			region: types.StringNull(), baseUrl: types.StringNull(),
			envRegion: "antarctica",
			wantErr:   true,
		},
		{
			name:   "unknown region env is ignored when base url is set",
			region: types.StringNull(), baseUrl: types.StringValue("https://api.onprem.example.com"),
			envRegion: "antarctica",
			want:      "https://api.onprem.example.com",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv("BUDDY_BASE_URL", test.envBaseUrl)
			t.Setenv("BUDDY_REGION", test.envRegion)
			got, err := resolveBaseUrl(test.baseUrl, test.region)
			if test.wantErr {
				if err == nil {
					t.Fatalf("resolveBaseUrl() = %q, want an error", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("resolveBaseUrl() unexpected error: %s", err)
			}
			if got != test.want {
				t.Errorf("resolveBaseUrl() = %q, want %q", got, test.want)
			}
		})
	}
}
