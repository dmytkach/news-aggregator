# First stage - Build the application
FROM golang:1.22-alpine AS build

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY ./internal ./internal
COPY ./server ./server

RUN go build -o /news-aggregator ./server/main.go

# Second stage - Create the final image
FROM alpine:3.20

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /root/

COPY --from=build /news-aggregator .

EXPOSE 8443

ENTRYPOINT ["./news-aggregator", "-port=:8443", "-path-to-source=/mnt/sources/sources.json", "-news-folder=/mnt/news", "-tls-cert=/etc/tls/certs/tls.crt", "-tls-key=/etc/tls/certs/tls.key"]
