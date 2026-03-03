.PHONY: all build build-instructions build-skills build-adr build-memory build-prompts build-graph test test-race test-integration cover lint fmt vet clean docker docker-instructions docker-skills docker-adr docker-memory docker-prompts docker-graph

GOBIN    ?= bin
MODULE   := github.com/Arkestone/mcp

all: lint test build

# Build targets
build: build-instructions build-skills build-adr build-memory build-prompts build-graph

build-instructions:
	go build -buildvcs=false -trimpath -o $(GOBIN)/mcp-instructions ./servers/mcp-instructions/cmd/mcp-instructions

build-skills:
	go build -buildvcs=false -trimpath -o $(GOBIN)/mcp-skills ./servers/mcp-skills/cmd/mcp-skills

build-adr:
	go build -buildvcs=false -trimpath -o $(GOBIN)/mcp-adr ./servers/mcp-adr/cmd/mcp-adr

build-memory:
	go build -buildvcs=false -trimpath -o $(GOBIN)/mcp-memory ./servers/mcp-memory/cmd/mcp-memory

build-prompts:
	go build -buildvcs=false -trimpath -o $(GOBIN)/mcp-prompts ./servers/mcp-prompts/cmd/mcp-prompts

build-graph:
	go build -buildvcs=false -trimpath -o $(GOBIN)/mcp-graph ./servers/mcp-graph/cmd/mcp-graph

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
docker: docker-instructions docker-skills docker-adr docker-memory docker-prompts docker-graph

docker-instructions:
	docker build -f servers/mcp-instructions/Dockerfile -t ghcr.io/arkestone/mcp-instructions:latest .

docker-skills:
	docker build -f servers/mcp-skills/Dockerfile -t ghcr.io/arkestone/mcp-skills:latest .

docker-adr:
	docker build -f servers/mcp-adr/Dockerfile -t ghcr.io/arkestone/mcp-adr:latest .

docker-memory:
	docker build -f servers/mcp-memory/Dockerfile -t ghcr.io/arkestone/mcp-memory:latest .

docker-prompts:
	docker build -f servers/mcp-prompts/Dockerfile -t ghcr.io/arkestone/mcp-prompts:latest .

docker-graph:
	docker build -f servers/mcp-graph/Dockerfile -t ghcr.io/arkestone/mcp-graph:latest .

# Clean
clean:
	rm -rf $(GOBIN) cover.out coverage.out coverage.html
