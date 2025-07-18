# Build a simple storage node that uses Unibase SDK patterns
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy the storage node source
COPY store-edge.go .

RUN go mod init storage-node
RUN go build -o store-edge store-edge.go

# Runtime stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates

# Create storage directory
RUN mkdir -p /root/.store

COPY --from=builder /app/store-edge /usr/local/bin/

# Set environment variables
ENV CHAIN_TYPE=bnb-testnet
ENV EXPOSE_PORT=8082

EXPOSE 8082

# Direct run (no init needed for this simple version)
CMD ["/usr/local/bin/store-edge"]