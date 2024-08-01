# First stage - Build the CronJob application
FROM golang:1.22-alpine AS build

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY cronjob/go.mod cronjob/go.sum ./
COPY cronjob/ ./cronjob/
COPY server/ ./server/
COPY internal/ ./internal/

WORKDIR /app/cronjob

RUN go mod download

RUN go build -o /cronjob ./main.go

# Second stage - Create the final image
FROM alpine:3.20

COPY --from=build /cronjob .

ENTRYPOINT ["/cronjob"]
