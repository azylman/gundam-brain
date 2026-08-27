FROM golang:1.22-alpine AS builder

WORKDIR /build
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o gundam-brain .

FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive
ENV PATH="/usr/local/bin:${PATH}"

# Install system dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        bash \
        git \
        jq \
        tar \
        gzip && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install Antigravity CLI (agy) directly into /usr/local/bin
RUN curl -fsSL https://antigravity.google/cli/install.sh | bash -s -- -d /usr/local/bin

WORKDIR /app
COPY --from=builder /build/gundam-brain /app/gundam-brain

EXPOSE 8080

ENTRYPOINT ["/app/gundam-brain"]
