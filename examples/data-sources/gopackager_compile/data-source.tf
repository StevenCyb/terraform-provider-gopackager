data "gopackager_compile" "example" {
  # Required
  ## Path to the main GoLang source or the root path of this file.
  source = "main.go"
  ## Output destination file.
  destination = "service/bootstrap"
  ## GOOS for compilation.
  goos = "linux"
  ## GOARCH for compilation.
  goarch = "amd64"

  # Optional
  ## Zip the compiled binary and additional resources.
  zip = true
  ## Additional resources to be zipped.
  ## {source_path = destination_path}
  zip_resources = {
    "static"  = "www/static"
    "LICENSE" = "LICENSE"
  }
}

output "example" {
  value = {
    # `output_path` provides the path and file name of the compiled binary.
    # If `zip = true`, this will refer to the zip file.
    output_path = data.gopackager_compile.example.output_path
    # `output_md5` provides the md5 hash of the compiled binary or compressed ZIP file as hexadecimal encoded.
    # There are multiple factors that can affect the hash, that means
    output_md5 = data.gopackager_compile.example.output_md5
    # `output_sha1` provides the SHA1 hash of the compiled binary or compressed ZIP file as hexadecimal encoded.
    # There are multiple factors that can affect the hash, that means
    output_sha1 = data.gopackager_compile.example.output_sha1
    # `output_sha256` provides the SHA256 hash of the compiled binary or compressed ZIP file as hexadecimal encoded.
    # There are multiple factors that can affect the hash, that means
    output_sha256 = data.gopackager_compile.example.output_sha256
    # `output_sha512` provides the SHA512 hash of the compiled binary or compressed ZIP file as hexadecimal encoded.
    # There are multiple factors that can affect the hash, that means
    output_sha512 = data.gopackager_compile.example.output_sha512
    # `output_sha256_base64` provides the Base64 encoded SHA256 hash of the compiled binary or compressed ZIP file.
    # There are multiple factors that can affect the hash, that means
    output_sha256_base64 = data.gopackager_compile.example.output_sha256_base64
    # `output_sha512_base64` provides the Base64 encoded SHA512 hash of the compiled binary or compressed ZIP file.
    # There are multiple factors that can affect the hash, that means
    output_sha512_base64 = data.gopackager_compile.example.output_sha512_base64
    # Last commit hash of the GoLang source code.
    # If returns the last commit of current branch or `unknown` on error.
    # Us this hash for if more consistent hash needed.
    output_git_hash = data.gopackager_compile.example.output_git_hash
  }
}

# Example on how to use it with AWS lambda.
resource "aws_lambda_function" "example" {
  function_name = "example"
  runtime       = "provided.al2023"
  handler       = "bootstrap"
  role          = aws_iam_role.lambda_role.arn
  timeout       = 15
  filename      = data.gopackager_compile.example.output_path
  # Lambda expect base64 encoded sha256 hash of the source code.
  source_code_hash = base64sha256(data.gopackager_compile.example.output_git_hash)
  memory_size      = 128
}
