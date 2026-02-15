FROM node:25-slim AS builder
WORKDIR /web
COPY /web/package.json /web/pnpm-lock.yaml ./
RUN npm install -g pnpm@latest-10
RUN pnpm install --frozen-lockfile
COPY /web ./

ARG PUBLIC_API_URL
ARG PUBLIC_WEB_URL

ENV PUBLIC_API_URL=${PUBLIC_API_URL}
ENV PUBLIC_WEB_URL=${PUBLIC_WEB_URL}

RUN pnpm build

FROM nginx:alpine
COPY --from=builder /web/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
