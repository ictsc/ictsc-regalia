### 依存関係のインストール
FROM node:16 AS deps

WORKDIR /app

COPY package.json yarn.lock ./

RUN yarn install --frozen-lockfile


## Font のビルド
FROM alpine:3.14 AS font-builder

WORKDIR /font

RUN apk add --no-cache git make g++ && \
    git clone --recursive https://github.com/google/woff2.git

WORKDIR /font/woff2

RUN make clean all

WORKDIR /font

RUN wget https://github.com/googlefonts/noto-cjk/releases/download/Sans2.004/01_NotoSansCJK-OTF-VF.zip && \
    unzip 01_NotoSansCJK-OTF-VF.zip && \
    rm 01_NotoSansCJK-OTF-VF.zip && \
    ./woff2/woff2_compress Variable/OTF/Subset/NotoSansJP-VF.otf


### Next.js のビルド
FROM node:16-alpine AS builder

ARG next_public_api_url
ARG next_public_answer_limit

ENV NEXT_PUBLIC_API_URL=$next_public_api_url
ENV NEXT_PUBLIC_ANSWER_LIMIT=$next_public_answer_limit

WORKDIR /app

COPY --from=deps /app/node_modules ./node_modules
COPY . .
COPY --from=font-builder /font/Variable/OTF/Subset/NotoSansJP-VF.woff2 ./pages/NotoSansJP-VF.woff2

RUN yarn build


### Next.js の実行
FROM node:16-alpine AS runner

ENV NODE_ENV=production

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT 3000

CMD ["node", "server.js"]