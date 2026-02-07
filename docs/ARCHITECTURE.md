# Architecture Overview

## System Architecture

Flow follows a clean, modular architecture with clear separation between frontend and backend.

```
┌─────────────────────────────────────────────────────────────┐
│                         Frontend                             │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐     │
│  │   React UI   │  │  React Flow  │  │   Zustand    │     │
│  │  Components  │  │    Editor    │  │    Store     │     │
│  └──────────────┘  └──────────────┘  └──────────────┘     │
│         │                  │                  │             │
│         └──────────────────┴──────────────────┘             │
│                          │                                   │
│                    ┌─────▼─────┐                           │
│                    │  API       │                           │
│                    │  Service   │                           │
│                    └─────┬─────┘                           │
└──────────────────────────┼───────────────────────────────┘
                           │ REST API
┌──────────────────────────▼───────────────────────────────┐
│                         Backend                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │
│  │   HTTP       │  │   Workflow   │  │  Repository  │   │
│  │   Handlers   │  │   Engine     │  │  (Storage)   │   │
│  └──────────────┘  └──────────────┘  └──────────────┘   │
│         │                  │                  │           │
│         └──────────────────┴──────────────────┘           │
└─────────────────────────────────────────────────────────┘
```

## Frontend Architecture

### Component Structure

```
src/
├── components/
│   ├── Layout.tsx          # Main layout wrapper
│   ├── nodes/              # React Flow custom nodes
│   │   └── CustomNode.tsx
│   ├── editor/             # Workflow editor components
│   │   └── NodeConfigPanel.tsx
│   └── ui/                 # Reusable UI components
│       ├── Button.tsx
│       ├── Input.tsx
│       └── Select.tsx
├── pages/
│   ├── WorkflowList.tsx    # List all workflows
│   └── WorkflowEditor.tsx  # Edit workflow
├── stores/
│   └── workflowStore.ts    # Global state management
├── types/
│   └── workflow.ts         # TypeScript types & Zod schemas
└── services/
    └── api.ts              # API client
```

### State Management

- **Zustand** for global state
- **React Query** for server state
- Local state for component-specific data

### Data Flow

1. User interacts with React Flow canvas
2. State updates in Zustand store
3. API service makes HTTP calls
4. Backend processes request
5. Response updates UI state

## Backend Architecture

### Project Structure

```
backend/
├── cmd/
│   └── server/
│       └── main.go         # Entry point
├── internal/
│   ├── api/
│   │   ├── server.go       # HTTP server setup
│   │   └── middleware.go   # Middleware functions
│   ├── engine/
│   │   └── engine.go       # Workflow execution engine
│   ├── models/
│   │   └── workflow.go     # Domain models
│   ├── repository/
│   │   └── repository.go   # Data access layer
│   └── config/
│       └── config.go       # Configuration
└── pkg/                    # Shared packages
```

### Clean Architecture Layers

1. **Presentation Layer** (api/)
   - HTTP handlers
   - Request/response formatting
   - Middleware

2. **Business Logic Layer** (engine/)
   - Workflow execution
   - Node processing
   - Data transformation

3. **Data Layer** (repository/)
   - Data persistence
   - CRUD operations
   - Database abstraction

### Design Patterns

#### Repository Pattern
Abstracts data access, making it easy to swap storage backends:

```go
type Repository interface {
    GetWorkflow(id string) (*Workflow, error)
    CreateWorkflow(workflow *Workflow) error
    // ... more methods
}
```

#### Dependency Injection
Dependencies are injected through constructors:

```go
func NewServer(cfg *Config, repo Repository) *Server {
    return &Server{
        config: cfg,
        repo:   repo,
    }
}
```

#### Middleware Pattern
Chain middleware for cross-cutting concerns:

```go
handler = corsMiddleware(origins)(handler)
handler = loggingMiddleware(handler)
handler = recoveryMiddleware(handler)
```

## Workflow Execution

### Execution Flow

1. **Request Received**: POST to `/api/workflows/:id/execute`
2. **Load Workflow**: Fetch workflow from repository
3. **Build Execution Graph**: Determine node execution order
4. **Execute Nodes**: Process each node sequentially
5. **Store Results**: Save execution results
6. **Return Response**: Send execution status to client

### Node Execution

Each node type has its own execution logic:

```go
func (e *Engine) executeNode(ctx context.Context, node *Node) (output, error) {
    switch node.Type {
    case HTTPRequest:
        return e.executeHTTPRequest(ctx, node)
    case Transform:
        return e.executeTransform(ctx, node)
    case Condition:
        return e.executeCondition(ctx, node)
    }
}
```

### Error Handling

- Errors are propagated up the chain
- Partial results are saved even on failure
- Context for cancellation and timeouts

## Security Architecture

### Frontend Security

- **Input Validation**: Zod schemas validate all user input
- **XSS Prevention**: React automatically escapes content
- **CORS**: Configured in backend
- **No Secrets**: All sensitive data in backend

### Backend Security

- **Request Validation**: Validate all incoming requests
- **Error Handling**: Don't expose internal errors
- **Rate Limiting**: Ready for implementation
- **CORS Middleware**: Restrict origins
- **Timeouts**: Prevent long-running requests
- **Recovery Middleware**: Graceful panic recovery

## Scalability

### Current State
- In-memory storage (single instance)
- Synchronous execution

### Future Enhancements
- PostgreSQL/MongoDB for persistence
- Redis for caching
- Message queue for async execution
- Horizontal scaling with load balancer
- Distributed execution workers

## Testing Strategy

### Frontend Tests
- **Unit Tests**: Component logic (Vitest)
- **Integration Tests**: Component interactions
- **E2E Tests**: User flows (Playwright)

### Backend Tests
- **Unit Tests**: Individual functions
- **Integration Tests**: API endpoints
- **Load Tests**: Performance testing

## CI/CD Pipeline

1. **Code Push**: Trigger GitHub Actions
2. **Lint**: Check code style
3. **Test**: Run all tests
4. **Build**: Create artifacts
5. **Security Scan**: Check for vulnerabilities
6. **Deploy**: (Ready for deployment setup)

## Monitoring & Observability

Ready for implementation:
- Structured logging
- Metrics collection
- Error tracking
- Performance monitoring
- Distributed tracing
