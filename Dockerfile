# ---- Build stage ----
FROM golang:1.23 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o zenful_shopping ./cmd/api

# ---- Runtime stage ----
FROM python:3.11-slim as runtime

# Install additional system packages if needed
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy Go binary and Python scripts from builder stage
COPY --from=builder /app/zenful_shopping .
COPY --from=builder /app/scripts ./scripts
COPY --from=builder /app/scripts/requirements.txt ./scripts/requirements.txt

# Install Python dependencies
RUN pip install --upgrade pip && \
    pip install -r ./scripts/requirements.txt
# Expose application port
EXPOSE 8080
# Start the Go backend
CMD ["./zenful_shopping"]

