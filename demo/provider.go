package demo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/magodo/terraform-provider-demo/client"
)

type provider struct {
	client client.Client
}

type providerData struct {
	FileSystem *filesystemData `tfsdk:"filesystem"`
	JSONServer *jsonserverData `tfsdk:"jsonserver"`
}

type filesystemData struct {
	Workdir types.String `tfsdk:"workdir"`
}

type jsonserverData struct {
	URL types.String `tfsdk:"url"`
}

func New() tfsdk.Provider {
	return &provider{}
}

func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description:         "The schema of the magodo/terraform-provider-demo provider",
		MarkdownDescription: "The schema of the `magodo/terraform-provider-demo` provider",
		Attributes: map[string]tfsdk.Attribute{
			"filesystem": {
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"workdir": {
						Type:                types.StringType,
						Description:         "The directory to store the json files",
						MarkdownDescription: "The directory to store the json files",
						Required:            true,
					},
				}),
				Description:         "Using the filesystem as the backend service",
				MarkdownDescription: "Using the filesystem as the backend service",
				Optional:            true,
			},
			"jsonserver": {
				Attributes: tfsdk.SingleNestedAttributes(map[string]tfsdk.Attribute{
					"url": {
						Type:                types.StringType,
						Description:         "The URL to the json-server",
						MarkdownDescription: "The URL to the json-server",
						Required:            true,
					},
				}),
				Description:         "Using the json-server as the backend service",
				MarkdownDescription: "Using the [json-server](https://github.com/typicode/json-server) as the backend service",
				Optional:            true,
			},
		},
	}, nil
}

func (p *provider) ValidateConfig(ctx context.Context, req tfsdk.ValidateProviderConfigRequest, resp *tfsdk.ValidateProviderConfigResponse) {
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	if config.FileSystem == nil && config.JSONServer == nil {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			`None of "filesystem" and "jsonserver" is specified`,
		)
		return
	}
	if config.FileSystem != nil && config.JSONServer != nil {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			`Only one of "filesystem" and "jsonserver" can be specified`,
		)
		return
	}
	return
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	switch {
	case config.FileSystem != nil:
		client, err := client.NewFsClient(config.FileSystem.Workdir.Value)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to new filesystem client",
				err.Error(),
			)
		}
		p.client = client
	case config.JSONServer != nil:
		client, err := client.NewJSONServerClient(config.JSONServer.URL.Value)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to new jsonserver client",
				err.Error(),
			)
		}
		p.client = client
	}
}

func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{}, nil
}

// GetDataSources returns a map of the data source types this provider
// supports.
func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{}, nil
}
