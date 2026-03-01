# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /mcp-instructions ./servers/mcp-instructions/cmd/mcp-instructions

FROM alpine:3.21

LABEL org.opencontainers.image.title="mcp-instructions" \
      org.opencontainers.image.description="MCP server for GitHub Copilot custom instructions" \
      org.opencontainers.image.source="https://github.com/Arkestone/mcp" \
      org.opencontainers.image.license="MIT"

# Health check: not applicable for stdio transport.
# For HTTP mode, uncomment and adjust:
# HEALTHCHECK --interval=30s --timeout=5s CMD wget -qO- http://localhost:8080/health || exit 1

RUN apk add --no-cache ca-certificates && \
    adduser -D -u 1000 mcp

USER mcp
COPY --from=build /mcp-instructions /usr/local/bin/mcp-instructions
ENTRYPOINT ["mcp-instructions"]
