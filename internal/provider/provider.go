package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// GoPackagerProviderModel describes the provider data model.
type GoPackagerProviderModel struct{}

// GoPackagerProvider defines the provider implementation.
type GoPackagerProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// New creates a new provider instance.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &GoPackagerProvider{
			version: version,
		}
	}
}

// Sets the provider metadata.
func (g *GoPackagerProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "gopackager"
	resp.Version = g.version
}

// Sets the provider schema.
func (g *GoPackagerProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	description := `Provides a resource to compile GoLang source code into a binary executable. This resource requires GoLang to be installed on the system. For more details see [gopackager_compile](https://registry.terraform.io/providers/StevenCyb/gopackager/latest/docs/data-sources/compile).`
	resp.Schema = schema.Schema{
		MarkdownDescription: description,
		Description:         description,
	}
}

// Configures the provider.
func (g *GoPackagerProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data GoPackagerProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Resources returns the provider resources.
// Currently there aren't any resources for this provider.
func (g *GoPackagerProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

// DataSources returns the provider data sources.
func (g *GoPackagerProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCompilerDataSource,
	}
}
