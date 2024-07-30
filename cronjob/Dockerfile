# First stage - Build the CronJob application
FROM golang:1.22-alpine AS build

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY ../go.mod ../go.sum ./
RUN go mod download

COPY . .

RUN go build -o /fetch-news ./main.go

# Second stage - Create the final image
FROM alpine:3.20

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /root/

COPY --from=build /fetch-news .

ENTRYPOINT ["/fetch-news"]
