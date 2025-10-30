package provider

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	fwpath "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/stevencyb/gopackager/internal/compiler"
	"github.com/stevencyb/gopackager/internal/git"
	"github.com/stevencyb/gopackager/internal/hasher"
	"github.com/stevencyb/gopackager/internal/packager"
)

// This is the global compiler instance.
// This instance is replaced by the mock instance during tests.
var globalCompiler compiler.CompilerI = compiler.New()

// This is the global ZIP packager instance.
// This instance is replaced by the mock instance during tests.
var globalZIPPackager packager.ZIPI = packager.New()

// This is the global hasher instance.
// This instance is replaced by the mock instance during tests.
var globalHasher hasher.HasherI = hasher.New()

// CompileDataSourceModel is the model for the compile data source.
type CompileDataSourceModel struct {
	// Input
	Source      types.String `tfsdk:"source"`
	Destination types.String `tfsdk:"destination"`
	GOOS        types.String `tfsdk:"goos"`
	GOARCH      types.String `tfsdk:"goarch"`
	// Optional
	ZIP            types.Bool   `tfsdk:"zip"`
	ZIPResources   types.Map    `tfsdk:"zip_resources"`
	GitTrigger     types.Bool   `tfsdk:"git_trigger"`
	GitTriggerPath types.String `tfsdk:"git_trigger_path"`
	// Output
	OutputPath         types.String `tfsdk:"output_path"`
	OutputMD5          types.String `tfsdk:"output_md5"`
	OutputSHA1         types.String `tfsdk:"output_sha1"`
	OutputSHA256       types.String `tfsdk:"output_sha256"`
	OutputSHA512       types.String `tfsdk:"output_sha512"`
	OutputSHA256Base64 types.String `tfsdk:"output_sha256_base64"`
	OutputSHA512Base64 types.String `tfsdk:"output_sha512_base64"`
	OutputGITHash      types.String `tfsdk:"output_git_hash"`
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
	description := `Compiles GoLang source code into a binary executable and optionally creates a ZIP with additional files.` +
		` This resource requires GoLang to be installed on the system.` +
		` The resource will automatically download the required dependencies and compile the source code.`

	resp.Schema = schema.Schema{
		MarkdownDescription: description,
		Description:         description,

		Attributes: map[string]schema.Attribute{
			// Required input
			"source": schema.StringAttribute{
				MarkdownDescription: "Path to the main file.",
				Required:            true,
			},
			"destination": schema.StringAttribute{
				MarkdownDescription: "Path for the compiled binary (or random UUID).",
				Required:            true,
			},
			"goos": schema.StringAttribute{
				MarkdownDescription: "GOOS for the compiled binary.",
				Required:            true,
			},
			"goarch": schema.StringAttribute{
				MarkdownDescription: "GOARCH for the compiled binary.",
				Required:            true,
			},
			// Output input
			"zip": schema.BoolAttribute{
				MarkdownDescription: "Zip the compiled binary and additional resources.",
				Optional:            true,
			},
			"zip_resources": schema.MapAttribute{
				MarkdownDescription: "Additional resources to include in the zip file. The binary is automatically included an copied to the root of the zip file.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"git_trigger": schema.BoolAttribute{
				MarkdownDescription: "Enable git trigger mode to only rebuild when files in git_trigger_path have changed since last compilation.",
				Optional:            true,
			},
			"git_trigger_path": schema.StringAttribute{
				MarkdownDescription: "Path to watch for changes when git_trigger is enabled. Defaults to current directory if not specified.",
				Optional:            true,
			},
			// Output
			"output_path": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Output path for the compiled binary or compressed ZIP file.",
			},
			"output_md5": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "MD5 hash of the compiled binary or compressed ZIP file.",
			},
			"output_sha1": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "SHA1 hash of the compiled binary or compressed ZIP file.",
			},
			"output_sha256": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "SHA256 hash of the compiled binary or compressed ZIP file.",
			},
			"output_sha512": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "SHA512 hash of the compiled binary or compressed ZIP file.",
			},
			"output_sha256_base64": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Base64 encoded SHA256 hash of the compiled binary or compressed ZIP file.",
			},
			"output_sha512_base64": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Base64 encoded SHA512 hash of the compiled binary or compressed ZIP file.",
			},
			"output_git_hash": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Last commit hash of the repository that changed `*.go`,`go.mod` or `go.sum` files.",
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
		return
	}

	// Check git trigger conditions if enabled
	if !data.GitTrigger.IsNull() && !data.GitTrigger.IsUnknown() && data.GitTrigger.ValueBool() {
		tflog.Trace(ctx, "Git trigger enabled, checking for changes")

		// Determine the path to monitor for changes
		triggerPath := "."
		if !data.GitTriggerPath.IsNull() && !data.GitTriggerPath.IsUnknown() {
			triggerPath = data.GitTriggerPath.ValueString()
		}

		// Determine expected output path for checking previous compilation
		expectedOutputPath := data.Destination.ValueString()
		if !data.ZIP.IsNull() && !data.ZIP.IsUnknown() && data.ZIP.ValueBool() {
			expectedOutputPath += ".zip"
		}

		// Check if we have a previous compilation commit recorded
		lastCompilationCommit, err := git.GetLastCompilationCommit(expectedOutputPath)
		if err == nil && lastCompilationCommit != "" {
			// Check if there have been changes in the trigger path since last compilation
			hasChanges, gitErr := git.HasChangedSinceCommit(triggerPath, lastCompilationCommit)
			if gitErr == nil {
				if !hasChanges {
					// No changes detected, check if output file still exists
					if content, readErr := globalHasher.ReadFile(expectedOutputPath); readErr == nil {
						tflog.Trace(ctx, "No changes detected since last compilation, using existing output")

						// Calculate hashes for existing file and return
						combinedHashes := globalHasher.CombinedHash(content)
						data.OutputPath = types.StringValue(expectedOutputPath)
						data.OutputMD5 = types.StringValue(combinedHashes.MD5)
						data.OutputSHA1 = types.StringValue(combinedHashes.SHA1)
						data.OutputSHA256 = types.StringValue(combinedHashes.SHA256)
						data.OutputSHA512 = types.StringValue(combinedHashes.SHA512)
						data.OutputSHA256Base64 = types.StringValue(combinedHashes.SHA256Base64)
						data.OutputSHA512Base64 = types.StringValue(combinedHashes.SHA512Base64)
						data.OutputGITHash = types.StringValue(lastCompilationCommit)

						resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
						return
					} else {
						tflog.Trace(ctx, "Output file missing despite no git changes, will recompile")
					}
				} else {
					tflog.Trace(ctx, "Changes detected in trigger path since last compilation, will recompile")
				}
			} else {
				tflog.Warn(ctx, "Failed to check git changes, proceeding with compilation", map[string]interface{}{
					"error": gitErr.Error(),
				})
			}
		} else {
			tflog.Trace(ctx, "No previous compilation commit found, will compile")
		}
	}

	tflog.Trace(ctx, "Compiling GoLang source code")

	outputPath, err := globalCompiler.Compile(*conf)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to compile binary.",
			"Compiling go code failed due '"+err.Error()+"'.",
		)

		return
	}

	if !data.ZIP.IsNull() && !data.ZIP.IsUnknown() && data.ZIP.ValueBool() {
		additionalFiles := map[string]string{}
		if !data.ZIPResources.IsNull() && !data.ZIPResources.IsUnknown() {
			if diag := data.ZIPResources.ElementsAs(ctx, &additionalFiles, false); diag.HasError() {
				resp.Diagnostics.Append(diag...)
			}
		}

		tflog.Trace(ctx, fmt.Sprintf("Zipping compiled binary with %d additional files", len(additionalFiles)))

		additionalFiles[outputPath] = filepath.Base(outputPath)
		outputPath += ".zip"

		if err = globalZIPPackager.Zip(outputPath, additionalFiles); err != nil {
			resp.Diagnostics.AddError(
				"Unable to create ZIP file.",
				"ZIP failed with: '"+err.Error()+"'.",
			)
		}
	}

	tflog.Trace(ctx, "Compute hashes for created file")
	content, err := globalHasher.ReadFile(outputPath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read compiled binary or ZIP file.",
			"Reading compiled binary or ZIP file failed with: '"+err.Error()+"'.",
		)

		return
	}

	combinedHashes := globalHasher.CombinedHash(content)

	data.OutputPath = types.StringValue(outputPath)
	data.OutputMD5 = types.StringValue(combinedHashes.MD5)
	data.OutputSHA1 = types.StringValue(combinedHashes.SHA1)
	data.OutputSHA256 = types.StringValue(combinedHashes.SHA256)
	data.OutputSHA512 = types.StringValue(combinedHashes.SHA512)
	data.OutputSHA256Base64 = types.StringValue(combinedHashes.SHA256Base64)
	data.OutputSHA512Base64 = types.StringValue(combinedHashes.SHA512Base64)
	// Determine which git hash to use and handle git trigger post-compilation
	var commitHash string
	var gitErr error

	if !data.GitTrigger.IsNull() && !data.GitTrigger.IsUnknown() && data.GitTrigger.ValueBool() {
		// Determine the path to monitor for changes
		triggerPath := "."
		if !data.GitTriggerPath.IsNull() && !data.GitTriggerPath.IsUnknown() {
			triggerPath = data.GitTriggerPath.ValueString()
		}

		// Get the last compilation commit to determine what triggered this build
		lastCompilationCommit, err := git.GetLastCompilationCommit(outputPath)
		if err == nil {
			// Get the commit that actually triggered this compilation
			commitHash, gitErr = git.GetTriggeringCommitHash(triggerPath, lastCompilationCommit)
			tflog.Trace(ctx, fmt.Sprintf("Triggering commit for path %s: %s (since %s)", triggerPath, commitHash, lastCompilationCommit))
		} else {
			// Fallback to latest commit for the path if no previous compilation
			commitHash, gitErr = git.LastCommitHashForPath(triggerPath)
			tflog.Trace(ctx, fmt.Sprintf("No previous compilation, using latest commit for path %s: %s", triggerPath, commitHash))
		}

		if gitErr == nil {
			// Get the current HEAD commit hash for tracking future changes
			currentHeadCommit, headErr := git.LastCommitHash()
			if headErr == nil {
				// Save the current HEAD commit as the last compilation commit for future comparisons
				saveErr := git.SaveLastCompilationCommit(outputPath, currentHeadCommit)
				if saveErr != nil {
					tflog.Warn(ctx, "Failed to save compilation commit for git trigger", map[string]interface{}{
						"error": saveErr.Error(),
					})
				} else {
					tflog.Trace(ctx, fmt.Sprintf("Saved HEAD commit: %s for output: %s", currentHeadCommit, outputPath))
				}
			}
		}
	} else {
		// Use the original git hash calculation for non-trigger mode
		commitHash, gitErr = git.LastCommitHash()
	}

	if gitErr != nil {
		data.OutputGITHash = types.StringValue("unknown")
	} else {
		data.OutputGITHash = types.StringValue(commitHash)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// ConfigValidators returns the config validators for this data source.
func (c *CompileDataSource) ConfigValidators(context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.RequiredTogether(
			fwpath.MatchRoot("source"),
			fwpath.MatchRoot("destination"),
			fwpath.MatchRoot("goos"),
			fwpath.MatchRoot("goarch"),
		),
	}
}
