FROM golang:1.24.2 AS notification-builder
WORKDIR /app
COPY . .
RUN cd cmd && CGO_ENABLED=0 GOOS=linux go build -o notification-service main.go

FROM alpine AS group
WORKDIR /
COPY --from=notification-builder /app/cmd/notification-service ./
ENTRYPOINT ["./notification-service"]