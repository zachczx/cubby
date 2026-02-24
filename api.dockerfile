FROM golang:1.26-alpine AS builder
ENV GO111MODULE=on
RUN apk add --no-cache build-base

WORKDIR /api

COPY /api/go.mod /api/go.sum ./
RUN go mod download

COPY /api/ ./

# Build
RUN GOOS=linux go build -o /api/cubby-api ./cmd/api

####################################################################################

FROM alpine:latest
WORKDIR /api

COPY --from=builder /api/cubby-api .

ARG DB_HOST
ARG DB_PORT
ARG DB_USER
ARG DB_PASSWORD
ARG DB_NAME
ARG API_LISTEN_ADDR
ARG STYTCH_PROJECT_ID
ARG STYTCH_SECRET
ARG COOKIE_DOMAIN
ARG PUBLIC_CORS_URL

ENV DB_HOST=${DB_HOST}
ENV DB_PORT=${DB_PORT}
ENV DB_USER=${DB_USER}
ENV DB_PASSWORD=${DB_PASSWORD}
ENV DB_NAME=${DB_NAME}
ENV API_LISTEN_ADDR=${API_LISTEN_ADDR}
ENV STYTCH_PROJECT_ID=${STYTCH_PROJECT_ID}
ENV STYTCH_SECRET=${STYTCH_SECRET}
ENV COOKIE_DOMAIN=${COOKIE_DOMAIN}
ENV PUBLIC_CORS_URL=${PUBLIC_CORS_URL}

RUN apk add --no-cache wget
EXPOSE ${API_LISTEN_ADDR}

# Run
CMD ["/api/cubby-api"]

