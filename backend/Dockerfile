FROM --platform=$BUILDPLATFORM golang:1.22.1 AS builder

WORKDIR /workdir

ENV GOCACHE=/tmp/go/cache
ENV CGO_ENABLED=0

ARG TARGETOS
ARG TARGETARCH
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH

COPY ./go.* ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .

# ---------------------anita---------------------

FROM --platform=$BUILDPLATFORM builder AS anita-builder

RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/tmp/go/cache \
  go build -o /anita -ldflags "-s -w" -trimpath

FROM gcr.io/distroless/base:latest AS anita

WORKDIR /app
EXPOSE 8080

COPY --from=anita-builder /anita .

HEALTHCHECK CMD ["curl", "localhost:8080/ping"]
ENTRYPOINT ["./anita"]
CMD ["run"]
