# Terraform Provider GoPackager (Terraform Plugin Framework)

[This Terraform provider](https://registry.terraform.io/providers/StevenCyb/gopackager/latest) is a helper to compile GoLang binaries with terraform.
In fact, Terraform is meant to be used to build infrastructure.
But in reality some small project with serverless code have a mono-repository.
Therefore I decided to build a custom Terraform provider to compile GoLang binaries to replace those "ugly" `local-exec` parts.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install)

## Documentations
* [GoPackager Provider](docs/index.md)
  * [Compile Datasource](docs/data-sources/compile.md)