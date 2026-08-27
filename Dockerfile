ARG BUILD_FROM=ghcr.io/home-assistant/amd64-base:latest
FROM golang:1.22-alpine AS builder

WORKDIR /build
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o gundam-brain .

FROM $BUILD_FROM

WORKDIR /app
COPY --from=builder /build/gundam-brain /app/gundam-brain

# Install runtime utilities
RUN apk add --no-cache ca-certificates bash curl

EXPOSE 8080

ENTRYPOINT ["/app/gundam-brain"]
