package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/stevencyb/terraform-provider-gopackager/internal/compiler"
)

// This is the global compiler instance.
// This instance is replaced by the mock instance during tests.
var globalCompiler compiler.CompilerI = compiler.New()

// CompileDataSourceModel is the model for the compile data source.
type CompileDataSourceModel struct {
	// Input
	Source      types.String `tfsdk:"source"`
	Destination types.String `tfsdk:"destination"`
	GOOS        types.String `tfsdk:"goos"`
	GOARCH      types.String `tfsdk:"goarch"`
	// Output
	BinaryLocation types.String `tfsdk:"binary_location"`
	BinaryHash     types.String `tfsdk:"binary_hash"`
}

// CompileDataSource is the data source for the compile resource.
type CompileDataSource struct{}

// New creates a new data source instance.
func NewCompilerDataSource() datasource.DataSource {
	return &CompileDataSource{}
}

// Sets the provider metadata.
func (c *CompileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_compile"
}

// Sets the provider schema.
func (c *CompileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	description := `Compiles GoLang source code into a binary executable. This resource requires GoLang to be installed on the system.`
	resp.Schema = schema.Schema{
		MarkdownDescription: description,
		Description:         description,

		Attributes: map[string]schema.Attribute{
			// Input
			"source": schema.StringAttribute{
				MarkdownDescription: "Path to the main file",
				Required:            true,
			},
			"destination": schema.StringAttribute{
				MarkdownDescription: "Path for the compiled binary (or random UUID)",
				Required:            true,
			},
			"goos": schema.StringAttribute{
				MarkdownDescription: "GOOS for the compiled binary",
				Required:            true,
			},
			"goarch": schema.StringAttribute{
				MarkdownDescription: "GOARCH for the compiled binary",
				Required:            true,
			},
			// Output
			"binary_location": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Binary location of compiled file.",
			},
			"binary_hash": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Binary hash for compiled file.",
			},
		},
	}
}

// Configures the provider.
func (c *CompileDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
}

// Read event for this data source.
func (c *CompileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CompileDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Checking configuration")

	conf := compiler.NewConfig().
		Source(data.Source.ValueString()).
		Destination(data.Destination.ValueString()).
		GOOS(data.GOOS.ValueString()).
		GOARCH(data.GOARCH.ValueString())
	if err := conf.Verify(); err != nil {
		resp.Diagnostics.AddError(
			"Invalid configuration.",
			"Expected configuration to be valid, but got '"+err.Error()+"'.",
		)
	}

	tflog.Trace(ctx, "Compiling GoLang source code")

	binaryLocation, hash, err := globalCompiler.Compile(*conf)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource",
			"Compiling go code failed due '"+err.Error()+"'.",
		)

		return
	}

	data.BinaryHash = types.StringValue(hash)
	data.BinaryLocation = types.StringValue(binaryLocation)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ValidateConfig validates the configuration.
func (c *CompileDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data CompileDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Source.IsNull() || data.Source.IsUnknown() {
		resp.Diagnostics.AddError(
			"Missing main file attribute.",
			"Expected main file to point to a GoLang main file.",
		)

		return
	}
	if data.GOOS.IsNull() || data.GOOS.IsUnknown() {
		resp.Diagnostics.AddError(
			"Missing GOOS attribute.",
			"Expected GOOS to be set to a supported value.",
		)

		return
	}
	if data.GOARCH.IsNull() || data.GOARCH.IsUnknown() {
		resp.Diagnostics.AddError(
			"Missing GOARCH attribute.",
			"Expected GOARCH to be set to a supported value.",
		)
	}

	conf := compiler.NewConfig().
		Source(data.Source.ValueString()).
		Destination(data.Destination.ValueString()).
		GOOS(data.GOOS.ValueString()).
		GOARCH(data.GOARCH.ValueString())
	if err := conf.Verify(); err != nil {
		resp.Diagnostics.AddError(
			"Invalid configuration.",
			"Expected configuration to be valid, but got '"+err.Error()+"'.",
		)
	}
}
