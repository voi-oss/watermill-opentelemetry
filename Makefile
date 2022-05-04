.PHONY: lint
lint:
	golangci-lint run --config .golangci.yaml

.PHONY: deps
deps:
	go get ./...

.PHONY: test
test:
	go test -v ./...
