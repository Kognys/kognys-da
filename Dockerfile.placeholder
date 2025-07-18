# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates curl

# Create storage directory
RUN mkdir -p /root/.store

# Create a directory for the binary
RUN mkdir -p /usr/local/bin

# NOTE: You need to add the store-edge binary to this image
# Option 1: Download from a release URL (update URL when available)
# RUN curl -L https://github.com/unibaseio/store-edge/releases/download/vX.X.X/store-edge-linux-amd64 -o /usr/local/bin/store-edge && chmod +x /usr/local/bin/store-edge

# Option 2: Copy from local file (uncomment and use this if you have the binary)
# COPY store-edge /usr/local/bin/store-edge
# RUN chmod +x /usr/local/bin/store-edge

# Option 3: Build from source (when repository becomes available)
# Uncomment the builder stage above and use COPY --from=builder

# Temporary: Create a placeholder script
RUN echo '#!/bin/sh\necho "Error: store-edge binary not found. Please add the binary to the Docker image."\necho "See Dockerfile comments for instructions."\nexit 1' > /usr/local/bin/store-edge && \
    chmod +x /usr/local/bin/store-edge

# Set environment variables
ENV CHAIN_TYPE=bnb-testnet
ENV EXPOSE_PORT=8082

# Expose the storage node port
EXPOSE 8082

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:${EXPOSE_PORT}/health || exit 1

# Initialize and run the storage node
CMD ["sh", "-c", "store-edge init && store-edge daemon run -b 0.0.0.0:${EXPOSE_PORT} -e ${EXPOSE_URL}"]