package demo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/magodo/terraform-provider-demo/client"
)

type resourceFooType struct{}

// GetSchema returns the schema for this resource.
func (r resourceFooType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:         "Resource Foo",
		MarkdownDescription: "Resource Foo",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.UseStateForUnknown(),
				},
			},
			"name": {
				Type:                types.StringType,
				Description:         "name",
				MarkdownDescription: "name",
				Required:            true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
				},
			},
			"age": {
				Type:                types.Int64Type,
				Description:         "age",
				MarkdownDescription: "age",
				Optional:            true,
			},
			"job": {
				Type:                types.StringType,
				Description:         "job",
				MarkdownDescription: "job",
				Optional:            true,
			},
		},
	}, nil
}

// NewResource instantiates a new Resource of this ResourceType.
func (r resourceFooType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return resourceFoo{p: *p.(*provider)}, nil
}

type resourceFoo struct {
	p provider
}

var _ tfsdk.Resource = resourceFoo{}

type foo struct {
	ID   types.String `tfsdk:"id"   json:"id,omitempty"`
	Name types.String `tfsdk:"name" json:"name,omitempty"`
	Age  types.Int64  `tfsdk:"age"  json:"age,omitempty"`
	Job  types.String `tfsdk:"job"  json:"job,omitempty"`
}

// Create is called when the provider must create a new resource. Config
// and planned state values should be read from the
// CreateResourceRequest and new state values set on the
// CreateResourceResponse.
func (r resourceFoo) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var plan foo
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Expand
	m := map[string]interface{}{
		"name": plan.Name.Value,
	}
	if !plan.Age.Null {
		m["age"] = plan.Age.Value
	}
	if !plan.Job.Null {
		m["job"] = plan.Job.Value
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
		foo{
			ID:   types.String{Value: id},
			Name: types.String{Null: true},
			Age:  types.Int64{Null: true},
			Job:  types.String{Null: true},
		},
	)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	rreq := tfsdk.ReadResourceRequest{
		State:        resp.State,
		ProviderMeta: req.ProviderMeta,
	}
	rresp := tfsdk.ReadResourceResponse{
		State:       resp.State,
		Diagnostics: resp.Diagnostics,
	}
	r.Read(ctx, rreq, &rresp)

	*resp = tfsdk.CreateResourceResponse{
		State:       rresp.State,
		Diagnostics: rresp.Diagnostics,
	}
}

// Read is called when the provider must read resource values in order
// to update state. Planned state values should be read from the
// ReadResourceRequest and new state values set on the
// ReadResourceResponse.
func (r resourceFoo) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state foo
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	b, err := r.p.client.Read(state.ID.Value)
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

	tflog.Debug(ctx, "read content", m)
	tflog.Debug(ctx, "state", state)

	// Flatten
	if name, ok := m["name"]; ok {
		state.Name = types.String{Value: name.(string)}
	}
	if age, ok := m["age"]; ok {
		state.Age = types.Int64{Value: int64(age.(float64))}
	}
	if job, ok := m["job"]; ok {
		state.Job = types.String{Value: job.(string)}
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
func (r resourceFoo) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var plan foo
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	// Expand
	m := map[string]interface{}{
		"name": plan.Name.Value,
	}
	if !plan.Age.Null {
		m["age"] = plan.Age.Value
	}
	if !plan.Job.Null {
		m["job"] = plan.Job.Value
	}
	b, err := json.Marshal(m)
	if err != nil {
		resp.Diagnostics.AddError(
			"Update failure",
			fmt.Sprintf("Failed to JSON encode the request: %v", err),
		)
		return
	}

	var state foo
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if err := r.p.client.Update(state.ID.Value, b); err != nil {
		resp.Diagnostics.AddError(
			"Update failure",
			fmt.Sprintf("Sending update request: %v", err),
		)
		return
	}

	rreq := tfsdk.ReadResourceRequest{
		State:        resp.State,
		ProviderMeta: req.ProviderMeta,
	}
	rresp := tfsdk.ReadResourceResponse{
		State:       resp.State,
		Diagnostics: resp.Diagnostics,
	}
	r.Read(ctx, rreq, &rresp)

	*resp = tfsdk.UpdateResourceResponse{
		State:       rresp.State,
		Diagnostics: rresp.Diagnostics,
	}
}

// Delete is called when the provider must delete the resource. Config
// values may be read from the DeleteResourceRequest.
func (r resourceFoo) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state foo
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if err := r.p.client.Delete(state.ID.Value); err != nil {
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
func (r resourceFoo) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}
