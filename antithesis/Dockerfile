# Antithesis Dockerfile for otel-datagen
# Multi-stage build for optimal image size and Antithesis compatibility

#############################################################################
# Stage 1: Build the Go binary
#############################################################################
FROM golang:1.24-bookworm AS builder

WORKDIR /src

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary with static linking for compatibility
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o otel-datagen ./cmd/otel-datagen

#############################################################################
# Stage 2: Runtime image for Antithesis
#############################################################################
FROM debian:bookworm-slim

# Install basic utilities needed for testing
RUN apt-get update && apt-get install -y \
    ca-certificates \
    curl \
    bash \
    jq \
    && rm -rf /var/lib/apt/lists/*

# Create the Antithesis test directory structure
RUN mkdir -p /opt/antithesis/test/v1/exercise

# Copy the binary from builder stage
COPY --from=builder /src/otel-datagen /usr/local/bin/otel-datagen
RUN chmod +x /usr/local/bin/otel-datagen

# Copy container-adapted test scripts to Antithesis directory structure
COPY antithesis/docker/tests/simple_traces_test /opt/antithesis/test/v1/exercise/simple_traces_test
COPY antithesis/docker/tests/simple_logs_test /opt/antithesis/test/v1/exercise/simple_logs_test
COPY antithesis/docker/tests/simple_metrics_test /opt/antithesis/test/v1/exercise/simple_metrics_test
COPY antithesis/docker/tests/complex_traces_test /opt/antithesis/test/v1/exercise/complex_traces_test
COPY antithesis/docker/tests/complex_logs_test /opt/antithesis/test/v1/exercise/complex_logs_test
COPY antithesis/docker/tests/complex_metrics_test /opt/antithesis/test/v1/exercise/complex_metrics_test

# Make all test scripts executable
RUN chmod +x /opt/antithesis/test/v1/exercise/*

# Create entrypoint script
COPY <<'EOF' /opt/antithesis/entrypoint.sh
#!/bin/bash
set -euo pipefail

# Signal setup completion to Antithesis
echo '{"antithesis_setup": { "status": "complete", "details": null }}' > "${ANTITHESIS_OUTPUT_DIR:-/tmp/antithesis}/sdk.jsonl"

sleep infinity
EOF

RUN chmod +x /opt/antithesis/entrypoint.sh

# Set working directory to root for test execution
WORKDIR /

# Set the entrypoint
ENTRYPOINT ["/opt/antithesis/entrypoint.sh"]