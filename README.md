# Unibase DA Storage Node - Railway Deployment

This repository contains the configuration to deploy a Unibase DA storage node on Railway.

## Prerequisites

- Railway account
- Railway CLI installed (optional, for local deployment)
- Access to the store-edge binary (see instructions below)

## Important Note

The `store-edge` binary is not publicly available. This repository provides:
1. A placeholder Dockerfile that you need to update with the actual binary
2. A mock server for testing the deployment setup

## Quick Deploy

### Option 1: Deploy Mock Server (For Testing)

1. Rename `Dockerfile.test` to `Dockerfile`
2. Push to your GitHub repository
3. Connect to Railway and deploy

### Option 2: Deploy Real Storage Node

1. Obtain the `store-edge` binary from Unibase
2. Update the `Dockerfile` with one of these methods:
   - Add the binary URL to download during build
   - Copy the binary to the repository and uncomment the COPY line
   - Add build instructions if you have access to source code
3. Push to GitHub and deploy via Railway

### Manual Deployment

1. Fork or clone this repository
2. Update the Dockerfile as needed
3. Connect your GitHub repository to Railway
4. Railway will automatically detect the Dockerfile and deploy

## Configuration

### Environment Variables

Set these in your Railway project settings:

- `CHAIN_TYPE`: Network type (default: `bnb-testnet`)
- `EXPOSE_PORT`: Port for the storage node (default: `8082`)
- `EXPOSE_URL`: External URL (Railway sets this automatically)
- `SK`: (Optional) Your secret key for authenticated operations

### Railway Configuration

The deployment uses:
- `Dockerfile`: Builds the storage node from source
- `railway.json`: Railway-specific configuration
- `railway.toml`: Service and environment configuration

## Architecture

The storage node:
1. Builds from the official Unibase store-edge repository
2. Initializes storage in `/root/.store`
3. Exposes port 8082 for external access
4. Automatically configures the external URL using Railway's domain

## Usage

Once deployed, your storage node will be accessible at:
```
https://your-app-name.railway.app
```

### Upload Example
```bash
curl -X POST https://your-app-name.railway.app/api/upload \
  -d '{"id":"test1","owner":"0xabcd","message":"sample message"}'
```

### Download Example
```bash
wget https://your-app-name.railway.app/api/download?name=<file_name>&owner=<owner_address>
```

## Monitoring

Railway provides built-in monitoring for:
- CPU and memory usage
- Network traffic
- Application logs

Access these through your Railway dashboard.

## Troubleshooting

1. **Storage node not starting**: Check logs in Railway dashboard
2. **Connection issues**: Ensure `EXPOSE_URL` is set correctly
3. **Out of memory**: Upgrade your Railway plan for more resources

## Resources

- [Unibase Documentation](https://openos-labs.gitbook.io/unibase-docs/)
- [Railway Documentation](https://docs.railway.app/)
- [Unibase GitHub](https://github.com/unibaseio)