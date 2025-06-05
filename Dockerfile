# ---- Build stage ----
FROM golang:1.23 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o zenful_shopping ./cmd/api

# ---- Runtime stage ----
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/zenful_shopping .
EXPOSE 8080
CMD ["./zenful_shopping"]

