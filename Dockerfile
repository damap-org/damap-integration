# syntax=docker/dockerfile:1

ARG GO_VERSION=1.24

FROM golang:${GO_VERSION}-bookworm AS build

WORKDIR /src

ENV CGO_ENABLED=0 GOTOOLCHAIN=local

# Download only the service dependency graph, not the release/lint tool graph.
COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -mod=readonly -trimpath -ldflags="-s -w" -o /damap-integration ./cmd/damap-integration

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /damap-integration /damap-integration

USER 65532:65532

ENV DAMAP_INTEGRATION_LISTEN=0.0.0.0:8080
EXPOSE 8080

ENTRYPOINT ["/damap-integration"]
