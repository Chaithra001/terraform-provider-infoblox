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

// PolicyTopologyModel is the Terraform model for PolicyTopology
type PolicyTopologyModel struct {
	Rules      types.List   `tfsdk:"rules"`
	TopologyId types.String `tfsdk:"topology_id"`
}

// PolicyTopologyAttrTypes contains the attribute types for PolicyTopologyModel
var PolicyTopologyAttrTypes = map[string]attr.Type{
	"rules":       types.ListType{ElemType: types.ObjectType{AttrTypes: PolicyTopologyRuleBindingAttrTypes}},
	"topology_id": types.StringType,
}

// PolicyTopologyResourceSchemaAttributes contains the schema attributes for PolicyTopologyModel
var PolicyTopologyResourceSchemaAttributes = map[string]schema.Attribute{
	"rules": schema.ListNestedAttribute{
		NestedObject: schema.NestedAttributeObject{
			Attributes: PolicyTopologyRuleBindingResourceSchemaAttributes,
		},
		Optional:            true,
		MarkdownDescription: "Ordered list of destination bindings, one per __TopologyRulePreset__ in the referenced __Topology__. Each binding maps a preset _name_ to the destination to use when that preset matches.",
	},
	"topology_id": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "The resource identifier.",
	},
}

// ExpandPolicyTopology converts a Terraform Object to SDK type
func ExpandPolicyTopology(ctx context.Context, o types.Object, diags *diag.Diagnostics) *uddidtc.PolicyTopology {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m PolicyTopologyModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

// Expand converts the Terraform model to SDK type
func (m *PolicyTopologyModel) Expand(ctx context.Context, diags *diag.Diagnostics) *uddidtc.PolicyTopology {
	if m == nil {
		return nil
	}
	to := &uddidtc.PolicyTopology{
		Rules:      flex.ExpandFrameworkListNestedBlock(ctx, m.Rules, diags, ExpandPolicyTopologyRuleBinding),
		TopologyId: flex.ExpandString(m.TopologyId),
	}
	return to
}

// FlattenPolicyTopology converts an SDK type to Terraform Object
func FlattenPolicyTopology(ctx context.Context, from *uddidtc.PolicyTopology, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(PolicyTopologyAttrTypes)
	}
	m := &PolicyTopologyModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, PolicyTopologyAttrTypes, m)
	diags.Append(d...)
	return t
}

// Flatten populates the Terraform model from SDK type
func (m *PolicyTopologyModel) Flatten(ctx context.Context, from *uddidtc.PolicyTopology, diags *diag.Diagnostics) {
	if from == nil || m == nil {
		return
	}
	m.Rules = flex.FlattenFrameworkListNestedBlock(ctx, from.Rules, PolicyTopologyRuleBindingAttrTypes, diags, FlattenPolicyTopologyRuleBinding)
	m.TopologyId = flex.FlattenString(from.TopologyId)
}
