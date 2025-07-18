# Unibase DA Storage Node API Documentation

**Base URL**: `https://kognys-da-production.up.railway.app`

## Endpoints

### 1. Health Check

Check if the storage node is running and healthy.

**Endpoint**: `GET /health`

**Request**:
```bash
curl https://kognys-da-production.up.railway.app/health
```

**Response**:
```json
{
  "chain_type": "bnb-testnet",
  "port": "8082",
  "status": "healthy",
  "timestamp": "2025-07-18T13:57:32.178634969Z",
  "type": "storage-node"
}
```

**Status Codes**:
- `200 OK`: Node is healthy

---

### 2. Node Information

Get information about the storage node.

**Endpoint**: `GET /api/info`

**Request**:
```bash
curl https://kognys-da-production.up.railway.app/api/info
```

**Response**:
```json
{
  "chainType": "bnb-testnet",
  "exposeURL": "",
  "name": "unibase-storage-node",
  "type": "store"
}
```

**Status Codes**:
- `200 OK`: Success

---

### 3. Upload Data

Upload data to the storage node.

**Endpoint**: `POST /api/upload`

**Headers**:
- `Content-Type: application/json`

**Request Body**:
```json
{
  "id": "unique-identifier",
  "owner": "0xOwnerAddress",
  "message": "Your data here",
  "timestamp": "2025-07-18T13:57:00Z",
  "anyField": "Any additional data"
}
```

**Example Request**:
```bash
curl -X POST https://kognys-da-production.up.railway.app/api/upload \
  -H "Content-Type: application/json" \
  -d '{
    "id": "railway-test-001",
    "owner": "0xtest",
    "message": "Hello from Railway deployment",
    "timestamp": "2025-07-18T13:57:00Z"
  }'
```

**Response**:
```json
{
  "id": "railway-test-001",
  "message": "Upload successful",
  "success": true
}
```

**Status Codes**:
- `200 OK`: Upload successful
- `400 Bad Request`: Invalid JSON or missing required fields
- `405 Method Not Allowed`: Wrong HTTP method

**Notes**:
- The `id` field is used as the key for retrieval. If not provided, a timestamp-based ID will be generated
- You can include any JSON-serializable data in the request body
- All fields are stored and returned on download

---

### 4. Download Data

Retrieve previously uploaded data.

**Endpoint**: `GET /api/download`

**Query Parameters**:
- `name` (required): The ID of the data to retrieve
- `owner` (optional): The owner address (for logging purposes)

**Example Request**:
```bash
curl "https://kognys-da-production.up.railway.app/api/download?name=railway-test-001&owner=0xtest"
```

**Response**:
```json
{
  "data": {
    "id": "railway-test-001",
    "message": "Hello from Railway deployment",
    "owner": "0xtest",
    "timestamp": "2025-07-18T13:57:00Z"
  },
  "message": "Download successful",
  "name": "railway-test-001",
  "owner": "0xtest"
}
```

**Status Codes**:
- `200 OK`: Data found and returned
- `400 Bad Request`: Missing `name` parameter
- `404 Not Found`: Data with specified name not found

---

## Examples

### Uploading AI Model Data

```bash
curl -X POST https://kognys-da-production.up.railway.app/api/upload \
  -H "Content-Type: application/json" \
  -d '{
    "id": "model-v1.0",
    "type": "ai-model",
    "data": {
      "weights": [0.1, 0.2, 0.3, 0.4, 0.5],
      "version": "1.0",
      "architecture": "transformer"
    },
    "metadata": {
      "created": "2025-07-18",
      "author": "research-team"
    }
  }'
```

### Uploading DePIN Device Data

```bash
curl -X POST https://kognys-da-production.up.railway.app/api/upload \
  -H "Content-Type: application/json" \
  -d '{
    "id": "device-001-data",
    "deviceId": "iot-sensor-001",
    "readings": {
      "temperature": 23.5,
      "humidity": 45,
      "pressure": 1013.25
    },
    "timestamp": "2025-07-18T14:00:00Z",
    "location": {
      "lat": 40.7128,
      "lon": -74.0060
    }
  }'
```

### Batch Upload Example

```bash
curl -X POST https://kognys-da-production.up.railway.app/api/upload \
  -H "Content-Type: application/json" \
  -d '{
    "id": "batch-001",
    "type": "batch-data",
    "items": [
      {"sensor": "A1", "value": 100},
      {"sensor": "A2", "value": 150},
      {"sensor": "A3", "value": 125}
    ],
    "processedAt": "2025-07-18T14:00:00Z"
  }'
```

## CORS Support

The API supports Cross-Origin Resource Sharing (CORS) with the following configuration:
- **Allowed Origins**: * (all origins)
- **Allowed Methods**: GET, POST, OPTIONS
- **Allowed Headers**: Content-Type

This allows the API to be accessed from web browsers and frontend applications.

## Rate Limits

Currently, there are no rate limits implemented. For production use, consider implementing rate limiting based on your requirements.

## Data Persistence

**Important**: The current implementation stores data in memory. This means:
- Data will be lost if the server restarts
- Storage is limited by available RAM
- For production use, consider implementing persistent storage (database, file system, or distributed storage)

## Error Handling

All endpoints return JSON-formatted error messages:

```json
{
  "error": "Description of the error"
}
```

Common HTTP status codes:
- `200`: Success
- `400`: Bad Request (invalid input)
- `404`: Not Found (resource doesn't exist)
- `405`: Method Not Allowed
- `500`: Internal Server Error