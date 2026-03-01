# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /mcp-graph ./servers/mcp-graph/cmd/mcp-graph

FROM alpine:3.21

LABEL org.opencontainers.image.title="mcp-graph" \
      org.opencontainers.image.description="MCP server for knowledge graph associations" \
      org.opencontainers.image.source="https://github.com/Arkestone/mcp" \
      org.opencontainers.image.license="MIT"

RUN apk add --no-cache ca-certificates && \
    adduser -D -u 1000 mcp

RUN mkdir -p /home/mcp/.local/share/mcp-graph && chown mcp:mcp /home/mcp/.local/share/mcp-graph

USER mcp
COPY --from=build /mcp-graph /usr/local/bin/mcp-graph
ENTRYPOINT ["mcp-graph"]
