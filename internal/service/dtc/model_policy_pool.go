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

// PolicyPoolModel is the Terraform model for PolicyPool
type PolicyPoolModel struct {
	Name   types.String `tfsdk:"name"`
	PoolId types.String `tfsdk:"pool_id"`
	Weight types.Int64  `tfsdk:"weight"`
}

// PolicyPoolAttrTypes contains the attribute types for PolicyPoolModel
var PolicyPoolAttrTypes = map[string]attr.Type{
	"name":    types.StringType,
	"pool_id": types.StringType,
	"weight":  types.Int64Type,
}

// PolicyPoolResourceSchemaAttributes contains the schema attributes for PolicyPoolModel
var PolicyPoolResourceSchemaAttributes = map[string]schema.Attribute{
	"name": schema.StringAttribute{
		Optional:            true,
		Computed:            true,
		MarkdownDescription: "Display name of __Pool__.",
	},
	"pool_id": schema.StringAttribute{
		Optional:            true,
		MarkdownDescription: "The resource identifier.",
	},
	"weight": schema.Int64Attribute{
		Optional:            true,
		MarkdownDescription: "Weight of __Pool__ to be used for load balancing. Unsigned integer, min 1; max 65535.",
	},
}

// ExpandPolicyPool converts a Terraform Object to SDK type
func ExpandPolicyPool(ctx context.Context, o types.Object, diags *diag.Diagnostics) *uddidtc.PolicyPool {
	if o.IsNull() || o.IsUnknown() {
		return nil
	}
	var m PolicyPoolModel
	diags.Append(o.As(ctx, &m, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return nil
	}
	return m.Expand(ctx, diags)
}

// Expand converts the Terraform model to SDK type
func (m *PolicyPoolModel) Expand(ctx context.Context, diags *diag.Diagnostics) *uddidtc.PolicyPool {
	if m == nil {
		return nil
	}
	to := &uddidtc.PolicyPool{
		Name:   flex.ExpandStringPointer(m.Name),
		PoolId: flex.ExpandString(m.PoolId),
		Weight: flex.ExpandInt64Pointer(m.Weight),
	}
	return to
}

// FlattenPolicyPool converts an SDK type to Terraform Object
func FlattenPolicyPool(ctx context.Context, from *uddidtc.PolicyPool, diags *diag.Diagnostics) types.Object {
	if from == nil {
		return types.ObjectNull(PolicyPoolAttrTypes)
	}
	m := &PolicyPoolModel{}
	m.Flatten(ctx, from, diags)
	t, d := types.ObjectValueFrom(ctx, PolicyPoolAttrTypes, m)
	diags.Append(d...)
	return t
}

// Flatten populates the Terraform model from SDK type
func (m *PolicyPoolModel) Flatten(ctx context.Context, from *uddidtc.PolicyPool, diags *diag.Diagnostics) {
	if from == nil || m == nil {
		return
	}
	m.Name = flex.FlattenStringPointer(from.Name)
	m.PoolId = flex.FlattenString(from.PoolId)
	m.Weight = flex.FlattenInt64Pointer(from.Weight)
}
