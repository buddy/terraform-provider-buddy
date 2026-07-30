package provider

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMaskSensitive(t *testing.T) {
	tests := []struct {
		name string
		msg  string
		want map[string]any
	}{
		{
			name: "a variable value is redacted but its name is kept",
			msg:  "API Request POST /workspaces/abc/variables\n{\"key\":\"MY_VAR\",\"value\":\"s3cret\",\"type\":\"VAR\"}",
			want: map[string]any{"key": "MY_VAR", "value": "***", "type": "VAR"},
		},
		{
			name: "a target's auth key is redacted",
			msg:  "API Request POST /workspaces/abc/targets\n{\"name\":\"t1\",\"auth\":{\"method\":\"SSH_KEY\",\"key\":\"-----BEGIN RSA PRIVATE KEY-----\",\"passphrase\":\"pp\"}}",
			want: map[string]any{"name": "t1", "auth": map[string]any{"method": "SSH_KEY", "key": "***", "passphrase": "***"}},
		},
		{
			name: "sensitive keys nested in an array are redacted",
			msg:  "API Response GET /x\n{\"variables\":[{\"key\":\"A\",\"value\":\"one\"},{\"key\":\"B\",\"value\":\"two\"}]}",
			want: map[string]any{"variables": []any{
				map[string]any{"key": "A", "value": "***"},
				map[string]any{"key": "B", "value": "***"},
			}},
		},
		{
			name: "a key that merely ends with a sensitive name is untouched",
			msg:  "API Response GET /x\n{\"ca_certificate\":\"public\",\"public_value\":\"public\"}",
			want: map[string]any{"ca_certificate": "public", "public_value": "public"},
		},
		{
			name: "non sensitive payloads pass through",
			msg:  "API Response GET /workspaces 200 OK\n{\"name\":\"ws\",\"note\":\"hello\"}",
			want: map[string]any{"name": "ws", "note": "hello"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := maskSensitive(test.msg)
			head, body, found := strings.Cut(got, "\n")
			if !found {
				t.Fatalf("maskSensitive() dropped the body: %s", got)
			}
			if wantHead, _, _ := strings.Cut(test.msg, "\n"); head != wantHead {
				t.Errorf("maskSensitive() head = %q, want %q", head, wantHead)
			}
			var parsed map[string]any
			if err := json.Unmarshal([]byte(body), &parsed); err != nil {
				t.Fatalf("maskSensitive() produced invalid json: %s", body)
			}
			wantJson, _ := json.Marshal(test.want)
			gotJson, _ := json.Marshal(parsed)
			if string(gotJson) != string(wantJson) {
				t.Errorf("maskSensitive()\n got: %s\nwant: %s", gotJson, wantJson)
			}
		})
	}
}

func TestMaskSensitiveNonJsonBody(t *testing.T) {
	// a body that does not parse still gets the plain text substitution
	got := maskSensitive("API Response POST /x 500\n{\"value\": \"leaked\", truncated")
	if strings.Contains(got, "leaked") {
		t.Errorf("maskSensitive() leaked a value from an unparseable body: %s", got)
	}
}

func TestMaskSensitiveNoBody(t *testing.T) {
	msg := "API Request GET /workspaces/abc/variables?key=NAME"
	if got := maskSensitive(msg); got != msg {
		t.Errorf("maskSensitive() = %q, want it unchanged", got)
	}
}
