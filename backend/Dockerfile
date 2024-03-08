FROM --platform=$BUILDPLATFORM golang:1.22.1 AS builder

ENV GOCACHE=/tmp/go/cache
ENV CGO_ENABLED=0

ARG TARGETOS
ARG TARGETARCH
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH

ENV GOPROXY=https://goproxy.io,direct

WORKDIR /workdir

COPY ./go.* ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .

# ---------------------anita---------------------

FROM --platform=$BUILDPLATFORM builder AS anita-builder

RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/tmp/go/cache \
  go build -o /anita -ldflags "-s -w" -trimpath ./cmd/anita

FROM gcr.io/distroless/static:nonroot AS anita

WORKDIR /app
EXPOSE 8080

COPY --from=anita-builder --chown=nonroot:nonroot /anita .

HEALTHCHECK CMD ["./anita", "ping"]
ENTRYPOINT ["./anita"]
CMD ["run"]
