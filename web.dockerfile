FROM node:24-slim
WORKDIR /web
COPY ./web/package.json ./web/pnpm-lock.yaml ./
RUN npm install -g corepack@latest &&\ 
    corepack enable &&\ 
    corepack prepare pnpm@latest --activate

RUN pnpm install --frozen-lockfile

COPY ./web .

ARG PUBLIC_API_URL
ARG FRONTEND_LISTEN_ADDR

ENV PUBLIC_API_URL=${PUBLIC_API_URL}
ENV FRONTEND_LISTEN_ADDR=${FRONTEND_LISTEN_ADDR}

RUN pnpm build

ENV NODE_ENV=production

# These are not necessary. Docker-compose settles it.
# ENV LISTEN_ADDR=localhost:3000
# EXPOSE 3000

CMD ["node", "build/index.js"]