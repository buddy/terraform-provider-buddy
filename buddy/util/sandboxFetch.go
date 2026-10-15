package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	sourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type sandboxFetchConfigModel struct {
	Type         types.String `tfsdk:"type"`
	Repository   types.String `tfsdk:"repository"`
	Ref          types.String `tfsdk:"ref"`
	Path         types.String `tfsdk:"path"`
	SourcePath   types.String `tfsdk:"source_path"`
	Artifact     types.String `tfsdk:"artifact"`
	BuildCommand types.String `tfsdk:"build_command"`
}

func ResourceSandboxFetchConfigModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"type": schema.StringAttribute{
			MarkdownDescription: "The fetch's type. Allowed: `PROJECT_REPO`, `PUBLIC_REPO`, `ARTIFACT`",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					buddy.SandboxFetchTypeProjectRepo,
					buddy.SandboxFetchTypePublicRepo,
					buddy.SandboxFetchTypeArtifact,
				),
			},
		},
		"repository": schema.StringAttribute{
			MarkdownDescription: "The URL of the git repository. Provide for `PUBLIC_REPO`",
			Optional:            true,
		},
		"ref": schema.StringAttribute{
			MarkdownDescription: "The branch, tag or commit to checkout. Defaults to the default branch. Provide for `PROJECT_REPO`, `PUBLIC_REPO`",
			Optional:            true,
		},
		"path": schema.StringAttribute{
			MarkdownDescription: "The target directory in the sandbox. Defaults to the sandbox's `app_dir`. Two fetches cannot share it",
			Optional:            true,
		},
		"source_path": schema.StringAttribute{
			MarkdownDescription: "The path inside the artifact to fetch instead of its whole content. Provide for `ARTIFACT`",
			Optional:            true,
		},
		"artifact": schema.StringAttribute{
			MarkdownDescription: "The artifact in format `name:version`. Provide for `ARTIFACT`",
			Optional:            true,
		},
		"build_command": schema.StringAttribute{
			MarkdownDescription: "The command run after fetching",
			Optional:            true,
		},
	}
}

func SandboxFetchConfigsToApi(ctx context.Context, l *types.List) (*[]*buddy.SandboxFetchOps, diag.Diagnostics) {
	var m []sandboxFetchConfigModel
	var diags diag.Diagnostics
	if !l.IsNull() && !l.IsUnknown() {
		diags = l.ElementsAs(ctx, &m, false)
	}
	fetches := make([]*buddy.SandboxFetchOps, len(m))
	for i, v := range m {
		f := buddy.SandboxFetchOps{
			Type: v.Type.ValueStringPointer(),
		}
		if !v.Repository.IsNull() && !v.Repository.IsUnknown() {
			f.Repository = v.Repository.ValueStringPointer()
		}
		if !v.Ref.IsNull() && !v.Ref.IsUnknown() {
			f.Ref = v.Ref.ValueStringPointer()
		}
		if !v.Path.IsNull() && !v.Path.IsUnknown() {
			f.Path = v.Path.ValueStringPointer()
		}
		if !v.SourcePath.IsNull() && !v.SourcePath.IsUnknown() {
			f.SourcePath = v.SourcePath.ValueStringPointer()
		}
		if !v.Artifact.IsNull() && !v.Artifact.IsUnknown() {
			f.Artifact = v.Artifact.ValueStringPointer()
		}
		if !v.BuildCommand.IsNull() && !v.BuildCommand.IsUnknown() {
			f.BuildCommand = v.BuildCommand.ValueStringPointer()
		}
		fetches[i] = &f
	}
	return &fetches, diags
}

type sandboxFetchModel struct {
	Id              types.String `tfsdk:"id"`
	Type            types.String `tfsdk:"type"`
	Repository      types.String `tfsdk:"repository"`
	Ref             types.String `tfsdk:"ref"`
	Path            types.String `tfsdk:"path"`
	SourcePath      types.String `tfsdk:"source_path"`
	Artifact        types.String `tfsdk:"artifact"`
	BuildCommand    types.String `tfsdk:"build_command"`
	Status          types.String `tfsdk:"status"`
	CurrentRevision types.String `tfsdk:"current_revision"`
	Behind          types.Int64  `tfsdk:"behind"`
	Outdated        types.Bool   `tfsdk:"outdated"`
	FetchedDate     types.String `tfsdk:"fetched_date"`
	ErrorMessage    types.String `tfsdk:"error_message"`
}

func sandboxFetchModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"id":               types.StringType,
		"type":             types.StringType,
		"repository":       types.StringType,
		"ref":              types.StringType,
		"path":             types.StringType,
		"source_path":      types.StringType,
		"artifact":         types.StringType,
		"build_command":    types.StringType,
		"status":           types.StringType,
		"current_revision": types.StringType,
		"behind":           types.Int64Type,
		"outdated":         types.BoolType,
		"fetched_date":     types.StringType,
		"error_message":    types.StringType,
	}
}

func (s *sandboxFetchModel) loadAPI(f *buddy.SandboxFetch) {
	s.Id = types.StringValue(f.Id)
	s.Type = types.StringValue(f.Type)
	s.Repository = types.StringValue(f.Repository)
	s.Ref = types.StringValue(f.Ref)
	s.Path = types.StringValue(f.Path)
	s.SourcePath = types.StringValue(f.SourcePath)
	s.Artifact = types.StringValue(f.Artifact)
	s.BuildCommand = types.StringValue(f.BuildCommand)
	s.Status = types.StringValue(f.Status)
	s.CurrentRevision = types.StringValue(f.CurrentRevision)
	if f.Behind != nil {
		s.Behind = types.Int64Value(int64(*f.Behind))
	} else {
		s.Behind = types.Int64Null()
	}
	if f.Outdated != nil {
		s.Outdated = types.BoolValue(*f.Outdated)
	} else {
		s.Outdated = types.BoolNull()
	}
	s.FetchedDate = types.StringValue(f.FetchedDate)
	s.ErrorMessage = types.StringValue(f.ErrorMessage)
}

func SandboxFetchesFromApi(ctx context.Context, fetches *[]*buddy.SandboxFetch) (basetypes.ListValue, diag.Diagnostics) {
	l := make([]*sandboxFetchModel, len(*fetches))
	for i, v := range *fetches {
		l[i] = &sandboxFetchModel{}
		l[i].loadAPI(v)
	}
	return types.ListValueFrom(ctx, types.ObjectType{AttrTypes: sandboxFetchModelAttrs()}, &l)
}

func sandboxFetchAttributeDescriptions() map[string]string {
	return map[string]string{
		"id":               "The fetch's ID",
		"type":             "The fetch's type",
		"repository":       "The URL of the git repository",
		"ref":              "The branch, tag or commit to checkout",
		"path":             "The target directory in the sandbox",
		"source_path":      "The path inside the artifact",
		"artifact":         "The artifact in format `name:version`",
		"build_command":    "The command run after fetching",
		"status":           "The fetch's status: `FETCHING`, `FETCHED`, `ERROR` or `STALE` (needs a refetch)",
		"current_revision": "The revision checked out in the sandbox (repository fetches)",
		"behind":           "How many commits the fetched revision is behind `ref`, null when it cannot be measured",
		"outdated":         "Whether the fetched artifact version differs from the current one, null when it cannot be measured",
		"fetched_date":     "When the source was last fetched, empty when it has not been fetched yet",
		"error_message":    "Why the last fetch failed",
	}
}

func ResourceSandboxFetchModelAttributes() map[string]schema.Attribute {
	attrs := map[string]schema.Attribute{}
	for name, desc := range sandboxFetchAttributeDescriptions() {
		switch sandboxFetchModelAttrs()[name] {
		case types.Int64Type:
			attrs[name] = schema.Int64Attribute{MarkdownDescription: desc, Computed: true}
		case types.BoolType:
			attrs[name] = schema.BoolAttribute{MarkdownDescription: desc, Computed: true}
		default:
			attrs[name] = schema.StringAttribute{MarkdownDescription: desc, Computed: true}
		}
	}
	return attrs
}

func SourceSandboxFetchModelAttributes() map[string]sourceschema.Attribute {
	attrs := map[string]sourceschema.Attribute{}
	for name, desc := range sandboxFetchAttributeDescriptions() {
		switch sandboxFetchModelAttrs()[name] {
		case types.Int64Type:
			attrs[name] = sourceschema.Int64Attribute{MarkdownDescription: desc, Computed: true}
		case types.BoolType:
			attrs[name] = sourceschema.BoolAttribute{MarkdownDescription: desc, Computed: true}
		default:
			attrs[name] = sourceschema.StringAttribute{MarkdownDescription: desc, Computed: true}
		}
	}
	return attrs
}
