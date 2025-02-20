# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM node:22.13.1-slim AS builder
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN npm install -g pnpm
WORKDIR /app
RUN --mount=type=bind,source=package.json,target=package.json \
    --mount=type=bind,source=pnpm-lock.yaml,target=pnpm-lock.yaml \
    pnpm fetch
COPY . .
RUN pnpm install --recursive --offline
RUN pnpm run --recursive --parallel --aggregate-output build

FROM nginx:1.27.4 AS frontend
COPY <<EOF /etc/nginx/conf.d/default.conf
server {
    listen 8080;
    server_name _;

    location / {
        root /var/www/contestant;
        index index.html;
        try_files $uri $uri/index.html /index.html;
    }

    location = /healthz {
        access_log off;
        return 200;
    }
}
EOF
COPY --from=builder /app/packages/contestant/dist /var/www/contestant
EXPOSE 8080
