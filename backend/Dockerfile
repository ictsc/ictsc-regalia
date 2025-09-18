# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang AS builder
ENV GOTOOLCHAIN=auto
ENV CGO_ENABLED=0
WORKDIR /src
RUN --mount=type=bind,source=go.mod,target=go.mod \
    --mount=type=bind,source=go.sum,target=go.sum \
    go mod download
ARG TARGETOS
ARG TARGETARCH
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH

# --- backend ---
FROM --platform=$BUILDPLATFORM builder AS scoreserver-builder
RUN --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    --mount=type=bind,source=.,target=. \
    go build -o /scoreserver -ldflags "-s -w" -trimpath ./cmd/scoreserver

FROM gcr.io/distroless/static-debian12:nonroot AS backend
WORKDIR /app
EXPOSE 8080 8081
COPY --from=scoreserver-builder --chown=nonroot:nonroot /scoreserver /
ENTRYPOINT ["/scoreserver"]

# --- batch ---
FROM --platform=$BUILDPLATFORM builder AS batch-builder
RUN --mount=type=cache,target=/root/.cache/go-build,sharing=locked \
    --mount=type=bind,source=.,target=. \
    go build -o /batch -ldflags "-s -w" -trimpath ./cmd/batch

FROM gcr.io/distroless/static-debian12:nonroot AS batch
WORKDIR /app
COPY --from=batch-builder --chown=nonroot:nonroot /batch /
ENTRYPOINT ["/batch"]

# -- toolbox ---
FROM debian:13-slim AS toolbox
# Install tools
RUN --mount=type=cache,target=/var/lib/apt/ \
    --mount=type=cache,target=/var/cache/apt/ \
    apt-get update && apt-get install -y \
      curl \
      postgresql-client
ARG TARGETOS
ARG TARGETARCH
RUN --mount=type=bind,source=tools/go.mod,target=go.mod \
<<-EOS
SQLDEF_VERSION="$(grep 'github.com/sqldef/sqldef v' go.mod | cut -d" " -f2 | head -n1)"
curl -fsSL -o psqldef.tar.gz "https://github.com/sqldef/sqldef/releases/download/${SQLDEF_VERSION}/psqldef_${TARGETOS}_${TARGETARCH}.tar.gz"
tar -xf psqldef.tar.gz -C /usr/local/bin
rm psqldef.tar.gz
chmod +x /usr/local/bin/psqldef
EOS
# Add a user
ARG UID=1000
RUN useradd -u "${UID}" -m -s /bin/bash -c "" -p "!" toolbox
USER toolbox
# Copy scripts
WORKDIR /home/toolbox
COPY --link schema.sql view.sql .
COPY --link scripts ./scripts
