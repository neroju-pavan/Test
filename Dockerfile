# -----------------------------
# 1) Build Stage
# -----------------------------
FROM golang:1.25.4-alpine AS builder

# Install git for go mod
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod/go.sum first to cache dependencies
# Copy go.mod first
COPY go.mod ./

# Download modules (go.sum may not exist yet)
RUN go mod tidy


# Copy source code
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o userservice .

# -----------------------------
# 2) Runtime Stage
# -----------------------------
FROM alpine:3.19

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/userservice .

# Expose port
EXPOSE 8083

# Run the binary
CMD ["./userservice"]
