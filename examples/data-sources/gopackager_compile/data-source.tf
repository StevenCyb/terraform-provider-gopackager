data "gopackager_compile" "example" {
  # Path to the main source file.
  source = "main.go"
  # Compiled output destination file.
  destination = "service"
  # GOOS for compilation.
  goos = "linux"
  # GOARCH for compilation.
  goarch = "amd64"
}

# `binary_location` provides the path and file name of the compiled binary.
output "binary_location" {
  value = data.gopackager_compile.example.binary_location
}

# `binary_hash` provides the hash of the compiled binary.
# There are multiple factors that can affect the hash, that means
# false positive changes are possible.
output "binary_hash" {
  value = data.gopackager_compile.example.binary_hash
}
