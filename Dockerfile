FROM golang:1.25.5-alpine AS builder
LABEL org.opencontainers.image.title="News-Parser"\
      org.opencontainers.image.description="service for parsing pages of Internet sites"\
      org.opencontainers.image.version="1.0.0"\
      org.opencontainers.image.authors="dmitriysyworov1986.com@gmail.com"\
      org.opencontainers.image.source="https://github.com/DmitriySyworov"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/parser-app ./cmd/main.go
RUN go build -o /app/migrator ./migration/auto.go
FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache \
    chromium \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont
COPY --from=builder /app/parser-app .
COPY --from=builder /app/migrator .
ENTRYPOINT ["./parser-app"]