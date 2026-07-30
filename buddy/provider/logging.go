package provider

import (
	"context"
	"encoding/json"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"regexp"
	"strings"
)

const maskedValue = "***"

// sensitiveJsonKeys are the API payload keys whose values must never reach the log,
// wherever in the payload they appear.
var sensitiveJsonKeys = map[string]bool{
	"access_key":             true,
	"api_key":                true,
	"certificate":            true,
	"client_secret":          true,
	"partner_token":          true,
	"passphrase":             true,
	"password":               true,
	"private_key":            true,
	"secret_key":             true,
	"token":                  true,
	"trigger_variable_value": true,
	"value":                  true,
}

// sensitiveNestedJsonKeys are keys that are only sensitive inside a given parent object.
// A target's SSH private key is `auth.key`, while a variable's `key` is its name and is
// what makes a log readable, so the two are told apart by where they sit.
var sensitiveNestedJsonKeys = map[string]string{
	"key": "auth",
}

var sensitiveJsonValue = regexp.MustCompile(`("(?:` + strings.Join(sortedSensitiveKeys(), "|") + `)"\s*:\s*)"(?:[^"\\]|\\.)*"`)

func sortedSensitiveKeys() []string {
	keys := make([]string, 0, len(sensitiveJsonKeys))
	for key := range sensitiveJsonKeys {
		keys = append(keys, key)
	}
	return keys
}

func isSensitiveKey(key string, parent string) bool {
	if sensitiveJsonKeys[key] {
		return true
	}
	requiredParent, ok := sensitiveNestedJsonKeys[key]
	return ok && requiredParent == parent
}

func maskJsonValue(value any, parent string) any {
	switch v := value.(type) {
	case map[string]any:
		for key, nested := range v {
			if isSensitiveKey(key, parent) {
				v[key] = maskedValue
			} else {
				v[key] = maskJsonValue(nested, key)
			}
		}
		return v
	case []any:
		for i, nested := range v {
			v[i] = maskJsonValue(nested, parent)
		}
		return v
	}
	return value
}

// maskSensitive redacts sensitive values in a logged payload.
//
// The message is a request or response line optionally followed by a JSON body. The body is
// masked structurally so that a key can be judged by its parent, falling back to a plain
// text substitution when it does not parse as JSON.
func maskSensitive(msg string) string {
	head, body, found := strings.Cut(msg, "\n")
	if !found {
		return msg
	}
	var parsed any
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		return head + "\n" + sensitiveJsonValue.ReplaceAllString(body, `${1}"`+maskedValue+`"`)
	}
	masked, err := json.MarshalIndent(maskJsonValue(parsed, ""), "", "\t")
	if err != nil {
		return head + "\n" + sensitiveJsonValue.ReplaceAllString(body, `${1}"`+maskedValue+`"`)
	}
	return head + "\n" + string(masked)
}

// newApiLogger routes the SDK's request logging into the provider's logger.
//
// The SDK must not write to stdout or stderr itself, Terraform reports anything a
// plugin puts there as `unexpected data` and the plugin protocol owns those streams.
// ctx is the one handed to Configure, it carries the logger for the provider's lifetime.
//
// Logged at info so that requests show under TF_LOG=INFO, DEBUG and TRACE alike.
func newApiLogger(ctx context.Context) buddy.LogFunc {
	return func(msg string) {
		tflog.Info(ctx, maskSensitive(msg))
	}
}
