FROM golang:1.22-alpine AS build

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY cronjob/go.mod cronjob/go.sum ./
RUN go mod download

COPY cronjob/ ./cronjob/
COPY server/ ./server/
COPY internal/ ./internal/

WORKDIR /app/cronjob
RUN go build -o /news-fetcher ./main.go

FROM alpine:3.20

COPY --from=build /news-fetcher /news-fetcher

ENTRYPOINT ["/news-fetcher"]
