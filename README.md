# Unibase DA Storage Node - Railway Deployment

This repository contains the configuration to deploy a Unibase DA storage node on Railway.

## Prerequisites

- Railway account
- Railway CLI installed (optional, for local deployment)

## Features

This storage node implementation:
1. Implements the Unibase DA storage API endpoints
2. Provides in-memory data storage (can be extended to persistent storage)
3. Supports upload/download operations
4. Compatible with Unibase SDK patterns
5. Ready for Railway deployment

## Quick Deploy

### Current Setup: Functional Storage Node

The repository now contains a fully functional storage node that implements the Unibase DA storage API. This node provides in-memory storage with all the required endpoints.

### API Endpoints

The storage node provides the following endpoints:

- `GET /health` - Health check endpoint
- `GET /api/info` - Node information
- `POST /api/upload` - Upload data
- `GET /api/download?name={id}` - Download data

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