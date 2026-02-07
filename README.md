# Flow - Workflow Automation Tool

A modern, enterprise-grade workflow automation platform similar to n8n, Zapier, and BuildShip. Built with React Flow, TypeScript, and Go.

## Features

- 🎨 Visual workflow editor with drag-and-drop interface
- 🔌 Multiple node types: HTTP requests, data transformations, conditions
- 🔄 Real-time workflow execution
- 🎯 Type-safe with TypeScript and Zod validation
- 🚀 Fast and efficient Go backend
- 🔒 Security-first design with comprehensive testing
- 🧪 Unit tests and E2E tests included
- 🐳 Docker support for easy deployment

## Tech Stack

### Frontend
- **React** 18 - UI library
- **React Flow** - Visual workflow editor
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling
- **Zod** - Runtime validation
- **Zustand** - State management
- **Vite** - Build tool
- **Vitest** - Unit testing
- **Playwright** - E2E testing

### Backend
- **Go** 1.21+ - Backend language
- **Standard library** - HTTP server (no framework needed)
- **Clean architecture** - Maintainable code structure
- In-memory storage (easily replaceable with PostgreSQL/MongoDB)

## Getting Started

### Prerequisites

- Node.js 18+
- Go 1.21+
- npm or yarn

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yesoreyeram/flow.git
cd flow
```

2. Install frontend dependencies:
```bash
cd frontend
npm install
```

3. Install backend dependencies:
```bash
cd backend
go mod download
```

### Development

Run frontend and backend separately:

**Frontend:**
```bash
cd frontend
npm run dev
```
Frontend will be available at http://localhost:3000

**Backend:**
```bash
cd backend
go run cmd/server/main.go
```
Backend API will be available at http://localhost:8080

### Using Docker

Build and run with Docker Compose:

```bash
docker-compose up --build
```

## Testing

### Frontend Tests

```bash
cd frontend

# Unit tests
npm run test

# E2E tests
npm run e2e

# Test with coverage
npm run test:coverage
```

### Backend Tests

```bash
cd backend

# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run with race detection
go test -race ./...
```

## Building for Production

### Frontend

```bash
cd frontend
npm run build
```

Build output will be in `frontend/dist/`

### Backend

```bash
cd backend
go build -o server cmd/server/main.go
```

Binary will be created as `server`

### Docker Build

```bash
docker build -t flow:latest .
docker run -p 8080:8080 flow:latest
```

## Project Structure

```
flow/
├── frontend/               # React frontend
│   ├── src/
│   │   ├── components/    # React components
│   │   │   ├── editor/    # Workflow editor components
│   │   │   ├── nodes/     # Custom node components
│   │   │   └── ui/        # Reusable UI components
│   │   ├── pages/         # Page components
│   │   ├── stores/        # Zustand state management
│   │   ├── types/         # TypeScript types & Zod schemas
│   │   ├── services/      # API services
│   │   └── test/          # Tests
│   └── e2e/               # E2E tests
│
├── backend/               # Go backend
│   ├── cmd/server/        # Main application entry
│   ├── internal/
│   │   ├── api/           # HTTP handlers & middleware
│   │   ├── engine/        # Workflow execution engine
│   │   ├── models/        # Data models
│   │   ├── repository/    # Data access layer
│   │   └── config/        # Configuration
│   └── pkg/               # Shared packages
│
├── .github/workflows/     # GitHub Actions CI/CD
├── docs/                  # Documentation
├── Dockerfile            # Production Dockerfile
└── docker-compose.yml    # Docker Compose config
```

## API Documentation

### Workflows

- `GET /api/workflows` - List all workflows
- `POST /api/workflows` - Create a workflow
- `GET /api/workflows/:id` - Get a workflow
- `PUT /api/workflows/:id` - Update a workflow
- `DELETE /api/workflows/:id` - Delete a workflow
- `POST /api/workflows/:id/execute` - Execute a workflow

### Executions

- `GET /api/executions/:id` - Get execution details
- `GET /api/workflows/:id/executions` - List workflow executions

### Health

- `GET /api/health` - Health check endpoint

## Node Types

### HTTP Request Node
Make HTTP requests to external APIs with support for:
- Multiple HTTP methods (GET, POST, PUT, DELETE, PATCH)
- Custom headers
- Request body
- Authentication (Basic, Bearer, API Key)

### Transform Node
Transform data using:
- JavaScript code
- jq expressions

### Condition Node
Conditional logic with:
- Multiple conditions
- AND/OR combinators
- Various operators (equals, contains, greaterThan, etc.)

### Trigger Node
Trigger workflows based on:
- Webhooks
- Scheduled events
- Manual triggers

## Security

- Input validation using Zod schemas
- CORS configuration
- Request timeouts
- Rate limiting ready
- Error recovery middleware
- Secure headers
- No hardcoded secrets

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run linters and tests
6. Submit a pull request

## License

MIT

## Support

For issues and questions, please use the GitHub issue tracker.