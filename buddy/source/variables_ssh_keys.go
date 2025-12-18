package source

import (
  "context"
  "github.com/buddy/api-go-sdk/buddy"
  "github.com/hashicorp/terraform-plugin-framework/datasource"
  "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
  "github.com/hashicorp/terraform-plugin-framework/diag"
  "github.com/hashicorp/terraform-plugin-framework/schema/validator"
  "github.com/hashicorp/terraform-plugin-framework/types"
  "regexp"
  "terraform-provider-buddy/buddy/util"
)

var (
  _ datasource.DataSource              = &variablesSshKeysSource{}
  _ datasource.DataSourceWithConfigure = &variablesSshKeysSource{}
)

func NewVariablesSshKeysSource() datasource.DataSource {
  return &variablesSshKeysSource{}
}

type variablesSshKeysSource struct {
  client *buddy.Client
}

type variablesSshKeysSourceModel struct {
  ID            types.String `tfsdk:"id"`
  Domain        types.String `tfsdk:"domain"`
  KeyRegex      types.String `tfsdk:"key_regex"`
  ProjectName   types.String `tfsdk:"project_name"`
  PipelineId    types.Int64  `tfsdk:"pipeline_id"`
  ActionId      types.Int64  `tfsdk:"action_id"`
  EnvironmentId types.String `tfsdk:"environment_id"`
  Variables     types.Set    `tfsdk:"variables"`
}

func (s *variablesSshKeysSourceModel) loadAPI(ctx context.Context, domain string, variables *[]*buddy.Variable) diag.Diagnostics {
  s.ID = types.StringValue(util.UniqueString())
  s.Domain = types.StringValue(domain)
  v, d := util.VariablesSshKeysModelFromApi(ctx, variables)
  s.Variables = v
  return d
}

func (s *variablesSshKeysSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
  resp.TypeName = req.ProviderTypeName + "_variables_ssh_keys"
}

func (s *variablesSshKeysSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
  if req.ProviderData == nil {
    return
  }
  s.client = req.ProviderData.(*buddy.Client)
}

func (s *variablesSshKeysSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
  resp.Schema = schema.Schema{
    MarkdownDescription: "List variables of SSH key type and optionally filter them by key, project name, pipeline or action\n\n" +
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
      "key_regex": schema.StringAttribute{
        MarkdownDescription: "The variable's key regular expression to match",
        Optional:            true,
        Validators: []validator.String{
          util.RegexpValidator(),
        },
      },
      "project_name": schema.StringAttribute{
        MarkdownDescription: "Get only from provided project",
        Optional:            true,
      },
      "pipeline_id": schema.Int64Attribute{
        MarkdownDescription: "Get only from provided pipeline",
        Optional:            true,
      },
      "action_id": schema.Int64Attribute{
        MarkdownDescription: "Get only from provided action",
        Optional:            true,
      },
      "environment_id": schema.StringAttribute{
        MarkdownDescription: "Get only from provided environment",
        Optional:            true,
      },
      "variables": schema.SetNestedAttribute{
        MarkdownDescription: "List of variables",
        Computed:            true,
        NestedObject: schema.NestedAttributeObject{
          Attributes: util.SourceVariableSshKeyModelAttributes(),
        },
      },
    },
  }
}

func (s *variablesSshKeysSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
  var data *variablesSshKeysSourceModel
  resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
  if resp.Diagnostics.HasError() {
    return
  }
  domain := data.Domain.ValueString()
  var keyRegex *regexp.Regexp
  if !data.KeyRegex.IsNull() && !data.KeyRegex.IsUnknown() {
    keyRegex = regexp.MustCompile(data.KeyRegex.ValueString())
  }
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
  if !data.EnvironmentId.IsNull() && !data.EnvironmentId.IsUnknown() {
    ops.EnvironmentId = data.EnvironmentId.ValueString()
  }
  variables, _, err := s.client.VariableService.GetList(domain, &ops)
  if err != nil {
    resp.Diagnostics.Append(util.NewDiagnosticApiError("get variables", err))
    return
  }
  var result []*buddy.Variable
  for _, v := range variables.Variables {
    if v.Type != buddy.VariableTypeSshKey {
      continue
    }
    if keyRegex != nil && !keyRegex.MatchString(v.Key) {
      continue
    }
    result = append(result, v)
  }
  resp.Diagnostics.Append(data.loadAPI(ctx, domain, &result)...)
  if resp.Diagnostics.HasError() {
    return
  }
  resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
