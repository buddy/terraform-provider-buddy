package util

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"regexp"
)

var _ validator.String = regexpValidator{}

type regexpValidator struct {
}

func (v regexpValidator) Description(_ context.Context) string {
	return "value must be a valid regex"
}

func (v regexpValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v regexpValidator) ValidateString(ctx context.Context, req validator.StringRequest, res *validator.StringResponse) {
	val := req.ConfigValue
	if !val.IsNull() && !val.IsUnknown() {
		str := val.ValueString()
		if _, err := regexp.Compile(str); err != nil {
			res.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
				req.Path,
				v.Description(ctx),
				str,
			))
		}
	}
}

func RegexpValidator() validator.String {
	return regexpValidator{}
}
