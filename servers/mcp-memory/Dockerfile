# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /mcp-memory ./servers/mcp-memory/cmd/mcp-memory

FROM alpine:3.21

LABEL org.opencontainers.image.title="mcp-memory" \
      org.opencontainers.image.description="MCP server for persistent memory storage" \
      org.opencontainers.image.source="https://github.com/Arkestone/mcp" \
      org.opencontainers.image.license="MIT"

RUN apk add --no-cache ca-certificates && \
    adduser -D -u 1000 mcp

RUN mkdir -p /home/mcp/.local/share/mcp-memory && chown mcp:mcp /home/mcp/.local/share/mcp-memory

USER mcp
COPY --from=build /mcp-memory /usr/local/bin/mcp-memory
ENTRYPOINT ["mcp-memory"]
