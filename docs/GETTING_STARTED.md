# Getting Started with Flow

This guide will help you get Flow up and running on your local machine.

## Prerequisites

Before you begin, ensure you have the following installed:

- **Node.js** 18 or higher ([Download](https://nodejs.org/))
- **Go** 1.21 or higher ([Download](https://go.dev/dl/))
- **Git** ([Download](https://git-scm.com/downloads))

Optional:
- **Docker** and **Docker Compose** for containerized deployment
- **Make** for using the Makefile commands

## Quick Start (5 minutes)

### Option 1: Using Make (Recommended)

```bash
# Clone the repository
git clone https://github.com/yesoreyeram/flow.git
cd flow

# Install all dependencies
make install

# Run tests to ensure everything is working
make test

# Start development servers (frontend + backend)
make dev
```

The frontend will be available at http://localhost:3000 and backend at http://localhost:8080.

### Option 2: Manual Setup

#### Backend Setup

```bash
# Navigate to backend directory
cd backend

# Install dependencies
go mod download

# Run tests
go test ./...

# Start the server
go run cmd/server/main.go
```

The backend will start on http://localhost:8080.

#### Frontend Setup

```bash
# Open a new terminal and navigate to frontend directory
cd frontend

# Install dependencies
npm install

# Run tests
npm run test

# Start the development server
npm run dev
```

The frontend will start on http://localhost:3000.

### Option 3: Using Docker

```bash
# Clone the repository
git clone https://github.com/yesoreyeram/flow.git
cd flow

# Start with Docker Compose
docker-compose up --build
```

Access the application at http://localhost:8080.

## Creating Your First Workflow

1. **Open the application** in your browser at http://localhost:3000

2. **Click "New Workflow"** button

3. **Add nodes to your workflow:**
   - Click the "Add Node" button
   - Select a node type (HTTP Request, Transform, Condition)
   - Configure the node by clicking on it

4. **Connect nodes:**
   - Drag from the bottom handle of one node to the top handle of another

5. **Save your workflow:**
   - Click the "Save" button in the toolbar

6. **Execute your workflow:**
   - Click the "Execute" button
   - View the results in the execution panel

## Example Workflow: Fetch and Transform Data

Let's create a simple workflow that fetches data from an API and transforms it.

### Step 1: Add HTTP Request Node

1. Click "Add Node" → Select "httpRequest"
2. Click on the node to configure it
3. Set the following:
   - URL: `https://jsonplaceholder.typicode.com/users`
   - Method: `GET`
4. Click "Save"

### Step 2: Add Transform Node

1. Click "Add Node" → Select "transform"
2. Connect the HTTP Request node to this Transform node
3. Click on the Transform node to configure it
4. Set the following:
   - Language: `javascript`
   - Code:
   ```javascript
   return {
     count: input.length,
     names: input.map(user => user.name)
   };
   ```
5. Click "Save"

### Step 3: Save and Execute

1. Click "Save" in the toolbar
2. Click "Execute" to run the workflow
3. View the results showing the count and list of names

## Development Workflow

### Running Tests

**Frontend Tests:**
```bash
cd frontend

# Unit tests (watch mode)
npm run test

# Run once with coverage
npm run test:coverage

# E2E tests
npm run e2e

# E2E tests with UI
npm run e2e:ui
```

**Backend Tests:**
```bash
cd backend

# Run all tests
go test ./...

# With coverage
go test -cover ./...

# With race detection
go test -race ./...

# Verbose output
go test -v ./...
```

### Code Quality

**Linting:**
```bash
# Frontend
cd frontend
npm run lint

# Backend
cd backend
golangci-lint run
```

**Formatting:**
```bash
# Frontend
cd frontend
npm run format

# Backend
cd backend
go fmt ./...
```

### Building for Production

**Frontend:**
```bash
cd frontend
npm run build
# Output in frontend/dist/
```

**Backend:**
```bash
cd backend
go build -o server cmd/server/main.go
# Binary created as 'server'
```

**Docker:**
```bash
docker build -t flow:latest .
docker run -p 8080:8080 flow:latest
```

## Configuration

### Environment Variables

**Backend:**
- `PORT` - Server port (default: 8080)
- `ENVIRONMENT` - Environment (development/production)
- `CORS_ORIGINS` - Allowed CORS origins (default: http://localhost:3000)

**Frontend:**
Environment variables should be prefixed with `VITE_`:
- `VITE_API_URL` - Backend API URL (default: /api)

### Creating .env files

**Backend (.env):**
```env
PORT=8080
ENVIRONMENT=development
CORS_ORIGINS=http://localhost:3000
```

**Frontend (.env):**
```env
VITE_API_URL=/api
```

## Project Structure Overview

```
flow/
├── frontend/          # React frontend application
├── backend/           # Go backend application
├── docs/             # Documentation
├── .github/          # GitHub Actions CI/CD
├── Makefile          # Development commands
├── Dockerfile        # Production Docker image
└── docker-compose.yml # Development containers
```

## Common Issues

### Port Already in Use

If you get an error that port 3000 or 8080 is already in use:

```bash
# Find and kill the process using the port (macOS/Linux)
lsof -ti:3000 | xargs kill -9
lsof -ti:8080 | xargs kill -9

# Windows
netstat -ano | findstr :3000
taskkill /PID <PID> /F
```

### Module Not Found (Frontend)

```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### Go Module Issues (Backend)

```bash
cd backend
go mod tidy
go mod download
```

### CORS Issues

Ensure the backend `CORS_ORIGINS` environment variable includes your frontend URL.

## Next Steps

- Read the [Architecture Documentation](../docs/ARCHITECTURE.md)
- Explore the [API Documentation](../docs/API.md)
- Check out example workflows in the `examples/` directory (coming soon)
- Join our community (Discord/Slack links coming soon)

## Getting Help

- **Issues**: [GitHub Issues](https://github.com/yesoreyeram/flow/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yesoreyeram/flow/discussions)
- **Documentation**: [docs/](../docs/)

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

MIT License - see [LICENSE](../LICENSE) for details.
