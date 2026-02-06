FROM golang:1.24-alpine AS first
ENV GO111MODULE=on
RUN apk add --no-cache build-base

WORKDIR /api

COPY ./api/go.mod ./api/go.sum ./
RUN go mod download

COPY ./api ./

# Build
RUN GOOS=linux go build -o /api/cubby-api .

####################################################################################

FROM alpine:latest
WORKDIR /api

COPY --from=first /api/cubby-api .
# COPY --from=first /api/locapi.db . # database is likely Postgres given pgx in go.mod, removing locapi.db which sounds like sqlite/local

# Receive values from docker-compose's "build.args"
ARG DB_HOST
ARG DB_PORT
ARG DB_USER
ARG DB_PASSWORD
ARG DB_NAME
ARG API_LISTEN_ADDR
ARG STYTCH_PROJECT_ID
ARG STYTCH_SECRET

ENV DB_HOST=${DB_HOST}
ENV DB_PORT=${DB_PORT}
ENV DB_USER=${DB_USER}
ENV DB_PASSWORD=${DB_PASSWORD}
ENV DB_NAME=${DB_NAME}
ENV API_LISTEN_ADDR=${API_LISTEN_ADDR}
ENV STYTCH_PROJECT_ID=${STYTCH_PROJECT_ID}
ENV STYTCH_SECRET=${STYTCH_SECRET}

RUN apk add --no-cache wget
EXPOSE ${API_LISTEN_ADDR}

# Run
CMD ["/api/cubby-api"]

