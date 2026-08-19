package dtc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/infobloxopen/terraform-provider-infoblox/internal/flex"
	uddidtc "github.com/infobloxopen/universal-ddi-go-client/dtc"
)

// TopologyRuleModel is the Terraform model for TopologyRule
type TopologyRuleModel struct {
	Code        types.String `tfsdk:"code"`
	Destination types.String `tfsdk:"destination"`
	Name        types.String `tfsdk:"name"`
	PoolId      types.String `tfsdk:"pool_id"`
	Source      types.String `tfsdk:"source"`
	Subnets     types.List   `tfsdk:"subnets"`
	Tags        types.List   `tfsdk:"tags"`
}

// TopologyRuleAttrTypes contains the attribute types for TopologyRuleModel
var TopologyRuleAttrTypes = map[string]attr.Type{
	"code":        types.StringType,
	"destination": types.StringType,
	"name":        types.StringType,
	"pool_id":     types.StringType,
	"source":      types.StringType,
	"subnets":     types.ListType{ElemType: types.StringType},
	"tags":        types.ListType{ElemType: types.ObjectType{AttrTypes: TagRuleAttrTypes}},
}

// TopologyRuleResourceSchemaAttributes contains the schema attributes for TopologyRuleModel
var TopologyRuleResourceSchemaAttributes = map[string]schema.Attribute{
	"code": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "Optional. DNS code to return if rule matches. Must be set if _destination_ is set to _code_.  Allowed values: - nodata - nxdomain  Defaults to _nodata_.",
	},
	"destination": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "Destination of __TopologyRule__.  Allowed values: - code - pool  Defaults to _code_.",
	},
	"name": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "Display name of __TopologyRule__.",
	},
	"pool_id": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "The resource identifier.",
	},
	"source": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "Type of source.  Allowed values: - subnet - tags - default  Defaults to _default_.",
	},
	"subnets": schema.ListAttribute{
		ElementType:         types.StringType,
		Optional:            true,
		MarkdownDescription: "Optional. List of subnets in CIDR format.  Must be set if _source_ is _subnet_, otherwise must be empty.",
	},
	"tags": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: TagRuleResourceSchemaAttributes,
		},
		Optional:            true,
		MarkdownDescription: "Optional. List of tag rules to match against a source object's effective tags. Effective tags = direct tags plus tags inherited from the IPAM parent chain (IPSpace → Address Block → Subnet); the closer level wins on key conflicts. All rules use AND semantics: an object must satisfy every __TagRule__ to match.  Must be set if _source_ is _tags_, otherwise must be empty.",
	},
}

// ExpandTopologyRule converts a Terraform Object to SDK type
func ExpandTopologyRule(ctx context.Context, o types.Object, diags *diag.Diagnostics) *uddidtc.TopologyRule {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m TopologyRuleModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

// Expand converts the Terraform model to SDK type
func (m *TopologyRuleModel) Expand(ctx context.Context, diags *diag.Diagnostics) *uddidtc.TopologyRule {
	if m == nil {
		return nil
	}
	to := &uddidtc.TopologyRule{
		Code:        flex.ExpandStringPointer(m.Code),
		Destination: flex.ExpandStringPointer(m.Destination),
		Name:        flex.ExpandString(m.Name),
		PoolId:      flex.ExpandStringPointer(m.PoolId),
		Source:      flex.ExpandStringPointer(m.Source),
		Subnets:     flex.ExpandFrameworkListString(ctx, m.Subnets, diags),
		Tags:        flex.ExpandFrameworkListNestedBlock(ctx, m.Tags, diags, ExpandTagRule),
	}
	return to
}

// FlattenTopologyRule converts an SDK type to Terraform Object
func FlattenTopologyRule(ctx context.Context, from *uddidtc.TopologyRule, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(TopologyRuleAttrTypes)
	}
	m := &TopologyRuleModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, TopologyRuleAttrTypes, m)
	diags.Append(d...)
	return t
}

// Flatten populates the Terraform model from SDK type
func (m *TopologyRuleModel) Flatten(ctx context.Context, from *uddidtc.TopologyRule, diags *diag.Diagnostics) {
	if from == nil || m == nil {
		return
	}
	m.Code = flex.FlattenStringPointer(from.Code)
	m.Destination = flex.FlattenStringPointer(from.Destination)
	m.Name = flex.FlattenString(from.Name)
	m.PoolId = flex.FlattenStringPointer(from.PoolId)
	m.Source = flex.FlattenStringPointer(from.Source)
	m.Subnets = flex.FlattenFrameworkListString(ctx, from.Subnets, diags)
	m.Tags = flex.FlattenFrameworkListNestedBlock(ctx, from.Tags, TagRuleAttrTypes, diags, FlattenTagRule)
}
