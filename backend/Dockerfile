# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.23.4 AS builder
WORKDIR /src
RUN --mount=type=cache,target=/go/pkg/mod,sharing=locked \
    --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    go mod download -x

ARG TARGETOS
ARG TARGETARCH
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
ENV CGO_ENABLED=0

# --- score server ---
FROM --platform=$BUILDPLATFORM builder AS scoreserver-builder
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    --mount=type=bind,source=.,target=. \
    go build -o /scoreserver -ldflags "-s -w" -trimpath ./cmd/scoreserver
RUN env

FROM gcr.io/distroless/static-debian12:nonroot AS scoreserver

WORKDIR /app
EXPOSE 8080 8081

COPY --from=scoreserver-builder --chown=nonroot:nonroot /scoreserver /

ENTRYPOINT ["/scoreserver"]
