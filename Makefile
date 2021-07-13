PHONY: lint
lint:
	golangci-lint run --config .golangci.yaml

.PHONY: deps
deps: vendor

.PHONY: vendor
vendor:
	go mod vendor
