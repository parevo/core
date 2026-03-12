.PHONY: test vet fmt lint build clean release

test:
	go test ./...

vet:
	go vet ./...

fmt:
	go fmt ./...

lint: vet
	@command -v golangci-lint >/dev/null 2>&1 || (echo "golangci-lint not installed, run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run ./...

build:
	go build ./...

clean:
	go clean -cache -testcache

release:
	@echo "Tag with: git tag v0.1.0 && git push origin v0.1.0"
	@echo "Or use goreleaser if configured"
