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

// PolicyTopologyRuleBindingModel is the Terraform model for PolicyTopologyRuleBinding
type PolicyTopologyRuleBindingModel struct {
	Code        types.String `tfsdk:"code"`
	Destination types.String `tfsdk:"destination"`
	Name        types.String `tfsdk:"name"`
	PoolId      types.String `tfsdk:"pool_id"`
}

// PolicyTopologyRuleBindingAttrTypes contains the attribute types for PolicyTopologyRuleBindingModel
var PolicyTopologyRuleBindingAttrTypes = map[string]attr.Type{
	"code":        types.StringType,
	"destination": types.StringType,
	"name":        types.StringType,
	"pool_id":     types.StringType,
}

// PolicyTopologyRuleBindingResourceSchemaAttributes contains the schema attributes for PolicyTopologyRuleBindingModel
var PolicyTopologyRuleBindingResourceSchemaAttributes = map[string]schema.Attribute{
	"code": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "Optional. DNS code to return if the referenced preset matches. Must be set if _destination_ is _code_.  Allowed values: - nodata - nxdomain  Defaults to _nodata_.",
	},
	"destination": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "Destination of the matched __TopologyRulePreset__.  Allowed values: - code - pool  Defaults to _code_.",
	},
	"name": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "Required. Name of the __TopologyRulePreset__ from the referenced __Topology__.",
	},
	"pool_id": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "The resource identifier.",
	},
}

// ExpandPolicyTopologyRuleBinding converts a Terraform Object to SDK type
func ExpandPolicyTopologyRuleBinding(ctx context.Context, o types.Object, diags *diag.Diagnostics) *uddidtc.PolicyTopologyRuleBinding {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m PolicyTopologyRuleBindingModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

// Expand converts the Terraform model to SDK type
func (m *PolicyTopologyRuleBindingModel) Expand(ctx context.Context, diags *diag.Diagnostics) *uddidtc.PolicyTopologyRuleBinding {
	if m == nil {
		return nil
	}
	to := &uddidtc.PolicyTopologyRuleBinding{
		Code:        flex.ExpandStringPointer(m.Code),
		Destination: flex.ExpandStringPointer(m.Destination),
		Name:        flex.ExpandString(m.Name),
		PoolId:      flex.ExpandStringPointer(m.PoolId),
	}
	return to
}

// FlattenPolicyTopologyRuleBinding converts an SDK type to Terraform Object
func FlattenPolicyTopologyRuleBinding(ctx context.Context, from *uddidtc.PolicyTopologyRuleBinding, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(PolicyTopologyRuleBindingAttrTypes)
	}
	m := &PolicyTopologyRuleBindingModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, PolicyTopologyRuleBindingAttrTypes, m)
	diags.Append(d...)
	return t
}

// Flatten populates the Terraform model from SDK type
func (m *PolicyTopologyRuleBindingModel) Flatten(ctx context.Context, from *uddidtc.PolicyTopologyRuleBinding, diags *diag.Diagnostics) {
	if from == nil || m == nil {
		return
	}
	m.Code = flex.FlattenStringPointer(from.Code)
	m.Destination = flex.FlattenStringPointer(from.Destination)
	m.Name = flex.FlattenString(from.Name)
	m.PoolId = flex.FlattenStringPointer(from.PoolId)
}
