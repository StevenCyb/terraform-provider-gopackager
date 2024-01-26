terraform {
  required_providers {
    gopackager = {
      version = "0.1.0"
      source  = "github.com/stevencyb/gopackager"
    }
  }
}

data "gopackager_compile" "example" {
  source      = "../main.go"
  destination = "compiled_example"
  goos        = "linux"
  goarch      = "amd64"

  zip = true
  zip_resources = {
    "../LICENSE"           = "LICENSE"
    "../internal/provider" = "provider"
  }
}

output "binary_location" {
  value = data.gopackager_compile.example.binary_location
}

output "binary_hash" {
  value = data.gopackager_compile.example.binary_hash
}
