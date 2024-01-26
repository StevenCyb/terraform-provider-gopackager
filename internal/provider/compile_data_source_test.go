package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stevencyb/gopackager/internal/compiler"
	"github.com/stevencyb/gopackager/internal/packager"
	"github.com/stretchr/testify/assert"
)

func TestAccDataSourceFrameworkSatisfaction(t *testing.T) {
	t.Parallel()

	var _ datasource.DataSource = &CompileDataSource{}
}

func TestAccCompileDataSource(t *testing.T) {
	t.Parallel()

	mockCompiler := compiler.MockCompiler{}
	mockPackager := packager.MockZIP{}
	globalCompiler = &mockCompiler
	globalZIPPackager = &mockPackager
	testAccProtoV6ProviderFactories := map[string]func() (tfprotov6.ProviderServer, error){
		"gopackager": providerserver.NewProtocol6WithError(New("test")()),
	}

	additionalZIPResources := map[string]string{
		"../../LICENSE": "LICENSE",
		"../provider":   "a/provider",
	}
	additionalZIPResourcesGen, diag := types.MapValueFrom(context.Background(), types.StringType, additionalZIPResources)
	assert.False(t, diag.HasError())
	additionalZIPResources["windows_amd64_binary"] = "windows_amd64_binary"

	initialDataSource := CompileDataSourceModel{
		Source:         types.StringValue("provider.go"),
		Destination:    types.StringValue("linux_amd64_binary"),
		GOOS:           types.StringValue("linux"),
		GOARCH:         types.StringValue("amd64"),
		BinaryLocation: types.StringValue("linux_amd64_binary"),
		BinaryHash:     types.StringValue("123"),
	}
	firstUpdate := CompileDataSourceModel{
		Source:         types.StringValue("provider.go"),
		Destination:    types.StringValue("windows_amd64_binary"),
		GOOS:           types.StringValue("windows"),
		GOARCH:         types.StringValue("amd64"),
		BinaryLocation: types.StringValue("windows_amd64_binary"),
		BinaryHash:     types.StringValue("321"),
	}
	secondUpdate := CompileDataSourceModel{
		Source:         types.StringValue("provider.go"),
		Destination:    types.StringValue("windows_amd64_binary"),
		GOOS:           types.StringValue("windows"),
		GOARCH:         types.StringValue("amd64"),
		BinaryLocation: types.StringValue("windows_amd64_binary"),
		BinaryHash:     types.StringValue("404"),
	}
	thirdUpdate := CompileDataSourceModel{
		Source:         types.StringValue("provider.go"),
		Destination:    types.StringValue("windows_amd64_binary"),
		GOOS:           types.StringValue("windows"),
		GOARCH:         types.StringValue("amd64"),
		BinaryLocation: types.StringValue("windows_amd64_binary"),
		BinaryHash:     types.StringValue("404"),
		ZIP:            types.BoolValue(true),
		ZIPResources:   additionalZIPResourcesGen,
	}

	mockCompiler.On("Compile",
		*compiler.NewConfig().
			Source(initialDataSource.Source.ValueString()).
			Destination(initialDataSource.Destination.ValueString()).
			GOOS(initialDataSource.GOOS.ValueString()).
			GOARCH(initialDataSource.GOARCH.ValueString()),
	).Times(3).
		Return(initialDataSource.BinaryLocation.ValueString(), initialDataSource.BinaryHash.ValueString(), nil)
	mockCompiler.On("Compile",
		*compiler.NewConfig().
			Source(firstUpdate.Source.ValueString()).
			Destination(firstUpdate.Destination.ValueString()).
			GOOS(firstUpdate.GOOS.ValueString()).
			GOARCH(firstUpdate.GOARCH.ValueString()),
	).Times(3).
		Return(firstUpdate.BinaryLocation.ValueString(), firstUpdate.BinaryHash.ValueString(), nil)
	mockCompiler.On("Compile",
		*compiler.NewConfig().
			Source(secondUpdate.Source.ValueString()).
			Destination(secondUpdate.Destination.ValueString()).
			GOOS(secondUpdate.GOOS.ValueString()).
			GOARCH(secondUpdate.GOARCH.ValueString()),
	).Times(3).
		Return(secondUpdate.BinaryLocation.ValueString(), secondUpdate.BinaryHash.ValueString(), nil)
	mockCompiler.On("Compile",
		*compiler.NewConfig().
			Source(thirdUpdate.Source.ValueString()).
			Destination(thirdUpdate.Destination.ValueString()).
			GOOS(thirdUpdate.GOOS.ValueString()).
			GOARCH(thirdUpdate.GOARCH.ValueString()),
	).Times(3).
		Return(thirdUpdate.BinaryLocation.ValueString(), thirdUpdate.BinaryHash.ValueString(), nil)
	mockPackager.On("Zip", thirdUpdate.BinaryLocation.ValueString()+".zip", additionalZIPResources).Times(3).Return("ff", nil)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: compilerDataSourceFromModel(t, initialDataSource),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "source", initialDataSource.Source.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "destination", initialDataSource.Destination.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "goos", initialDataSource.GOOS.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "goarch", initialDataSource.GOARCH.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "binary_location", initialDataSource.BinaryLocation.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "binary_hash", initialDataSource.BinaryHash.ValueString()),
				),
			},
			// First update testing
			{
				Config: compilerDataSourceFromModel(t, firstUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "source", firstUpdate.Source.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "destination", firstUpdate.Destination.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "goos", firstUpdate.GOOS.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "goarch", firstUpdate.GOARCH.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "binary_location", firstUpdate.BinaryLocation.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "binary_hash", firstUpdate.BinaryHash.ValueString()),
				),
			},
			// Second update testing
			{
				Config: compilerDataSourceFromModel(t, secondUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "source", secondUpdate.Source.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "destination", secondUpdate.Destination.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "goos", secondUpdate.GOOS.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "goarch", secondUpdate.GOARCH.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "binary_location", secondUpdate.BinaryLocation.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "binary_hash", secondUpdate.BinaryHash.ValueString()),
				),
			},
			// Third update testing
			{
				Config: compilerDataSourceFromModel(t, thirdUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "source", thirdUpdate.Source.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "destination", thirdUpdate.Destination.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "goos", thirdUpdate.GOOS.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "goarch", thirdUpdate.GOARCH.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "binary_location", thirdUpdate.BinaryLocation.ValueString()+".zip"),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "binary_hash", "ff"),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "zip", "true"),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "zip_resources.../../LICENSE", "LICENSE"),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "zip_resources.../provider", "a/provider"),
				),
			},
		},
	})
}

func compilerDataSourceFromModel(t *testing.T, model CompileDataSourceModel) string {
	t.Helper()

	zip := ""
	zipResource := ""

	if !model.ZIP.IsNull() && !model.ZIP.IsUnknown() && model.ZIP.ValueBool() {
		zip = `zip = true`
	}

	if !model.ZIPResources.IsNull() && !model.ZIPResources.IsUnknown() {
		additionalFiles := map[string]string{}
		zipResource += "zip_resources = {\n"

		diag := model.ZIPResources.ElementsAs(context.Background(), &additionalFiles, false)
		assert.False(t, diag.HasError())
		for k, v := range additionalFiles {
			zipResource += fmt.Sprintf("		\"%s\" = \"%s\"\n", k, v)
		}

		zipResource += "	}"
	}

	return fmt.Sprintf(`
data "gopackager_compile" "test" {
	source = %s
	destination = %s
	goos = %s
	goarch = %s
	%s
	%s
}
	`, model.Source.String(), model.Destination.String(), model.GOOS.String(), model.GOARCH.String(), zip, zipResource)
}
