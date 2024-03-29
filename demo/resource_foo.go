package demo

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/magodo/terraform-provider-demo/client"
)

type resourceFoo struct {
	p *Provider
}

type fooData struct {
	ID              types.String  `tfsdk:"id"`
	String          types.String  `tfsdk:"string"`
	Int64           types.Int64   `tfsdk:"int64"`
	Float64         types.Float64 `tfsdk:"float64"`
	Number          types.Number  `tfsdk:"number"`
	Bool            types.Bool    `tfsdk:"bool"`
	ListNestedBlock types.List    `tfsdk:"list_nested_block"`
	SetNestedBlock  types.Set     `tfsdk:"set_nested_block"`
}

type nestedData struct {
	Name types.String `tfsdk:"name"`
	Age  types.Int64  `tfsdk:"age"`
}

var _ resource.Resource = resourceFoo{}

// Metadata implements resource.Resource.
func (resourceFoo) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_foo"
}

// Schema implements resource.Resource.
func (resourceFoo) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Resource Foo",
		MarkdownDescription: "Resource Foo",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"string": schema.StringAttribute{
				Optional: true,
			},
			"int64": schema.Int64Attribute{
				Optional: true,
			},
			"float64": schema.Float64Attribute{
				Optional: true,
			},
			"number": schema.NumberAttribute{
				Optional: true,
			},
			"bool": schema.BoolAttribute{
				Optional: true,
			},
		},
		Blocks: map[string]schema.Block{
			"list_nested_block": schema.ListNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Optional: true,
						},
						"age": schema.Int64Attribute{
							Optional: true,
						},
					},
				},
			},
			"set_nested_block": schema.SetNestedBlock{
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Optional: true,
						},
						"age": schema.Int64Attribute{
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func (r *resourceFoo) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	provider, ok := req.ProviderData.(*Provider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("got: %T.", req.ProviderData),
		)
		return
	}
	r.p = provider
}

// Create is called when the provider must create a new resource. Config
// and planned state values should be read from the
// CreateResourceRequest and new state values set on the
// CreateResourceResponse.
func (r resourceFoo) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan fooData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Expand
	m := map[string]interface{}{}
	if !plan.String.IsNull() {
		m["string"] = plan.String.ValueString()
	}
	if !plan.Int64.IsNull() {
		m["int64"] = plan.Int64.ValueInt64()
	}
	if !plan.Float64.IsNull() {
		m["float64"] = plan.Float64.ValueFloat64()
	}
	if !plan.Number.IsNull() {
		m["number"], _ = plan.Number.ValueBigFloat().Float64()
	}
	if !plan.Bool.IsNull() {
		m["bool"] = plan.Bool.ValueBool()
	}
	if !plan.ListNestedBlock.IsNull() {
		var blks []nestedData
		diags := plan.ListNestedBlock.ElementsAs(ctx, &blks, false)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}
		m["list_nested_block"] = expandNestedObject(blks)
	}
	if !plan.SetNestedBlock.IsNull() {
		var blks []nestedData
		diags := plan.SetNestedBlock.ElementsAs(ctx, &blks, false)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}
		m["set_nested_block"] = expandNestedObject(blks)
	}
	b, err := json.Marshal(m)
	if err != nil {
		resp.Diagnostics.AddError(
			"Creation failure",
			fmt.Sprintf("Failed to JSON encode the request: %v", err),
		)
		return
	}
	id, err := r.p.client.Create(b)
	if err != nil {
		resp.Diagnostics.AddError(
			"Creation failure",
			fmt.Sprintf("Sending create request: %v", err),
		)
		return
	}
	diags = resp.State.Set(ctx,
		fooData{
			ID:              types.StringValue(id),
			String:          types.StringNull(),
			Int64:           types.Int64Null(),
			Float64:         types.Float64Null(),
			Number:          types.NumberNull(),
			Bool:            types.BoolNull(),
			ListNestedBlock: types.ListNull(types.ObjectType{AttrTypes: map[string]attr.Type{"name": types.StringType, "age": types.Int64Type}}),
			SetNestedBlock:  types.SetNull(types.ObjectType{AttrTypes: map[string]attr.Type{"name": types.StringType, "age": types.Int64Type}}),
		},
	)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	rreq := resource.ReadRequest{
		State:        resp.State,
		ProviderMeta: req.ProviderMeta,
	}
	rresp := resource.ReadResponse{
		State:       resp.State,
		Diagnostics: resp.Diagnostics,
	}
	r.Read(ctx, rreq, &rresp)

	*resp = resource.CreateResponse{
		State:       rresp.State,
		Diagnostics: rresp.Diagnostics,
	}
}

// Read is called when the provider must read resource values in order
// to update state. Planned state values should be read from the
// ReadResourceRequest and new state values set on the
// ReadResourceResponse.
func (r resourceFoo) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state fooData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	b, err := r.p.client.Read(state.ID.ValueString())
	if err != nil {
		if err == client.ErrNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Read failure",
			fmt.Sprintf("Sending read request: %v", err),
		)
		return
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		resp.Diagnostics.AddError(
			"Read failure",
			fmt.Sprintf("Failed to JSON decode the response: %v", err),
		)
		return
	}

	// Flatten
	if v, ok := m["string"]; ok {
		state.String = types.StringValue(v.(string))
	}
	if v, ok := m["int64"]; ok {
		state.Int64 = types.Int64Value(int64(v.(float64)))
	}
	if v, ok := m["float64"]; ok {
		state.Float64 = types.Float64Value(v.(float64))
	}
	if v, ok := m["number"]; ok {
		state.Number = types.NumberValue(big.NewFloat(v.(float64)))
	}
	if v, ok := m["bool"]; ok {
		state.Bool = types.BoolValue(v.(bool))
	}
	if v, ok := m["list_nested_block"]; ok {
		state.ListNestedBlock = types.ListValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{"name": types.StringType, "age": types.Int64Type}}, flattenNestedObject(v.([]interface{})))
	}
	if v, ok := m["set_nested_block"]; ok {
		state.SetNestedBlock = types.SetValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{"name": types.StringType, "age": types.Int64Type}}, flattenNestedObject(v.([]interface{})))
	}
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
}

// Update is called to update the state of the resource. Config, planned
// state, and prior state values should be read from the
// UpdateResourceRequest and new state values set on the
// UpdateResourceResponse.
func (r resourceFoo) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan fooData
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Expand
	m := map[string]interface{}{}
	if !plan.String.IsNull() {
		m["string"] = plan.String.ValueString()
	}
	if !plan.Int64.IsNull() {
		m["int64"] = plan.Int64.ValueInt64()
	}
	if !plan.Float64.IsNull() {
		m["float64"] = plan.Float64.ValueFloat64()
	}
	if !plan.Number.IsNull() {
		m["number"], _ = plan.Number.ValueBigFloat().Float64()
	}
	if !plan.Bool.IsNull() {
		m["bool"] = plan.Bool.ValueBool()
	}
	if !plan.ListNestedBlock.IsNull() {
		var blks []nestedData
		diags := plan.ListNestedBlock.ElementsAs(ctx, &blks, false)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}
		m["list_nested_block"] = expandNestedObject(blks)
	}
	if !plan.SetNestedBlock.IsNull() {
		var blks []nestedData
		diags := plan.SetNestedBlock.ElementsAs(ctx, &blks, false)
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}
		m["set_nested_block"] = expandNestedObject(blks)
	}
	b, err := json.Marshal(m)
	if err != nil {
		resp.Diagnostics.AddError(
			"Update failure",
			fmt.Sprintf("Failed to JSON encode the request: %v", err),
		)
		return
	}

	var state fooData
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if err := r.p.client.Update(state.ID.ValueString(), b); err != nil {
		resp.Diagnostics.AddError(
			"Update failure",
			fmt.Sprintf("Sending update request: %v", err),
		)
		return
	}

	rreq := resource.ReadRequest{
		State:        resp.State,
		ProviderMeta: req.ProviderMeta,
	}
	rresp := resource.ReadResponse{
		State:       resp.State,
		Diagnostics: resp.Diagnostics,
	}
	r.Read(ctx, rreq, &rresp)

	*resp = resource.UpdateResponse{
		State:       rresp.State,
		Diagnostics: rresp.Diagnostics,
	}
}

// Delete is called when the provider must delete the resource. Config
// values may be read from the DeleteResourceRequest.
func (r resourceFoo) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state fooData
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if err := r.p.client.Delete(state.ID.ValueString()); err != nil {
		if err == client.ErrNotFound {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Delete failure",
			fmt.Sprintf("Sending delete request: %v", err),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState is called when the provider must import the resource.
//
// If import is not supported, it is recommended to use the
// ResourceImportStateNotImplemented() call in this method.
//
// If setting an attribute with the import identifier, it is recommended
// to use the ResourceImportStatePassthroughID() call in this method.
func (r resourceFoo) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func expandNestedObject(l []nestedData) []interface{} {
	var output []interface{}

	for _, d := range l {
		m := map[string]interface{}{}
		if !d.Name.IsNull() {
			m["name"] = d.Name.ValueString()
		}
		if !d.Age.IsNull() {
			m["age"] = d.Age.ValueInt64()
		}
		output = append(output, m)
	}
	return output
}

func flattenNestedObject(l []interface{}) []attr.Value {
	var elements []attr.Value

	for _, v := range l {
		m := v.(map[string]interface{})

		name := types.StringNull()
		if v, ok := m["name"]; ok {
			name = types.StringValue(v.(string))
		}
		age := types.Int64Null()
		if v, ok := m["age"]; ok {
			age = types.Int64Value(int64(v.(float64)))
		}

		obj, diags := types.ObjectValue(
			map[string]attr.Type{
				"name": types.StringType,
				"age":  types.Int64Type,
			},
			map[string]attr.Value{
				"name": name,
				"age":  age,
			},
		)
		if diags.HasError() {
			panic(fmt.Sprintf("%v", diags.Errors()))
		}

		elements = append(elements, obj)
	}

	return elements
}
