package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stevencyb/terraform-provider-gopackager/internal/compiler"
)

func TestAccDataSourceFrameworkSatisfaction(t *testing.T) {
	t.Parallel()

	var _ datasource.DataSource = &CompileDataSource{}
}

func TestAccCompileDataSource(t *testing.T) {
	t.Parallel()

	mockCompiler := compiler.MockCompiler{}
	globalCompiler = &mockCompiler
	testAccProtoV6ProviderFactories := map[string]func() (tfprotov6.ProviderServer, error){
		"gopackager": providerserver.NewProtocol6WithError(New("test")()),
	}

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

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: compilerDataSourceFromModel("test", initialDataSource),
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
				Config: compilerDataSourceFromModel("test", firstUpdate),
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
				Config: compilerDataSourceFromModel("test", secondUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "source", secondUpdate.Source.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "destination", secondUpdate.Destination.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "goos", secondUpdate.GOOS.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "goarch", secondUpdate.GOARCH.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "binary_location", secondUpdate.BinaryLocation.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "binary_hash", secondUpdate.BinaryHash.ValueString()),
				),
			},
		},
	})
}

func compilerDataSourceFromModel(name string, model CompileDataSourceModel) string {
	return fmt.Sprintf(`
data "gopackager_compile" "%s" {
	source = %s
	destination = %s
	goos = %s
	goarch = %s
}
	`, name, model.Source.String(), model.Destination.String(), model.GOOS.String(), model.GOARCH.String())
}
