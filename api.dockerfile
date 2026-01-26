FROM golang:alpine AS first
ENV GO111MODULE=on
RUN apk add build-base
WORKDIR /backend
COPY ./backend/go.mod ./backend/go.sum ./
COPY . ./
RUN go mod download

COPY ./backend ./

# Build
RUN GOOS=linux go build -o /backend/roamichi-ni

####################################################################################

FROM alpine
WORKDIR /backend
COPY --from=first /backend/roamichi-ni .
COPY --from=first /backend/locapi.db .

# Receive values from docker-compose's "build.args"
ARG DB_HOST
ARG DB_PORT
ARG DB_USER
ARG DB_PASSWORD
ARG DB_NAME
ARG BACKEND_LISTEN_ADDR
ARG GORILLA_SECRET_KEY
ARG STYTCH_PROJECT_ID
ARG STYTCH_SECRET

ENV DB_HOST=${DB_HOST}
ENV DB_PORT=${DB_PORT}
ENV DB_USER=${DB_USER}
ENV DB_PASSWORD=${DB_PASSWORD}
ENV DB_NAME=${DB_NAME}
ENV BACKEND_LISTEN_ADDR=${BACKEND_LISTEN_ADDR}
ENV GORILLA_SECRET_KEY=${GORILLA_SECRET_KEY}
ENV STYTCH_PROJECT_ID=${STYTCH_PROJECT_ID}
ENV STYTCH_SECRET=${STYTCH_SECRET}

RUN apk add wget
EXPOSE ${BACKEND_LISTEN_ADDR}

# Run
CMD ["/backend/roamichi-ni"]

