# ----- Builder ----- #
FROM golang:1.25 AS builder

WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# build
# -ldflags="-w -s": Remove debug information to reduce binary size
# CGO_ENABLED=0: Create a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/server/main.go

# ----- runner ----- #
FROM alpine:latest

WORKDIR /app

# ca-certificates: Required for external HTTPS communication and TLS verification.
# tzdata: Required to display log times in the correct timezone (e.g., JST).
RUN apk --no-cache add ca-certificates tzdata

# Create a non-root user
# Avoid running as root to minimize risk if the container is compromised.
RUN adduser -D -g '' appuser

COPY --from=builder /app/server .

# Set permissions for volume mount points
# Ensure access rights to the certs directory mounted by docker-compose.
RUN mkdir -p /app/certs && \
    chown -R appuser:appuser /app/certs && \
    chown -R appuser:appuser /app

# Switch user
USER appuser

# Expose port and startup command
EXPOSE 8080

CMD ["./server"]
