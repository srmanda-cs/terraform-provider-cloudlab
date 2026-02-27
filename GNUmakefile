default: build

.PHONY: build
build:
	go build -v ./...

.PHONY: test
test:
	go test -v -count=1 -timeout 120s ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: fmt
fmt:
	gofmt -w -s .

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: generate
generate:
	go generate ./...

.PHONY: clean
clean:
	rm -f terraform-provider-cloudlab

.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/srmanda-cs/cloudlab/0.1.0/linux_amd64
	cp terraform-provider-cloudlab ~/.terraform.d/plugins/registry.terraform.io/srmanda-cs/cloudlab/0.1.0/linux_amd64/
