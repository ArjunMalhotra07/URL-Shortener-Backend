FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -g '' appuser && \
    mkdir -p /certs
ADD https://truststore.pki.rds.amazonaws.com/global/global-bundle.pem /certs/global-bundle.pem
RUN chmod 644 /certs/global-bundle.pem
WORKDIR /app
COPY --from=builder /app/server .
COPY --from=builder /app/db/migrations ./db/migrations
USER appuser
EXPOSE 8080
ENTRYPOINT ["./server"]
