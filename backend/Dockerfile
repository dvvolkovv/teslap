# Multi-stage build for all TeslaPay Go services.
# Build a specific service with: docker build --build-arg SERVICE=gateway .
# Available services: gateway, auth-service, account-service, payment-service,
#                     card-service, ledger-service, kyc-service, notification-service

# ---- Stage 1: Builder ----
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git for private module fetches if needed.
RUN apk add --no-cache git

# Cache dependency downloads separately from source builds.
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source tree.
COPY . .

# SERVICE build arg selects which cmd/ sub-directory to compile.
ARG SERVICE
RUN test -n "$SERVICE" || (echo "ERROR: ARG SERVICE is required" && exit 1)

# Build a statically linked binary; CGO disabled for minimal Alpine image.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -trimpath -ldflags="-s -w" \
    -o /service ./cmd/${SERVICE}/

# ---- Stage 2: Runtime ----
FROM alpine:3.19

# ca-certificates: needed for HTTPS outbound calls (Kafka TLS, OTEL, etc.)
# tzdata: needed so time.LoadLocation works inside the container.
RUN apk add --no-cache ca-certificates tzdata

# Run as a non-root user.
RUN addgroup -S teslapay && adduser -S teslapay -G teslapay

# Copy the compiled binary from the builder stage.
COPY --from=builder /service /service

RUN chmod +x /service

USER teslapay

# The actual port is controlled by the PORT env var; this EXPOSE is documentation.
EXPOSE 8080

CMD ["/service"]
