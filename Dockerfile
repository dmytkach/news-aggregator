# First stage - Build the application
FROM golang:1.22-alpine AS build

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /news-aggregator ./server/main.go

# Second stage - Create the final image
FROM alpine:latest

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /root/

COPY --from=build /news-aggregator .
COPY --from=build /app/certificates ./certificates
COPY --from=build /app/server-resources ./server-resources
COPY --from=build /app/server-news ./server-news
ENV FETCH_INTERVAL=1h0m
EXPOSE 8080

CMD ["./news-aggregator"]
