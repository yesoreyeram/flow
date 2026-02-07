# API Reference

## Base URL

```
http://localhost:8080/api
```

## Authentication

Currently, the API does not require authentication. In production, add JWT or API key authentication.

## Workflows

### List Workflows

```http
GET /workflows
```

**Response:**
```json
[
  {
    "id": "wf-1234567890",
    "name": "My Workflow",
    "description": "A sample workflow",
    "nodes": [...],
    "edges": [...],
    "createdAt": "2024-01-01T00:00:00Z",
    "updatedAt": "2024-01-01T00:00:00Z",
    "version": 1
  }
]
```

### Get Workflow

```http
GET /workflows/:id
```

**Parameters:**
- `id` (string, required): Workflow ID

**Response:**
```json
{
  "id": "wf-1234567890",
  "name": "My Workflow",
  "description": "A sample workflow",
  "nodes": [
    {
      "id": "node-1",
      "type": "httpRequest",
      "position": {"x": 100, "y": 100},
      "data": {
        "label": "HTTP Request",
        "config": {
          "url": "https://api.example.com/data",
          "method": "GET"
        }
      }
    }
  ],
  "edges": [],
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z",
  "version": 1
}
```

### Create Workflow

```http
POST /workflows
```

**Request Body:**
```json
{
  "name": "New Workflow",
  "description": "Optional description",
  "nodes": [],
  "edges": []
}
```

**Response:**
```json
{
  "id": "wf-1234567890",
  "name": "New Workflow",
  "description": "Optional description",
  "nodes": [],
  "edges": [],
  "createdAt": "2024-01-01T00:00:00Z",
  "updatedAt": "2024-01-01T00:00:00Z",
  "version": 1
}
```

### Update Workflow

```http
PUT /workflows/:id
```

**Request Body:**
```json
{
  "name": "Updated Workflow",
  "description": "Updated description",
  "nodes": [...],
  "edges": [...]
}
```

**Response:** Same as Get Workflow

### Delete Workflow

```http
DELETE /workflows/:id
```

**Response:** 204 No Content

### Execute Workflow

```http
POST /workflows/:id/execute
```

**Request Body:**
```json
{
  "input": {
    "key": "value"
  }
}
```

**Response:**
```json
{
  "id": "exec-1234567890",
  "workflowId": "wf-1234567890",
  "status": "completed",
  "startedAt": "2024-01-01T00:00:00Z",
  "completedAt": "2024-01-01T00:00:01Z",
  "results": [
    {
      "nodeId": "node-1",
      "status": "success",
      "output": {
        "data": "result"
      },
      "duration": 100
    }
  ]
}
```

## Executions

### Get Execution

```http
GET /executions/:id
```

**Response:**
```json
{
  "id": "exec-1234567890",
  "workflowId": "wf-1234567890",
  "status": "completed",
  "startedAt": "2024-01-01T00:00:00Z",
  "completedAt": "2024-01-01T00:00:01Z",
  "results": [...]
}
```

### List Workflow Executions

```http
GET /workflows/:id/executions
```

**Response:**
```json
[
  {
    "id": "exec-1234567890",
    "workflowId": "wf-1234567890",
    "status": "completed",
    "startedAt": "2024-01-01T00:00:00Z",
    "completedAt": "2024-01-01T00:00:01Z",
    "results": [...]
  }
]
```

## Health Check

### Health

```http
GET /health
```

**Response:**
```json
{
  "status": "ok",
  "time": "2024-01-01T00:00:00Z"
}
```

## Error Responses

All errors follow this format:

```json
{
  "message": "Error description",
  "code": "ERROR_CODE"
}
```

### Common Error Codes

- `400` - Bad Request (Invalid input)
- `404` - Not Found (Resource doesn't exist)
- `500` - Internal Server Error

## Node Types

### HTTP Request Node

**Config Schema:**
```typescript
{
  url: string;           // Required
  method: "GET" | "POST" | "PUT" | "DELETE" | "PATCH";
  headers?: Record<string, string>;
  body?: any;
}
```

**Example:**
```json
{
  "url": "https://api.example.com/users",
  "method": "POST",
  "headers": {
    "Content-Type": "application/json",
    "Authorization": "Bearer token"
  },
  "body": {
    "name": "John Doe"
  }
}
```

### Transform Node

**Config Schema:**
```typescript
{
  code: string;          // Required
  language: "javascript" | "jq";
}
```

**Example:**
```json
{
  "code": "return input.map(x => ({ ...x, processed: true }))",
  "language": "javascript"
}
```

### Condition Node

**Config Schema:**
```typescript
{
  conditions: Array<{
    field: string;
    operator: "equals" | "notEquals" | "contains" | "greaterThan" | "lessThan";
    value: any;
  }>;
  combinator: "AND" | "OR";
}
```

**Example:**
```json
{
  "conditions": [
    {
      "field": "status",
      "operator": "equals",
      "value": "active"
    },
    {
      "field": "count",
      "operator": "greaterThan",
      "value": 10
    }
  ],
  "combinator": "AND"
}
```

## Rate Limits

Currently no rate limits. For production:
- 100 requests per minute per IP
- 1000 workflow executions per hour

## CORS

Configure allowed origins via `CORS_ORIGINS` environment variable.

Default: `http://localhost:3000`
