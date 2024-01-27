data "gopackager_compile" "example" {
  # Required
  ## Path to the main source file.
  source = "main.go"
  ## Compiled output destination file.
  destination = "service"
  ## GOOS for compilation.
  goos = "linux"
  ## GOARCH for compilation.
  goarch = "amd64"

  # Optional
  ## Zip the compiled binary and additional resources.
  zip = true
  ## Additional resources to be zipped.
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
    # `output_md5` provides the md5 hash of the compiled binary.
    # There are multiple factors that can affect the hash, that means
    output_md5 = data.gopackager_compile.example.output_md5
    # `output_sha1` provides the SHA1 hash of the compiled binary.
    # There are multiple factors that can affect the hash, that means
    output_sha1 = data.gopackager_compile.example.output_sha1
    # `output_sha256` provides the SHA256 hash of the compiled binary.
    # There are multiple factors that can affect the hash, that means
    output_sha256 = data.gopackager_compile.example.output_sha256
    # `output_sha512` provides the SHA512 hash of the compiled binary.
    # There are multiple factors that can affect the hash, that means
    output_sha512 = data.gopackager_compile.example.output_sha512
    # `output_sha256_base64` provides the Base64 ecndoded SHA256 hash of the compiled binary.
    # There are multiple factors that can affect the hash, that means
    output_sha256_base64 = data.gopackager_compile.example.output_sha256_base64
    # `output_sha512_base64` provides the Base64 ecndoded SHA512 hash of the compiled binary.
    # There are multiple factors that can affect the hash, that means
    output_sha512_base64 = data.gopackager_compile.example.output_sha512_base64
  }
}

# Example on how to use it with AWS lambda.
resource "aws_lambda_function" "example" {
  function_name    = "example"
  runtime          = "go1.x"
  handler          = "service"
  role             = aws_iam_role.lambda_role.arn
  timeout          = 15
  filename         = data.gopackager_compile.example.output_path
  source_code_hash = data.gopackager_compile.example.output_md5
  memory_size      = 128
}
