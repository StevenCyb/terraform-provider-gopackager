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

# `binary_location` provides the path and file name of the compiled binary.
# If `zip = true`, this will refer to the zip file.
output "binary_location" {
  value = data.gopackager_compile.example.binary_location
}

# `binary_hash` provides the hash of the compiled binary.
# There are multiple factors that can affect the hash, that means
# false positive changes are possible.
output "binary_hash" {
  value = data.gopackager_compile.example.binary_hash
}

# Example on how to use it with AWS lambda.
resource "aws_lambda_function" "example" {
  function_name    = "example"
  runtime          = "go1.x"
  handler          = "service"
  role             = aws_iam_role.lambda_role.arn
  timeout          = 15
  filename         = data.gopackager_compile.example.binary_location
  source_code_hash = data.gopackager_compile.example.binary_hash
  memory_size      = 128
}
