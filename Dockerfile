ARG GO_IMAGE=golang:1.24
ARG APP_PACKAGE=""
ARG APP_VERSION="0.1.0-dev"

FROM --platform=$BUILDPLATFORM golang AS dev

WORKDIR /app

COPY . .

RUN go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.0.2 && \
    go install go.uber.org/nilaway/cmd/nilaway@latest && \
    go install golang.org/x/vuln/cmd/govulncheck@latest

# Keep the container running in dev mode and allow any
# kind of command execution.
ENTRYPOINT ["tail", "-f", "/dev/null"]

FROM --platform=$BUILDPLATFORM dev AS builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /app

RUN go mod tidy && go mod vendor
RUN go generate ./...

RUN --mount=type=cache,id="go-build-cache-${TARGETOS}-${TARGETARCH}",sharing=private,target=/root/.cache/go-build \
    --mount=type=cache,id="go-package-cache",sharing=shared,target=/go/pkg \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
    -ldflags "-X '$APP_PACKAGE/main.version=$APP_VERSION' -extldflags '-static'" \
    -o civic main.go

FROM scratch AS production

COPY --from=builder /app/civic /usr/bin/civic

ENTRYPOINT ["civic"]
CMD ["-h"]
