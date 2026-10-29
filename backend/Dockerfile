# syntax=docker/dockerfile:1

FROM golang:1.27.3-bookworm AS base

ENV GOTOOLCHAIN=auto
WORKDIR /workspace

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

FROM base AS development

COPY . .

EXPOSE 8080
CMD ["go", "tool", "air", "-c", ".air.toml"]

FROM base AS builder

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /api ./cmd/api

# Preview-only, PR-scoped database migration image. Never receives shared DB credentials.
FROM postgres:16-alpine AS migrations
COPY db/migrations/ /migrations/
COPY scripts/seed-dev.sql /migrations/seed-dev.sql

FROM gcr.io/distroless/static-debian12:nonroot AS production

WORKDIR /app
COPY --from=builder --chown=nonroot:nonroot /api /api

EXPOSE 8080
ENTRYPOINT ["/api"]
