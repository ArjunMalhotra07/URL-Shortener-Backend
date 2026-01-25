FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/server ./server
EXPOSE 8080
ENTRYPOINT ["/app/server"]
