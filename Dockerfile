# Frontend Dockerfile
FROM node:18-alpine AS frontend-build

WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm ci

COPY frontend/ ./
RUN npm run build

# Backend Dockerfile
FROM golang:1.21-alpine AS backend-build

WORKDIR /app/backend

COPY backend/go.* ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy backend binary
COPY --from=backend-build /server /app/server

# Copy frontend build
COPY --from=frontend-build /app/frontend/dist /app/public

# Expose ports
EXPOSE 8080

# Run
CMD ["/app/server"]
