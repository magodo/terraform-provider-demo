package demo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type provider struct{}

func New() tfsdk.Provider {
	return &provider{}
}

func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{}, nil
}

func (p *provider) Configure(_ context.Context, _ tfsdk.ConfigureProviderRequest, _ *tfsdk.ConfigureProviderResponse) {
	return
}

func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{}, nil
}

// GetDataSources returns a map of the data source types this provider
// supports.
func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{}, nil
}
