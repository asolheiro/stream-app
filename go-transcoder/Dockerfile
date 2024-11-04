FROM golang:1.23-alpine

LABEL authors="asolheiro"

RUN apk add --no-cache bash ffmpeg

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./ .

CMD go run ./cmd/mediaconverter/main.go