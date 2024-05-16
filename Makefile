HOSTNAME=github.com
NAMESPACE=stevencyb
NAME=gopackager
BINARY=terraform-provider-${NAME}
VERSION=0.3.0
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
	@echo "Cleanup test directory"
	@rm -Rf test/.terraform
	@rm -f test/.terraform.lock.hcl
	@echo "Prepare test directory"
	@sed -i -r -E 's/"[0-9]+\.[0-9]+\.[0-9]+"/"${VERSION}"/g' 'test/main.tf'
	@cd test && terraform init