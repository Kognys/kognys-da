#!/bin/sh

# Create directories if they don't exist
mkdir -p ~/.store/datastore
mkdir -p ~/.store/keystore

# Create a basic config if it doesn't exist
if [ ! -f ~/.store/config.json ]; then
    cat > ~/.store/config.json << EOF
{
  "wallet": {
    "address": "0x0000000000000000000000000000000000000000"
  },
  "chain": {
    "type": "${CHAIN_TYPE:-bnb-testnet}"
  },
  "api": {
    "endpoint": "0.0.0.0:${EXPOSE_PORT:-8082}",
    "expose": "${EXPOSE_URL:-http://localhost:8082}"
  },
  "remote": {
    "url": "https://api.unibase.io"
  }
}
EOF
    echo "Created default config"
fi

echo "Storage node initialized"