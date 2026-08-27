ARG BUILD_FROM=ubuntu:24.04
FROM golang:1.22 AS builder

WORKDIR /build
COPY go.mod ./
COPY main.go ./
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o gundam-brain .

FROM $BUILD_FROM

ENV DEBIAN_FRONTEND=noninteractive
ENV PATH="/root/.local/bin:/usr/local/bin:${PATH}"

# Install system dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        bash \
        git \
        jq && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install Antigravity CLI (agy)
RUN curl -fsSL https://antigravity.google/cli/install.sh | bash || true

# Ensure binary is in standard path
RUN if [ -f /root/.local/bin/agy ]; then ln -sf /root/.local/bin/agy /usr/local/bin/agy; fi

WORKDIR /app
COPY --from=builder /build/gundam-brain /app/gundam-brain

EXPOSE 8080

ENTRYPOINT ["/app/gundam-brain"]
