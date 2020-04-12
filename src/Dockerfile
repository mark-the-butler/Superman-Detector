FROM golang:alpine AS builder

ENV GO111MODULE=on \
    GOOS=linux

RUN apk add --update gcc musl-dev

WORKDIR /build

COPY . .

RUN go mod download

RUN go build -o Superman-Detector .

WORKDIR /dist

RUN cp /build/Superman-Detector .

FROM alpine:latest

COPY --from=builder /dist /

RUN apk add --update sqlite

ENTRYPOINT ["/Superman-Detector"]