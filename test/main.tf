terraform {
  required_providers {
    example = {
      version = "0.3.0"
      source  = "github.com/stevencyb/gopackager"
    }
  }
}

# Variations
data "gopackager_compile" "example_local" {
  provider = example

  source      = "../main.go"
  destination = "example_local"
  goarch      = "amd64"
  goos        = "linux"
}

data "gopackager_compile" "example_local_zip" {
  provider = example

  source      = "../main.go"
  destination = "example_local_zip"
  goarch      = "amd64"
  goos        = "linux"

  zip = true
  zip_resources = {
    "../README.md" = "README.md"
  }
}

data "gopackager_compile" "example_git" {
  provider = example

  source      = "main.go"
  destination = "example_git"
  goarch      = "amd64"
  goos        = "linux"

  git_repository = "https://github.com/StevenCyb/terraform-provider-gopackager.git"
  # branch       = "main"
}

# TODO zip case

# Outputs
output "example_local" {
  value = {
    output_path          = data.gopackager_compile.example_local.output_path
    output_md5           = data.gopackager_compile.example_local.output_md5
    output_sha1          = data.gopackager_compile.example_local.output_sha1
    output_sha256        = data.gopackager_compile.example_local.output_sha256
    output_sha512        = data.gopackager_compile.example_local.output_sha512
    output_sha256_base64 = data.gopackager_compile.example_local.output_sha256_base64
    output_sha512_base64 = data.gopackager_compile.example_local.output_sha512_base64
  }
}

output "example_local_zip" {
  value = {
    output_path          = data.gopackager_compile.example_local_zip.output_path
    output_md5           = data.gopackager_compile.example_local_zip.output_md5
    output_sha1          = data.gopackager_compile.example_local_zip.output_sha1
    output_sha256        = data.gopackager_compile.example_local_zip.output_sha256
    output_sha512        = data.gopackager_compile.example_local_zip.output_sha512
    output_sha256_base64 = data.gopackager_compile.example_local_zip.output_sha256_base64
    output_sha512_base64 = data.gopackager_compile.example_local_zip.output_sha512_base64
  }
}

output "example_git" {
  value = {
    output_path          = data.gopackager_compile.example_git.output_path
    output_md5           = data.gopackager_compile.example_git.output_md5
    output_sha1          = data.gopackager_compile.example_git.output_sha1
    output_sha256        = data.gopackager_compile.example_git.output_sha256
    output_sha512        = data.gopackager_compile.example_git.output_sha512
    output_sha256_base64 = data.gopackager_compile.example_git.output_sha256_base64
    output_sha512_base64 = data.gopackager_compile.example_git.output_sha512_base64
  }
}
