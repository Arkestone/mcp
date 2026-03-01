# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /mcp-skills ./servers/mcp-skills/cmd/mcp-skills

FROM alpine:3.21

LABEL org.opencontainers.image.title="mcp-skills" \
      org.opencontainers.image.description="MCP server for Copilot skills" \
      org.opencontainers.image.source="https://github.com/Arkestone/mcp" \
      org.opencontainers.image.license="MIT"

RUN apk add --no-cache ca-certificates && \
    adduser -D -u 1000 mcp

USER mcp
COPY --from=build /mcp-skills /usr/local/bin/mcp-skills
ENTRYPOINT ["mcp-skills"]
