# First stage - Build the application
FROM golang:1.22-alpine AS build

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY server/certificates ./certificates
COPY ./internal ./internal
COPY ./server ./server

RUN go build -o /news-aggregator ./server/main.go

# Second stage - Create the final image
FROM alpine:3.20

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /root/

COPY --from=build /news-aggregator .
COPY --from=build /app/certificates ./server/certificates

# Expose port
EXPOSE 8443

ENTRYPOINT ["./news-aggregator"]
CMD ["./news-aggregator", "-fetch_interval=1h", "-port=:8443", "-cert=/certificates/cert.pem", "-key=/certificates/key.pem", "-path_to_source=sources.json", "-news_folder=server-news/"]
