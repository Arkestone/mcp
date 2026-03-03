# syntax=docker/dockerfile:1

FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG TARGETOS TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /mcp-prompts ./servers/mcp-prompts/cmd/mcp-prompts

FROM gcr.io/distroless/static-debian12:nonroot

LABEL org.opencontainers.image.title="mcp-prompts" \
      org.opencontainers.image.description="MCP server for GitHub Copilot prompt and chat mode files" \
      org.opencontainers.image.source="https://github.com/Arkestone/mcp" \
      org.opencontainers.image.license="MIT"

COPY --from=build /mcp-prompts /mcp-prompts
ENTRYPOINT ["/mcp-prompts"]
