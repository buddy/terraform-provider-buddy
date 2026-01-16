package source

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"strconv"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ datasource.DataSource              = &variableSshKeySource{}
	_ datasource.DataSourceWithConfigure = &variableSshKeySource{}
)

func NewVariableSshKeySource() datasource.DataSource {
	return &variableSshKeySource{}
}

type variableSshKeySource struct {
	client *buddy.Client
}

type variableSshKeySourceModel struct {
	ID             types.String `tfsdk:"id"`
	Domain         types.String `tfsdk:"domain"`
	Key            types.String `tfsdk:"key"`
	VariableId     types.Int64  `tfsdk:"variable_id"`
	ProjectName    types.String `tfsdk:"project_name"`
	PipelineId     types.Int64  `tfsdk:"pipeline_id"`
	ActionId       types.Int64  `tfsdk:"action_id"`
	EnvironmentId  types.String `tfsdk:"environment_id"`
	Encrypted      types.Bool   `tfsdk:"encrypted"`
	Settable       types.Bool   `tfsdk:"settable"`
	Description    types.String `tfsdk:"description"`
	Value          types.String `tfsdk:"value"`
	PublicValue    types.String `tfsdk:"public_value"`
	KeyFingerprint types.String `tfsdk:"key_fingerprint"`
	Checksum       types.String `tfsdk:"checksum"`
	FileChmod      types.String `tfsdk:"file_chmod"`
	FilePath       types.String `tfsdk:"file_path"`
	FilePlace      types.String `tfsdk:"file_place"`
}

func (s *variableSshKeySourceModel) loadAPI(domain string, variable *buddy.Variable, ops *buddy.VariableGetListQuery) {
	s.ID = types.StringValue(util.ComposeDoubleId(domain, strconv.Itoa(variable.Id)))
	s.Domain = types.StringValue(domain)
	s.Key = types.StringValue(variable.Key)
	s.VariableId = types.Int64Value(int64(variable.Id))
	s.Encrypted = types.BoolValue(variable.Encrypted)
	s.Settable = types.BoolValue(variable.Settable)
	s.Description = types.StringValue(variable.Description)
	s.Value = types.StringValue(variable.Value)
	s.PublicValue = types.StringValue(variable.PublicValue)
	s.KeyFingerprint = types.StringValue(variable.KeyFingerprint)
	s.Checksum = types.StringValue(variable.Checksum)
	s.FileChmod = types.StringValue(variable.FileChmod)
	s.FilePath = types.StringValue(variable.FilePath)
	s.FilePlace = types.StringValue(variable.FilePlace)
	if variable.Project != nil {
		s.ProjectName = types.StringValue(variable.Project.Name)
	} else if ops != nil && ops.ProjectName != "" {
		s.ProjectName = types.StringValue(ops.ProjectName)
	} else {
		s.ProjectName = types.StringNull()
	}
	if variable.Pipeline != nil {
		s.PipelineId = types.Int64Value(int64(variable.Pipeline.Id))
	} else if ops != nil && ops.PipelineId != 0 {
		s.PipelineId = types.Int64Value(int64(ops.PipelineId))
	} else {
		s.PipelineId = types.Int64Null()
	}
	if variable.Action != nil {
		s.ActionId = types.Int64Value(int64(variable.Action.Id))
	} else if ops != nil && ops.ActionId != 0 {
		s.ActionId = types.Int64Value(int64(ops.ActionId))
	} else {
		s.ActionId = types.Int64Null()
	}
	if variable.Environment != nil {
		s.EnvironmentId = types.StringValue(variable.Environment.Id)
	} else if ops != nil && ops.EnvironmentId != "" {
		s.EnvironmentId = types.StringValue(ops.EnvironmentId)
	} else {
		s.EnvironmentId = types.StringNull()
	}
}

func (s *variableSshKeySource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variable_ssh_key"
}

func (s *variableSshKeySource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *variableSshKeySource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get variables of SSH key type by key or variable ID\n\n" +
			"Token scope required: `WORKSPACE`, `VARIABLE_INFO`",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Terraform resource identifier for this item",
				Computed:            true,
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The workspace's URL handle",
				Required:            true,
				Validators:          util.StringValidatorsDomain(),
			},
			"key": schema.StringAttribute{
				MarkdownDescription: "The variable's name",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("variable_id"),
						path.MatchRoot("key"),
					}...),
				},
			},
			"variable_id": schema.Int64Attribute{
				MarkdownDescription: "The variable's ID",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("variable_id"),
						path.MatchRoot("key"),
					}...),
				},
			},
			"project_name": schema.StringAttribute{
				MarkdownDescription: "The variable's project name",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("key"),
					}...),
				},
			},
			"pipeline_id": schema.Int64Attribute{
				MarkdownDescription: "The variable's pipeline ID",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.AlsoRequires(path.Expressions{
						path.MatchRoot("key"),
					}...),
				},
			},
			"action_id": schema.Int64Attribute{
				MarkdownDescription: "The variable's action ID",
				Optional:            true,
				Computed:            true,
				Validators: []validator.Int64{
					int64validator.AlsoRequires(path.Expressions{
						path.MatchRoot("key"),
					}...),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "The variable's environment ID",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("key"),
					}...),
				},
			},
			"encrypted": schema.BoolAttribute{
				MarkdownDescription: "Is the variable's value encrypted, always true for buddy_variable_ssh_key",
				Computed:            true,
			},
			"settable": schema.BoolAttribute{
				MarkdownDescription: "Is the variable's value changeable, always false for buddy_variable_ssh_key",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The variable's description",
				Computed:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "The variable's value, always encrypted for buddy_variable_ssh_key",
				Computed:            true,
				Sensitive:           true,
			},
			"public_value": schema.StringAttribute{
				MarkdownDescription: "The variable's public key",
				Computed:            true,
			},
			"key_fingerprint": schema.StringAttribute{
				MarkdownDescription: "The variable's fingerprint",
				Computed:            true,
			},
			"checksum": schema.StringAttribute{
				MarkdownDescription: "The variable's checksum",
				Computed:            true,
			},
			"file_chmod": schema.StringAttribute{
				MarkdownDescription: "The variable's file permission in an action's container",
				Computed:            true,
			},
			"file_path": schema.StringAttribute{
				MarkdownDescription: "The variable's path in the action's container",
				Computed:            true,
			},
			"file_place": schema.StringAttribute{
				MarkdownDescription: "Should the variable's be copied to an action's container in **file_path** (`CONTAINER`, `NONE`)",
				Computed:            true,
			},
		},
	}
}

func (s *variableSshKeySource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *variableSshKeySourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
	ops := buddy.VariableGetListQuery{}
	var variable *buddy.Variable
	var err error
	if !data.VariableId.IsNull() && !data.VariableId.IsUnknown() {
		var httpRes *http.Response
		varId := int(data.VariableId.ValueInt64())
		variable, httpRes, err = s.client.VariableService.Get(domain, varId)
		if err != nil {
			if util.IsResourceNotFound(httpRes, err) {
				resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("variable"))
				return
			}
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get variable", err))
			return
		}
		if variable.Type != buddy.VariableTypeSshKey {
			resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("variable"))
			return
		}
	} else {
		key := data.Key.ValueString()
		if !data.ProjectName.IsNull() && !data.ProjectName.IsUnknown() {
			ops.ProjectName = data.ProjectName.ValueString()
		}
		if !data.PipelineId.IsNull() && !data.PipelineId.IsUnknown() {
			ops.PipelineId = int(data.PipelineId.ValueInt64())
		}
		if !data.ActionId.IsNull() && !data.ActionId.IsUnknown() {
			ops.ActionId = int(data.ActionId.ValueInt64())
		}
		if !data.EnvironmentId.IsNull() && !data.EnvironmentId.IsUnknown() {
			ops.EnvironmentId = data.EnvironmentId.ValueString()
		}
		var variables *buddy.Variables
		variables, _, err = s.client.VariableService.GetList(domain, &ops)
		if err != nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get variables", err))
			return
		}
		for _, v := range variables.Variables {
			if v.Type != buddy.VariableTypeSshKey {
				continue
			}
			if v.Key == key {
				variable = v
				break
			}
		}
		if variable == nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("variable"))
			return
		}
	}
	data.loadAPI(domain, variable, &ops)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
