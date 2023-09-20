package demo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/magodo/terraform-provider-demo/client"
)

type Provider struct {
	client client.Client
}

var _ provider.Provider = &Provider{}

type providerData struct {
	FileSystem types.Object `tfsdk:"filesystem"`
	JSONServer types.Object `tfsdk:"jsonserver"`
}

type filesystemData struct {
	Workdir types.String `tfsdk:"workdir"`
}

type jsonserverData struct {
	URL types.String `tfsdk:"url"`
}

func New() provider.Provider {
	return &Provider{}
}

func (*Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "demo"
}

func (p *Provider) Schema(_ context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "The schema of the magodo/terraform-provider-demo provider",
		MarkdownDescription: "The schema of the `magodo/terraform-provider-demo` provider",
		Attributes: map[string]schema.Attribute{
			"filesystem": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"workdir": schema.StringAttribute{
						Description:         "The directory to store the json files",
						MarkdownDescription: "The directory to store the json files",
						Required:            true,
					},
				},
				Description:         "Using the filesystem as the backend service",
				MarkdownDescription: "Using the filesystem as the backend service",
			},
			"jsonserver": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"url": schema.StringAttribute{
						Description:         "The URL to the json-server",
						MarkdownDescription: "The URL to the json-server",
						Required:            true,
					},
				},
				Description:         "Using the json-server as the backend service",
				MarkdownDescription: "Using the [json-server](https://github.com/typicode/json-server) as the backend service",
			},
		},
	}
}

func (p *Provider) ValidateConfig(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	if config.FileSystem.IsNull() && config.JSONServer.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			`None of "filesystem" and "jsonserver" is specified`,
		)
		return
	}
	if !config.FileSystem.IsNull() && !config.JSONServer.IsNull() {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			`Only one of "filesystem" and "jsonserver" can be specified`,
		)
		return
	}
	return
}

func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config providerData
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	switch {
	case !config.FileSystem.IsNull():
		var fs filesystemData
		diags := config.FileSystem.As(ctx, &fs, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}
		client, err := client.NewFsClient(fs.Workdir.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to new filesystem client",
				err.Error(),
			)
		}
		p.client = client
	case !config.JSONServer.IsNull():
		var jsonserver jsonserverData
		diags := config.FileSystem.As(ctx, &jsonserver, basetypes.ObjectAsOptions{})
		resp.Diagnostics.Append(diags...)
		if diags.HasError() {
			return
		}
		client, err := client.NewJSONServerClient(jsonserver.URL.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to new jsonserver client",
				err.Error(),
			)
		}
		p.client = client
	}

	resp.ResourceData = p
}

func (*Provider) DataSources(context.Context) []func() datasource.DataSource {
	return nil
}

func (*Provider) Resources(context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource {
			return &resourceFoo{}
		},
	}
}
