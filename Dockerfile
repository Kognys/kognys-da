# Mock storage node for testing Railway deployment
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy the mock server code
COPY mock-server.go .

RUN go mod init mock-storage-node
RUN go build -o store-edge mock-server.go

# Runtime stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates

# Create storage directory
RUN mkdir -p /root/.store

COPY --from=builder /app/store-edge /usr/local/bin/store-edge

# Set environment variables
ENV CHAIN_TYPE=bnb-testnet
ENV EXPOSE_PORT=8082

EXPOSE 8082

# Simple start command for the mock server
CMD ["/usr/local/bin/store-edge"]