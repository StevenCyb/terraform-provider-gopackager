HOSTNAME=github.com
NAMESPACE=stevencyb
NAME=gopackager
BINARY=terraform-provider-${NAME}
VERSION=0.2.5
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

local_mac_install:
	@echo "Cleaning up"
	@rm -Rf ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS}_${OS_ARCH} || true
	@echo "Building ${BINARY}"
	@go build -o ${BINARY} main.go
	@echo "Install"
	@mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS}_${OS_ARCH}
	@mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS}_${OS_ARCH}
	@echo "Installed to ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS}_${OS_ARCH}"
