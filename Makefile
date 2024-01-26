HOSTNAME=github.com
NAMESPACE=stevencyb
NAME=gopackager
BINARY=terraform-provider-${NAME}
VERSION=v0.1.0
# // THIS NEEDS to be adjusted for your OS and ARCHITECTURE
OS=darwin
OS_ARCH=arm64

default: testacc

doc:
	@echo "Generating docs"
	@go generate

# Run acceptance tests
.PHONY: testacc
testacc:
	@TF_ACC=1 go test ./... -v -timeout 120m

local_linux_install:
	@echo "Cleaning up"
	@rm ${BINARY} || true
	@rm -Rf ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS}_${OS_ARCH} || true
	@echo "Building ${BINARY}"
	@go build -o ${BINARY} main.go
	@echo "Install"
	@mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS}_${OS_ARCH}
	@mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS}_${OS_ARCH}
	@echo "Installed to ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS}_${OS_ARCH}"
