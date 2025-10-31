package provider

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stevencyb/gopackager/internal/compiler"
	"github.com/stevencyb/gopackager/internal/hasher"
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
	mockHasher := hasher.MockHasher{}
	globalCompiler = &mockCompiler
	globalZIPPackager = &mockPackager
	globalHasher = &mockHasher
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
		Source:             types.StringValue("provider.go"),
		Destination:        types.StringValue("linux_amd64_binary"),
		GOOS:               types.StringValue("linux"),
		GOARCH:             types.StringValue("amd64"),
		OutputPath:         types.StringValue("linux_amd64_binary"),
		OutputMD5:          types.StringValue("md5hash"),
		OutputSHA1:         types.StringValue("sha1hash"),
		OutputSHA256:       types.StringValue("sha256hash"),
		OutputSHA512:       types.StringValue("sha512hash"),
		OutputSHA256Base64: types.StringValue("sha256base64hash"),
		OutputSHA512Base64: types.StringValue("sha512base64hash"),
	}
	firstUpdate := CompileDataSourceModel{
		Source:             types.StringValue("provider.go"),
		Destination:        types.StringValue("windows_amd64_binary"),
		GOOS:               types.StringValue("windows"),
		GOARCH:             types.StringValue("amd64"),
		OutputPath:         types.StringValue("windows_amd64_binary"),
		OutputMD5:          types.StringValue("md5hash"),
		OutputSHA1:         types.StringValue("sha1hash"),
		OutputSHA256:       types.StringValue("sha256hash"),
		OutputSHA512:       types.StringValue("sha512hash"),
		OutputSHA256Base64: types.StringValue("sha256base64hash"),
		OutputSHA512Base64: types.StringValue("sha512base64hash"),
	}
	secondUpdate := CompileDataSourceModel{
		Source:             types.StringValue("provider.go"),
		Destination:        types.StringValue("windows_amd64_binary"),
		GOOS:               types.StringValue("windows"),
		GOARCH:             types.StringValue("amd64"),
		OutputPath:         types.StringValue("windows_amd64_binary"),
		OutputMD5:          types.StringValue("md5hash"),
		OutputSHA1:         types.StringValue("sha1hash"),
		OutputSHA256:       types.StringValue("sha256hash"),
		OutputSHA512:       types.StringValue("sha512hash"),
		OutputSHA256Base64: types.StringValue("sha256base64hash"),
		OutputSHA512Base64: types.StringValue("sha512base64hash"),
	}
	thirdUpdate := CompileDataSourceModel{
		Source:             types.StringValue("provider.go"),
		Destination:        types.StringValue("windows_amd64_binary"),
		GOOS:               types.StringValue("windows"),
		GOARCH:             types.StringValue("amd64"),
		OutputPath:         types.StringValue("windows_amd64_binary"),
		OutputMD5:          types.StringValue("md5hash"),
		OutputSHA1:         types.StringValue("sha1hash"),
		OutputSHA256:       types.StringValue("sha256hash"),
		OutputSHA512:       types.StringValue("sha512hash"),
		OutputSHA256Base64: types.StringValue("sha256base64hash"),
		OutputSHA512Base64: types.StringValue("sha512base64hash"),
		ZIP:                types.BoolValue(true),
		ZIPResources:       additionalZIPResourcesGen,
	}

	mockHasher.On("ReadFile", initialDataSource.OutputPath.ValueString()).Times(3).Return([]byte("123"), nil)
	mockHasher.On("CombinedHash", []byte("123")).Times(3).Return(hasher.CombinedHash{
		MD5:          initialDataSource.OutputMD5.ValueString(),
		SHA1:         initialDataSource.OutputSHA1.ValueString(),
		SHA256:       initialDataSource.OutputSHA256.ValueString(),
		SHA512:       initialDataSource.OutputSHA512.ValueString(),
		SHA256Base64: initialDataSource.OutputSHA256Base64.ValueString(),
		SHA512Base64: initialDataSource.OutputSHA512Base64.ValueString(),
	}, nil)

	basePath := filepath.Dir(initialDataSource.Source.ValueString())
	mockHasher.On("HashDir", basePath).Return(&hasher.CombinedHash{
		MD5:          initialDataSource.OutputMD5.ValueString(),
		SHA1:         initialDataSource.OutputSHA1.ValueString(),
		SHA256:       initialDataSource.OutputSHA256.ValueString(),
		SHA512:       initialDataSource.OutputSHA512.ValueString(),
		SHA256Base64: initialDataSource.OutputSHA256Base64.ValueString(),
		SHA512Base64: initialDataSource.OutputSHA512Base64.ValueString(),
	}, nil)
	mockCompiler.On("Compile",
		*compiler.NewConfig().
			Source(initialDataSource.Source.ValueString()).
			Destination(initialDataSource.Destination.ValueString()).
			GOOS(initialDataSource.GOOS.ValueString()).
			GOARCH(initialDataSource.GOARCH.ValueString()),
	).Times(3).
		Return(initialDataSource.OutputPath.ValueString(), nil)

	mockHasher.On("ReadFile", firstUpdate.OutputPath.ValueString()).Times(6).Return([]byte("333"), nil)
	mockHasher.On("CombinedHash", []byte("333")).Times(6).Return(hasher.CombinedHash{
		MD5:          firstUpdate.OutputMD5.ValueString(),
		SHA1:         firstUpdate.OutputSHA1.ValueString(),
		SHA256:       firstUpdate.OutputSHA256.ValueString(),
		SHA512:       firstUpdate.OutputSHA512.ValueString(),
		SHA256Base64: firstUpdate.OutputSHA256Base64.ValueString(),
		SHA512Base64: firstUpdate.OutputSHA512Base64.ValueString(),
	}, nil)
	mockCompiler.On("Compile",
		*compiler.NewConfig().
			Source(firstUpdate.Source.ValueString()).
			Destination(firstUpdate.Destination.ValueString()).
			GOOS(firstUpdate.GOOS.ValueString()).
			GOARCH(firstUpdate.GOARCH.ValueString()),
	).Times(3).
		Return(firstUpdate.OutputPath.ValueString(), nil)

	// ReadFile and CombinedHash are reused
	mockCompiler.On("Compile",
		*compiler.NewConfig().
			Source(secondUpdate.Source.ValueString()).
			Destination(secondUpdate.Destination.ValueString()).
			GOOS(secondUpdate.GOOS.ValueString()).
			GOARCH(secondUpdate.GOARCH.ValueString()),
	).Times(3).
		Return(secondUpdate.OutputPath.ValueString(), nil)

	mockPackager.On("Zip", thirdUpdate.OutputPath.ValueString()+".zip", additionalZIPResources).Times(3).Return(nil)
	mockHasher.On("ReadFile", thirdUpdate.OutputPath.ValueString()+".zip").Times(3).Return([]byte("666"), nil)
	mockHasher.On("CombinedHash", []byte("666")).Times(3).Return(hasher.CombinedHash{
		MD5:          firstUpdate.OutputMD5.ValueString(),
		SHA1:         firstUpdate.OutputSHA1.ValueString(),
		SHA256:       firstUpdate.OutputSHA256.ValueString(),
		SHA512:       firstUpdate.OutputSHA512.ValueString(),
		SHA256Base64: firstUpdate.OutputSHA256Base64.ValueString(),
		SHA512Base64: firstUpdate.OutputSHA512Base64.ValueString(),
	}, nil)
	mockCompiler.On("Compile",
		*compiler.NewConfig().
			Source(thirdUpdate.Source.ValueString()).
			Destination(thirdUpdate.Destination.ValueString()).
			GOOS(thirdUpdate.GOOS.ValueString()).
			GOARCH(thirdUpdate.GOARCH.ValueString()),
	).Times(3).
		Return(thirdUpdate.OutputPath.ValueString(), nil)

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
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_path", initialDataSource.OutputPath.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_md5", initialDataSource.OutputMD5.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_sha1", initialDataSource.OutputSHA1.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_sha256", initialDataSource.OutputSHA256.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_sha512", initialDataSource.OutputSHA512.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_sha256_base64", initialDataSource.OutputSHA256Base64.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_sha512_base64", initialDataSource.OutputSHA512Base64.ValueString()),
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
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_path", firstUpdate.OutputPath.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_md5", firstUpdate.OutputMD5.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_sha1", firstUpdate.OutputSHA1.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_sha256", firstUpdate.OutputSHA256.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_sha512", firstUpdate.OutputSHA512.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_sha256_base64", firstUpdate.OutputSHA256Base64.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_sha512_base64", firstUpdate.OutputSHA512Base64.ValueString()),
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
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_path", secondUpdate.OutputPath.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_md5", secondUpdate.OutputMD5.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_sha1", secondUpdate.OutputSHA1.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_sha256", secondUpdate.OutputSHA256.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_sha512", secondUpdate.OutputSHA512.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_sha256_base64", secondUpdate.OutputSHA256Base64.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_sha512_base64", secondUpdate.OutputSHA512Base64.ValueString()),
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
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_path", thirdUpdate.OutputPath.ValueString()+".zip"),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_md5", thirdUpdate.OutputMD5.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_sha1", thirdUpdate.OutputSHA1.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_sha256", thirdUpdate.OutputSHA256.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_sha512", thirdUpdate.OutputSHA512.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_sha256_base64", thirdUpdate.OutputSHA256Base64.ValueString()),
					resource.TestCheckResourceAttr("data.gopackager_compile.test", "output_sha512_base64", thirdUpdate.OutputSHA512Base64.ValueString()),
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
