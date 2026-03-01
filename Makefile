.PHONY: all build build-instructions build-skills test test-race test-integration cover lint fmt vet clean docker docker-instructions docker-skills

GOBIN    ?= bin
MODULE   := github.com/Arkestone/mcp

all: lint test build

# Build targets
build: build-instructions build-skills

build-instructions:
	go build -buildvcs=false -trimpath -o $(GOBIN)/mcp-instructions ./instructions/cmd/mcp-instructions

build-skills:
	go build -buildvcs=false -trimpath -o $(GOBIN)/mcp-skills ./skills/cmd/mcp-skills

# Test targets
test:
	go test -count=1 ./...

test-race:
	go test -race -count=1 ./...

test-integration:
	go test -tags integration -race -count=1 ./...

cover:
	go test -buildvcs=false -coverprofile=cover.out ./...
	go tool cover -func=cover.out
	@rm -f cover.out

cover-html:
	go test -buildvcs=false -coverprofile=cover.out ./...
	go tool cover -html=cover.out -o coverage.html
	@rm -f cover.out
	@echo "Coverage report: coverage.html"

# Quality targets
lint:
	@which golangci-lint > /dev/null 2>&1 || echo "golangci-lint not installed, skipping lint"
	@which golangci-lint > /dev/null 2>&1 && golangci-lint run ./... || true

fmt:
	gofmt -s -w .

vet:
	go vet ./...

# Docker targets
docker: docker-instructions docker-skills

docker-instructions:
	docker build -f instructions/Dockerfile -t ghcr.io/arkestone/mcp-instructions:latest .

docker-skills:
	docker build -f skills/Dockerfile -t ghcr.io/arkestone/mcp-skills:latest .

# Clean
clean:
	rm -rf $(GOBIN) cover.out coverage.out coverage.html
