# --- Build stage --------------------------------------------------------
FROM golang:1.22-alpine AS builder

WORKDIR /build

COPY go.mod ./
COPY main.go ./
RUN go mod tidy
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /gundam-brain .

# --- Runtime stage -------------------------------------------------------
FROM ubuntu:24.04

# Install base utilities, certificates, and runtime interpreters for MCP servers (Node.js, Python)
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    bash \
    git \
    jq \
    tar \
    gzip \
    python3 \
    python3-pip \
    python3-venv \
    nodejs \
    npm \
    && rm -rf /var/lib/apt/lists/*

# Pre-install common MCP servers into global npm cache for instant stdio startup
RUN npm install -g --no-audit --no-fund \
    @modelcontextprotocol/server-github \
    @pasympa/discord-mcp \
    @iqai/mcp-discord

# Install official Antigravity CLI binary
RUN curl -fsSL https://antigravity.google/cli/install.sh | bash -s -- -d /usr/local/bin

WORKDIR /app

COPY --from=builder /gundam-brain /usr/local/bin/gundam-brain

EXPOSE 8080

CMD ["/usr/local/bin/gundam-brain"]
