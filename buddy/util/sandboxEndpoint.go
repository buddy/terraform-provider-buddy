package util

import (
	"context"
	"github.com/buddy/api-go-sdk/buddy"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func sandboxEndpointHttpModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"verify_certificate":    types.BoolType,
		"compression":           types.BoolType,
		"http2":                 types.BoolType,
		"log_requests":          types.BoolType,
		"rewrite_host_header":   types.StringType,
		"whitelist_user_agents": types.SetType{ElemType: types.StringType},
		"request_headers":       types.MapType{ElemType: types.StringType},
		"response_headers":      types.MapType{ElemType: types.StringType},
		"login":                 types.StringType,
		"password":              types.StringType,
		"circuit_breaker":       types.Int32Type,
		"tls_ca":                types.StringType,
	}
}

func sandboxEndpointTlsModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"terminate_at":   types.StringType,
		"private_key":    types.StringType,
		"certificate":    types.StringType,
		"ca_certificate": types.StringType,
	}
}

func sandboxEndpointModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"endpoint":  types.StringType,
		"type":      types.StringType,
		"region":    types.StringType,
		"whitelist": types.SetType{ElemType: types.StringType},
		"timeout":   types.Int32Type,
		"http":      types.ObjectType{AttrTypes: sandboxEndpointHttpModelAttrs()},
		"tls":       types.ObjectType{AttrTypes: sandboxEndpointTlsModelAttrs()},
	}
}

type sandboxEndpointTlsModel struct {
	TerminateAt   types.String `tfsdk:"terminate_at"`
	PrivateKey    types.String `tfsdk:"private_key"`
	Certificate   types.String `tfsdk:"certificate"`
	CaCertificate types.String `tfsdk:"ca_certificate"`
}

type sandboxEndpointHttpModel struct {
	VerifyCertificate   types.Bool   `tfsdk:"verify_certificate"`
	Compression         types.Bool   `tfsdk:"compression"`
	Http2               types.Bool   `tfsdk:"http2"`
	LogRequests         types.Bool   `tfsdk:"log_requests"`
	RewriteHostHeader   types.String `tfsdk:"rewrite_host_header"`
	WhitelistUserAgents types.Set    `tfsdk:"whitelist_user_agents"`
	RequestHeaders      types.Map    `tfsdk:"request_headers"`
	ResponseHeaders     types.Map    `tfsdk:"response_headers"`
	Login               types.String `tfsdk:"login"`
	Password            types.String `tfsdk:"password"`
	CircuitBreaker      types.Int32  `tfsdk:"circuit_breaker"`
	TlsCa               types.String `tfsdk:"tls_ca"`
}

func (r *sandboxEndpointTlsModel) loadAPI(e *buddy.SandboxEndpointTls) {
	if e.TerminateAt != nil {
		r.TerminateAt = types.StringValue(*e.TerminateAt)
	} else {
		r.TerminateAt = types.StringNull()
	}
	if e.PrivateKey != nil {
		r.PrivateKey = types.StringValue(*e.PrivateKey)
	} else {
		r.PrivateKey = types.StringNull()
	}
	if e.Certificate != nil {
		r.Certificate = types.StringValue(*e.Certificate)
	} else {
		r.Certificate = types.StringNull()
	}
	if e.CaCertificate != nil {
		r.CaCertificate = types.StringValue(*e.CaCertificate)
	} else {
		r.CaCertificate = types.StringNull()
	}
}

func (r *sandboxEndpointTlsModel) toAPI() *buddy.SandboxEndpointTls {
	var e buddy.SandboxEndpointTls
	if !r.PrivateKey.IsNull() && !r.PrivateKey.IsUnknown() {
		e.PrivateKey = r.PrivateKey.ValueStringPointer()
	}
	if !r.Certificate.IsNull() && !r.Certificate.IsUnknown() {
		e.Certificate = r.Certificate.ValueStringPointer()
	}
	if !r.CaCertificate.IsNull() && !r.CaCertificate.IsUnknown() {
		e.CaCertificate = r.CaCertificate.ValueStringPointer()
	}
	if !r.TerminateAt.IsNull() && !r.TerminateAt.IsUnknown() {
		e.TerminateAt = r.TerminateAt.ValueStringPointer()
	}
	return &e
}

func (r *sandboxEndpointHttpModel) toAPI(ctx context.Context) (*buddy.SandboxEndpointHttp, diag.Diagnostics) {
	var diags diag.Diagnostics
	var e buddy.SandboxEndpointHttp
	if !r.VerifyCertificate.IsNull() && !r.VerifyCertificate.IsUnknown() {
		e.VerifyCertificate = r.VerifyCertificate.ValueBoolPointer()
	}
	if !r.Compression.IsNull() && !r.Compression.IsUnknown() {
		e.Compression = r.Compression.ValueBoolPointer()
	}
	if !r.Http2.IsNull() && !r.Http2.IsUnknown() {
		e.Http2 = r.Http2.ValueBoolPointer()
	}
	if !r.LogRequests.IsNull() && !r.LogRequests.IsUnknown() {
		e.LogRequests = r.LogRequests.ValueBoolPointer()
	}
	if !r.RewriteHostHeader.IsNull() && !r.RewriteHostHeader.IsUnknown() {
		e.RewriteHostHeader = r.RewriteHostHeader.ValueStringPointer()
	}
	if !r.WhitelistUserAgents.IsNull() && !r.WhitelistUserAgents.IsUnknown() {
		wh, d := StringSetToApi(ctx, &r.WhitelistUserAgents)
		diags.Append(d...)
		e.WhitelistUserAgents = wh
	}
	if !r.RequestHeaders.IsNull() && !r.ResponseHeaders.IsUnknown() {
		h, d := MapStringToApi(ctx, &r.RequestHeaders)
		diags.Append(d...)
		e.RequestHeaders = h
	}
	if !r.ResponseHeaders.IsNull() && !r.ResponseHeaders.IsUnknown() {
		h, d := MapStringToApi(ctx, &r.ResponseHeaders)
		diags.Append(d...)
		e.ResponseHeaders = h
	}
	if !r.Login.IsNull() && !r.Login.IsUnknown() {
		e.Login = r.Login.ValueStringPointer()
	}
	if !r.Password.IsNull() && !r.Password.IsUnknown() {
		e.Password = r.Password.ValueStringPointer()
	}
	if !r.CircuitBreaker.IsNull() && !r.CircuitBreaker.IsUnknown() {
		e.CircuitBreaker = PointerInt32(r.CircuitBreaker.ValueInt32())
	}
	if !r.TlsCa.IsNull() && !r.TlsCa.IsUnknown() {
		e.TlsCa = r.TlsCa.ValueStringPointer()
	}
	return &e, diags
}

func (r *sandboxEndpointHttpModel) loadAPI(ctx context.Context, e *buddy.SandboxEndpointHttp) diag.Diagnostics {
	var diags diag.Diagnostics
	if e.VerifyCertificate != nil {
		r.VerifyCertificate = types.BoolValue(*e.VerifyCertificate)
	} else {
		r.VerifyCertificate = types.BoolNull()
	}
	if e.Compression != nil {
		r.Compression = types.BoolValue(*e.Compression)
	} else {
		r.Compression = types.BoolNull()
	}
	if e.Http2 != nil {
		r.Http2 = types.BoolValue(*e.Http2)
	} else {
		r.Http2 = types.BoolNull()
	}
	if e.LogRequests != nil {
		r.LogRequests = types.BoolValue(*e.LogRequests)
	} else {
		r.LogRequests = types.BoolNull()
	}
	if e.RewriteHostHeader != nil {
		r.RewriteHostHeader = types.StringValue(*e.RewriteHostHeader)
	} else {
		r.RewriteHostHeader = types.StringNull()
	}
	if e.WhitelistUserAgents != nil {
		wh, d := types.SetValueFrom(ctx, types.StringType, *e.WhitelistUserAgents)
		diags.Append(d...)
		r.WhitelistUserAgents = wh
	} else {
		wh, d := types.SetValueFrom(ctx, types.StringType, []string{})
		diags.Append(d...)
		r.WhitelistUserAgents = wh
	}
	if e.RequestHeaders != nil {
		rh, d := types.MapValueFrom(ctx, types.StringType, *e.RequestHeaders)
		diags.Append(d...)
		r.RequestHeaders = rh
	} else {
		rh, d := types.MapValueFrom(ctx, types.StringType, map[string]string{})
		diags.Append(d...)
		r.RequestHeaders = rh
	}
	if e.ResponseHeaders != nil {
		rh, d := types.MapValueFrom(ctx, types.StringType, *e.ResponseHeaders)
		diags.Append(d...)
		r.ResponseHeaders = rh
	} else {
		rh, d := types.MapValueFrom(ctx, types.StringType, map[string]string{})
		diags.Append(d...)
		r.ResponseHeaders = rh
	}
	if e.Login != nil {
		r.Login = types.StringValue(*e.Login)
	} else {
		r.Login = types.StringNull()
	}
	if e.Password != nil {
		r.Password = types.StringValue(*e.Password)
	} else {
		r.Password = types.StringNull()
	}
	if e.CircuitBreaker != nil {
		r.CircuitBreaker = types.Int32Value(int32(*e.CircuitBreaker))
	} else {
		r.CircuitBreaker = types.Int32Null()
	}
	if e.TlsCa != nil {
		r.TlsCa = types.StringValue(*e.TlsCa)
	} else {
		r.TlsCa = types.StringNull()
	}
	return diags
}

type sandboxEndpointModel struct {
	Endpoint  types.String `tfsdk:"endpoint"`
	Type      types.String `tfsdk:"type"`
	Region    types.String `tfsdk:"region"`
	Whitelist types.Set    `tfsdk:"whitelist"`
	Timeout   types.Int32  `tfsdk:"timeout"`
	Http      types.Object `tfsdk:"http"`
	Tls       types.Object `tfsdk:"tls"`
}

func (r *sandboxEndpointModel) loadAPI(ctx context.Context, e *buddy.SandboxEndpoint) diag.Diagnostics {
	var diags diag.Diagnostics
	if e.Endpoint != nil {
		r.Endpoint = types.StringValue(*e.Endpoint)
	} else {
		r.Endpoint = types.StringNull()
	}
	if e.Type != nil {
		r.Type = types.StringValue(*e.Type)
	} else {
		r.Type = types.StringNull()
	}
	if e.Region != nil {
		r.Region = types.StringValue(*e.Region)
	} else {
		r.Region = types.StringNull()
	}
	if e.Whitelist != nil {
		wh, d := types.SetValueFrom(ctx, types.StringType, *e.Whitelist)
		diags.Append(d...)
		r.Whitelist = wh
	} else {
		wh, d := types.SetValueFrom(ctx, types.StringType, []string{})
		diags.Append(d...)
		r.Whitelist = wh
	}
	if e.Timeout != nil {
		r.Timeout = types.Int32Value(int32(*e.Timeout))
	} else {
		r.Timeout = types.Int32Null()
	}
	if e.Http != nil {
		var httpModel sandboxEndpointHttpModel
		d := httpModel.loadAPI(ctx, e.Http)
		diags.Append(d...)
		http, d := types.ObjectValueFrom(ctx, sandboxEndpointHttpModelAttrs(), httpModel)
		diags.Append(d...)
		r.Http = http
	} else {
		r.Http = types.ObjectNull(sandboxEndpointHttpModelAttrs())
	}
	if e.Tls != nil {
		var tlsModel sandboxEndpointTlsModel
		tlsModel.loadAPI(e.Tls)
		tls, d := types.ObjectValueFrom(ctx, sandboxEndpointTlsModelAttrs(), tlsModel)
		diags.Append(d...)
		r.Tls = tls
	} else {
		r.Tls = types.ObjectNull(sandboxEndpointTlsModelAttrs())
	}
	return diags
}

func ResourceSandboxEndpointTlsModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"terminate_at": schema.StringAttribute{
			Optional: true,
			Computed: true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					buddy.SandboxEndpointTlsTerminateAtRegion,
					buddy.SandboxEndpointTlsTerminateAtAgent,
					buddy.SandboxEndpointTlsTerminateAtTarget,
				),
			},
		},
		"private_key": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"certificate": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"ca_certificate": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
	}
}

func ResourceSandboxEndpointHttpModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"verify_certificate": schema.BoolAttribute{
			Optional: true,
			Computed: true,
		},
		"compression": schema.BoolAttribute{
			Optional: true,
			Computed: true,
		},
		"http2": schema.BoolAttribute{
			Optional: true,
			Computed: true,
		},
		"log_requests": schema.BoolAttribute{
			Optional: true,
			Computed: true,
		},
		"rewrite_host_header": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"whitelist_user_agents": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
		},
		"request_headers": schema.MapAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
		},
		"response_headers": schema.MapAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
		},
		"login": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"password": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
		"circuit_breaker": schema.Int32Attribute{
			Optional: true,
			Computed: true,
		},
		"tls_ca": schema.StringAttribute{
			Optional: true,
			Computed: true,
		},
	}
}

func ResourceSandboxEndpointModelAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"endpoint": schema.StringAttribute{
			Required: true,
		},
		"type": schema.StringAttribute{
			Required: true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					buddy.SandboxEndpointTypeTcp,
					buddy.SandboxEndpointTypeHttp,
					buddy.SandboxEndpointTypeTls,
				),
			},
		},
		"region": schema.StringAttribute{
			Optional: true,
			Computed: true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					buddy.SandboxEndpointRegionEu,
					buddy.SandboxEndpointRegionUs,
				),
			},
		},
		"whitelist": schema.SetAttribute{
			ElementType: types.StringType,
			Optional:    true,
			Computed:    true,
		},
		"timeout": schema.Int32Attribute{
			Optional: true,
			Computed: true,
		},
		"http": schema.SingleNestedAttribute{
			Optional:   true,
			Computed:   true,
			Attributes: ResourceSandboxEndpointHttpModelAttributes(),
		},
		"tls": schema.SingleNestedAttribute{
			Optional:   true,
			Computed:   true,
			Attributes: ResourceSandboxEndpointTlsModelAttributes(),
		},
	}
}

func SandboxEndpointsToApi(ctx context.Context, m *types.Map) (*[]buddy.SandboxEndpoint, diag.Diagnostics) {
	var e map[string]sandboxEndpointModel
	diags := m.ElementsAs(ctx, &e, false)
	endpoints := make([]buddy.SandboxEndpoint, len(e))
	i := 0
	for n, v := range e {
		endpoint := buddy.SandboxEndpoint{}
		endpoint.Name = &n
		if !v.Endpoint.IsNull() && !v.Endpoint.IsUnknown() {
			endpoint.Endpoint = v.Endpoint.ValueStringPointer()
		}
		if !v.Type.IsNull() && !v.Type.IsUnknown() {
			endpoint.Type = v.Type.ValueStringPointer()
		}
		if !v.Region.IsNull() && !v.Region.IsUnknown() {
			endpoint.Region = v.Region.ValueStringPointer()
		}
		if !v.Whitelist.IsNull() && !v.Whitelist.IsUnknown() {
			wh, d := StringSetToApi(ctx, &v.Whitelist)
			diags.Append(d...)
			endpoint.Whitelist = wh
		}
		if !v.Timeout.IsNull() && !v.Timeout.IsUnknown() {
			endpoint.Timeout = PointerInt32(v.Timeout.ValueInt32())
		}
		if !v.Http.IsNull() && !v.Http.IsUnknown() {
			var ehm sandboxEndpointHttpModel
			d := v.Http.As(ctx, &ehm, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})
			diags.Append(d...)
			endpointHttp, d := ehm.toAPI(ctx)
			diags.Append(d...)
			endpoint.Http = endpointHttp
		}
		if !v.Tls.IsNull() && !v.Tls.IsUnknown() {
			var etm sandboxEndpointTlsModel
			d := v.Tls.As(ctx, &etm, basetypes.ObjectAsOptions{
				UnhandledNullAsEmpty:    false,
				UnhandledUnknownAsEmpty: false,
			})
			diags.Append(d...)
			endpointTls := etm.toAPI()
			endpoint.Tls = endpointTls
		}
		endpoints[i] = endpoint
		i += 1
	}
	return &endpoints, diags
}

func SandboxEndpointsFromApi(ctx context.Context, endpoints *[]buddy.SandboxEndpoint) (basetypes.MapValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	m := map[string]sandboxEndpointModel{}
	if endpoints != nil {
		for _, v := range *endpoints {
			if v.Name != nil {
				e := sandboxEndpointModel{}
				d := e.loadAPI(ctx, &v)
				diags.Append(d...)
				m[*v.Name] = e
			}
		}
	}
	b, d := types.MapValueFrom(ctx, types.ObjectType{AttrTypes: sandboxEndpointModelAttrs()}, &m)
	diags.Append(d...)
	return b, d
}
