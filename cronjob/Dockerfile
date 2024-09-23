# First stage: build stage
FROM golang:1.22-alpine AS build

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY go.mod go.sum ./

COPY cronjob/go.mod cronjob/go.sum ./cronjob/

RUN go mod download

COPY cronjob/ ./cronjob/
COPY server/ ./server/
COPY internal/ ./internal/

WORKDIR /app/cronjob
RUN go build -o /news-fetcher ./main.go

# Second stage: final stage
FROM alpine:3.20

COPY --from=build /news-fetcher /news-fetcher

ENTRYPOINT ["/news-fetcher"]
