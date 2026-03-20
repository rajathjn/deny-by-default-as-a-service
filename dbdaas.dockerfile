FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# Disable CGO for static compilation and build the Go binary
# This results in a binary without external runtime dependencies
RUN CGO_ENABLED=0 go build -o dbdaas .

# --- Final stage ---
FROM scratch

WORKDIR /app

# Copy the built binary from the 'builder' stage to the final image
COPY --from=builder /app/dbdaas .

# Expose the port application listens
EXPOSE 8080

# Command to run the executable when the container starts
CMD ["./dbdaas"]
