# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM node:24.13.0-slim AS builder
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"
RUN npm install -g pnpm@10.33.0
WORKDIR /app
RUN --mount=type=bind,source=package.json,target=package.json \
    --mount=type=bind,source=pnpm-lock.yaml,target=pnpm-lock.yaml \
    pnpm fetch
COPY . .
RUN pnpm install --recursive --offline
RUN pnpm run build

FROM nginx:1.31.4 AS frontend
ENV API_UPSTREAM=http://api:8080
ENV NGINX_RESOLVER=127.0.0.11
COPY <<EOF /etc/nginx/templates/default.conf.template
server {
    listen 8080;
    listen [::]:8080;
    server_name _;
    absolute_redirect off;

    location = /healthz {
        access_log off;
        return 200;
    }

    # Resolve lazily so ingress-only deployments can start without an api DNS record.
    resolver \${NGINX_RESOLVER} ipv6=off valid=10s;
    set \$api_upstream \${API_UPSTREAM};
    location /api/ {
        proxy_pass \$api_upstream\$request_uri;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_buffering off;
        proxy_cache off;
        proxy_read_timeout 1h;
    }

    location = /admin {
        return 308 /admin/;
    }

    location ^~ /admin/ {
        root /var/www;
        index index.html;
        try_files \$uri \$uri/ /admin/index.html;
    }

    location / {
        root /var/www/contestant/;
        index index.html;
        try_files \$uri \$uri/index.html /index.html;
    }
}
EOF
COPY --from=builder /app/packages/contestant/.output/public /var/www/contestant
COPY --from=builder /app/packages/admin/.output/public /var/www/admin
EXPOSE 8080
