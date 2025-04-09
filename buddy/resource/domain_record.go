package resource

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-buddy/buddy/util"
)

var (
	_ resource.Resource                = &domainRecordResource{}
	_ resource.ResourceWithConfigure   = &domainRecordResource{}
	_ resource.ResourceWithImportState = &domainRecordResource{}
)

func NewDomainRecordResource() resource.Resource {
	return &domainRecordResource{}
}

type domainRecordResource struct {
	client *buddy.Client
}

type domainRecordResourceModel struct {
	ID              types.String `tfsdk:"id"`
	WorkspaceDomain types.String `tfsdk:"workspace_domain"`
	Domain          types.String `tfsdk:"domain"`
	Type            types.String `tfsdk:"type"`
	Ttl             types.Int64  `tfsdk:"ttl"`
	Value           types.List   `tfsdk:"value"`
}

func (r *domainRecordResourceModel) decomposeId() (string, string, string, error) {
	workspaceDomain, domain, typ, err := util.DecomposeTripleId(r.ID.ValueString())
	if err != nil {
		return "", "", "", err
	}
	return workspaceDomain, domain, typ, nil
}

func (r *domainRecordResourceModel) loadAPI(ctx context.Context, workspaceDomain string, domain string, record *buddy.Record) diag.Diagnostics {
	var diags diag.Diagnostics
	r.ID = types.StringValue(util.ComposeTripleId(workspaceDomain, domain, record.Type))
	r.WorkspaceDomain = types.StringValue(workspaceDomain)
	r.Domain = types.StringValue(domain)
	r.Type = types.StringValue(record.Type)
	r.Ttl = types.Int64Value(int64(record.Ttl))
	value, d := types.ListValueFrom(ctx, types.StringType, &record.Values)
	diags.Append(d...)
	r.Value = value
	return diags
}

func (r *domainRecordResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain_record"
}

func (r *domainRecordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage a domain record\n\n" +
			"Token scope required: `ZONE_READ, ZONE_WRITE`",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The Terraform resource identifier for this item",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"workspace_domain": schema.StringAttribute{
				MarkdownDescription: "The workspace's URL handle",
				Required:            true,
				Validators:          util.StringValidatorsDomain(),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"domain": schema.StringAttribute{
				MarkdownDescription: "The record's full domain name",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "The record's type",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"A",
						"AAAA",
						"CNAME",
						"MX",
						"TXT",
						"SRV",
						"SPF",
						"NAPTR",
						"CAA",
						"NS",
						"SOA",
					),
				},
			},
			"ttl": schema.Int64Attribute{
				MarkdownDescription: "The record's ttl",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(300),
			},
			"value": schema.ListAttribute{
				MarkdownDescription: "The record's value list",
				Required:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (r *domainRecordResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*buddy.Client)
}

func (r *domainRecordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *domainRecordResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	workspaceDomain := data.WorkspaceDomain.ValueString()
	domain := data.Domain.ValueString()
	typ := data.Type.ValueString()
	ttl := int(data.Ttl.ValueInt64())
	value, d := util.StringListToApi(ctx, &data.Value)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	ops := buddy.RecordUpsertOps{
		Ttl:    &ttl,
		Values: value,
	}
	record, _, err := r.client.DomainService.UpsertRecord(workspaceDomain, domain, typ, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("upsert record", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, workspaceDomain, domain, record)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *domainRecordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *domainRecordResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	workspaceDomain, domain, typ, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("domain_record", err))
		return
	}
	record, httpResp, err := r.client.DomainService.GetRecord(workspaceDomain, domain, typ)
	if err != nil {
		if util.IsResourceNotFound(httpResp, err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.Append(util.NewDiagnosticApiError("get record", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, workspaceDomain, domain, record)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *domainRecordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *domainRecordResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	workspaceDomain, domain, typ, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("domain_record", err))
		return
	}
	ttl := int(data.Ttl.ValueInt64())
	value, d := util.StringListToApi(ctx, &data.Value)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	ops := buddy.RecordUpsertOps{
		Ttl:    &ttl,
		Values: value,
	}
	record, _, err := r.client.DomainService.UpsertRecord(workspaceDomain, domain, typ, &ops)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("update record", err))
		return
	}
	resp.Diagnostics.Append(data.loadAPI(ctx, workspaceDomain, domain, record)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *domainRecordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *domainRecordResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	workspaceDomain, domain, typ, err := data.decomposeId()
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticDecomposeError("domain_record", err))
		return
	}
	_, err = r.client.DomainService.DeleteRecord(workspaceDomain, domain, typ)
	if err != nil {
		resp.Diagnostics.Append(util.NewDiagnosticApiError("delete record", err))
	}
}

func (r *domainRecordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
