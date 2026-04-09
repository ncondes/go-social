# ---- Build Stage ----
# Use the official Go image to compile the binary.
# This image is only used during build and is NOT included in the final image.
FROM golang:1.25 AS builder
WORKDIR /app

# Copy dependency files first to leverage Docker layer caching.
# go mod download only re-runs when go.mod or go.sum change.
COPY go.mod go.sum ./
RUN go mod download

# Install swag to generate Swagger docs before building.
# Required because cmd/api imports the generated docs package.
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy the rest of the source code.
COPY . .

# Generate Swagger documentation.
RUN swag init -g main.go -d ./cmd/api/,./internal/handlers,./internal/dtos,./internal/domain -o ./docs

# Build the binary. CGO_ENABLED=0 produces a fully static binary
# compatible with the minimal scratch image used in the run stage.
RUN CGO_ENABLED=0 GOOS=linux go build -o api cmd/api/*.go

# ---- Run Stage ----
# Use scratch (empty base image) for the smallest and most secure final image.
# Only the compiled binary and TLS certificates are included.
FROM scratch
WORKDIR /app

# Copy TLS certificates so the app can make outbound HTTPS calls (e.g. to GCP APIs).
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the compiled binary from the build stage.
COPY --from=builder /app/api .

# Cloud Run injects PORT env var and routes traffic to it.
# 8080 is the default expected port.
EXPOSE 8080

CMD ["./api"]
