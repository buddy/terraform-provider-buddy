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
	_ datasource.DataSource              = &variableSource{}
	_ datasource.DataSourceWithConfigure = &variableSource{}
)

func NewVariableSource() datasource.DataSource {
	return &variableSource{}
}

type variableSource struct {
	client *buddy.Client
}

type variableSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Domain      types.String `tfsdk:"domain"`
	Key         types.String `tfsdk:"key"`
	VariableId  types.Int64  `tfsdk:"variable_id"`
	ProjectName types.String `tfsdk:"project_name"`
	PipelineId  types.Int64  `tfsdk:"pipeline_id"`
	ActionId    types.Int64  `tfsdk:"action_id"`
	Encrypted   types.Bool   `tfsdk:"encrypted"`
	Settable    types.Bool   `tfsdk:"settable"`
	Description types.String `tfsdk:"description"`
	Value       types.String `tfsdk:"value"`
}

func (s *variableSourceModel) loadAPI(domain string, variable *buddy.Variable) {
	s.ID = types.StringValue(util.ComposeDoubleId(domain, strconv.Itoa(variable.Id)))
	s.Domain = types.StringValue(domain)
	s.Key = types.StringValue(variable.Key)
	s.VariableId = types.Int64Value(int64(variable.Id))
	s.Encrypted = types.BoolValue(variable.Encrypted)
	s.Settable = types.BoolValue(variable.Settable)
	s.Description = types.StringValue(variable.Description)
	s.Value = types.StringValue(variable.Value)
}

func (s *variableSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_variable"
}

func (s *variableSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	s.client = req.ProviderData.(*buddy.Client)
}

func (s *variableSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get variable by key or variable ID\n\n" +
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
				Validators: []validator.String{
					stringvalidator.AlsoRequires(path.Expressions{
						path.MatchRoot("key"),
					}...),
				},
			},
			"pipeline_id": schema.Int64Attribute{
				MarkdownDescription: "The variable's pipeline ID",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AlsoRequires(path.Expressions{
						path.MatchRoot("key"),
					}...),
				},
			},
			"action_id": schema.Int64Attribute{
				MarkdownDescription: "The variable's action ID",
				Optional:            true,
				Validators: []validator.Int64{
					int64validator.AlsoRequires(path.Expressions{
						path.MatchRoot("key"),
					}...),
				},
			},
			"encrypted": schema.BoolAttribute{
				MarkdownDescription: "Is the variable's value encrypted",
				Computed:            true,
			},
			"settable": schema.BoolAttribute{
				MarkdownDescription: "Is the variable's value changeable",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The variable's description",
				Computed:            true,
			},
			"value": schema.StringAttribute{
				MarkdownDescription: "The variable's value. Encrypted if **encrypted** == true",
				Computed:            true,
				Sensitive:           true,
			},
		},
	}
}

func (s *variableSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *variableSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	domain := data.Domain.ValueString()
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
		if variable.Type != buddy.VariableTypeVar {
			resp.Diagnostics.Append(util.NewDiagnosticApiNotFound("variable"))
			return
		}
	} else {
		key := data.Key.ValueString()
		ops := buddy.VariableGetListQuery{}
		if !data.ProjectName.IsNull() && !data.ProjectName.IsUnknown() {
			ops.ProjectName = data.ProjectName.ValueString()
		}
		if !data.PipelineId.IsNull() && !data.PipelineId.IsUnknown() {
			ops.PipelineId = int(data.PipelineId.ValueInt64())
		}
		if !data.ActionId.IsNull() && !data.ActionId.IsUnknown() {
			ops.ActionId = int(data.ActionId.ValueInt64())
		}
		var variables *buddy.Variables
		variables, _, err = s.client.VariableService.GetList(domain, &ops)
		if err != nil {
			resp.Diagnostics.Append(util.NewDiagnosticApiError("get variables", err))
			return
		}
		for _, v := range variables.Variables {
			if v.Type != buddy.VariableTypeVar {
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
	data.loadAPI(domain, variable)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
