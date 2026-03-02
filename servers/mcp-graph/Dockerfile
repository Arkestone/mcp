# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Create data directory owned by distroless nonroot user (UID/GID 65532)
RUN mkdir -p /home/nonroot/.local/share/mcp-graph && \
    chown 65532:65532 /home/nonroot/.local/share/mcp-graph

ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /mcp-graph ./servers/mcp-graph/cmd/mcp-graph

FROM gcr.io/distroless/static-debian12:nonroot

LABEL org.opencontainers.image.title="mcp-graph" \
      org.opencontainers.image.description="MCP server for knowledge graph associations" \
      org.opencontainers.image.source="https://github.com/Arkestone/mcp" \
      org.opencontainers.image.license="MIT"

COPY --from=build /home/nonroot/.local/share/mcp-graph /home/nonroot/.local/share/mcp-graph
COPY --from=build /mcp-graph /mcp-graph
ENTRYPOINT ["/mcp-graph"]
