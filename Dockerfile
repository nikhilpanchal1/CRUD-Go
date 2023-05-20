#FROM golang:1.19.0
FROM golang:1.20.4-bullseye

WORKDIR /usr/src/app

#hot reload
RUN go install github.com/cosmtrek/air@latest

COPY . .
RUN go mod tidy
