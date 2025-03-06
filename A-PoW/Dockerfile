FROM golang:1.23-alpine

ARG APP_DIR

WORKDIR /app

COPY . .

RUN go mod download

EXPOSE 5000