# Security Best Practices

Flow is built with security as a top priority. This document outlines the security features and best practices.

## Security Features

### Input Validation

**Frontend:**
- All user inputs are validated using Zod schemas before submission
- Type-safe TypeScript ensures compile-time type checking
- React's built-in XSS protection through JSX

**Backend:**
- All API requests are validated
- Type checking via Go's strong typing system
- JSON parsing with error handling

### CORS Configuration

Cross-Origin Resource Sharing (CORS) is properly configured:

```go
// Configurable via environment variable
corsOrigins := []string{"http://localhost:3000"}

// Only allowed origins can access the API
w.Header().Set("Access-Control-Allow-Origin", origin)
```

**Production Setup:**
```bash
export CORS_ORIGINS="https://yourdomain.com,https://app.yourdomain.com"
```

### Error Handling

**Never expose internal errors:**
```go
// Bad
return fmt.Errorf("database connection failed: %v", err)

// Good
log.Printf("Internal error: %v", err)
return errors.New("internal server error")
```

**Client receives safe messages:**
```json
{
  "message": "Failed to create workflow",
  "code": "CREATION_ERROR"
}
```

### Request Timeouts

All operations have timeouts to prevent resource exhaustion:

```go
// HTTP server timeouts
srv := &http.Server{
    ReadTimeout:  15 * time.Second,
    WriteTimeout: 15 * time.Second,
    IdleTimeout:  60 * time.Second,
}

// HTTP client timeouts
httpClient := &http.Client{
    Timeout: 30 * time.Second,
}
```

### Panic Recovery

Middleware catches panics to prevent server crashes:

```go
defer func() {
    if err := recover(); err != nil {
        log.Printf("Panic recovered: %v", err)
        http.Error(w, "Internal server error", 500)
    }
}()
```

## Security Checklist

### Development

- [ ] Never commit secrets to version control
- [ ] Use environment variables for sensitive data
- [ ] Validate all user inputs
- [ ] Use parameterized queries (when adding database)
- [ ] Implement rate limiting
- [ ] Log security events
- [ ] Use HTTPS in production

### Deployment

- [ ] Set secure environment variables
- [ ] Enable HTTPS/TLS
- [ ] Configure firewall rules
- [ ] Set up monitoring and alerting
- [ ] Regular security updates
- [ ] Backup data regularly
- [ ] Implement authentication/authorization

### Code Review

- [ ] Check for SQL injection vulnerabilities
- [ ] Validate file upload handling
- [ ] Review authentication logic
- [ ] Check for XSS vulnerabilities
- [ ] Verify CSRF protection
- [ ] Review error messages
- [ ] Check for sensitive data in logs

## Authentication & Authorization

Currently, Flow does not include authentication. For production use, implement:

### JWT Authentication (Recommended)

**Backend:**
```go
// Add JWT middleware
func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        if token == "" {
            http.Error(w, "Unauthorized", 401)
            return
        }
        
        // Validate JWT token
        claims, err := validateJWT(token)
        if err != nil {
            http.Error(w, "Invalid token", 401)
            return
        }
        
        // Add user to context
        ctx := context.WithValue(r.Context(), "user", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

**Frontend:**
```typescript
// Add token to requests
axios.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});
```

### Role-Based Access Control (RBAC)

Implement permissions for different user roles:

```go
type Role string

const (
    RoleAdmin  Role = "admin"
    RoleEditor Role = "editor"
    RoleViewer Role = "viewer"
)

func requireRole(role Role) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            user := getUserFromContext(r.Context())
            if user.Role != role && user.Role != RoleAdmin {
                http.Error(w, "Forbidden", 403)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

## Rate Limiting

Implement rate limiting to prevent abuse:

```go
import "golang.org/x/time/rate"

type RateLimiter struct {
    limiters map[string]*rate.Limiter
    mu       sync.RWMutex
}

func rateLimitMiddleware(next http.Handler) http.Handler {
    rl := NewRateLimiter()
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := getIP(r)
        if !rl.Allow(ip) {
            http.Error(w, "Rate limit exceeded", 429)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

## SQL Injection Prevention

When adding database support, always use parameterized queries:

```go
// Bad - SQL Injection vulnerable
query := fmt.Sprintf("SELECT * FROM workflows WHERE id = '%s'", id)

// Good - Parameterized query
query := "SELECT * FROM workflows WHERE id = ?"
db.Query(query, id)
```

## XSS Prevention

Frontend XSS protection:

```typescript
// React automatically escapes JSX content
const userInput = "<script>alert('xss')</script>";
return <div>{userInput}</div>; // Safe - will render as text

// When using dangerouslySetInnerHTML, sanitize first
import DOMPurify from 'dompurify';
const clean = DOMPurify.sanitize(dirty);
return <div dangerouslySetInnerHTML={{ __html: clean }} />;
```

## HTTPS/TLS Configuration

Production should always use HTTPS:

```go
// Use Let's Encrypt
srv := &http.Server{
    Addr:      ":443",
    TLSConfig: &tls.Config{
        MinVersion: tls.VersionTLS12,
        CipherSuites: []uint16{
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
        },
    },
}
srv.ListenAndServeTLS("cert.pem", "key.pem")
```

## Secrets Management

Never hardcode secrets:

```go
// Bad
const apiKey = "sk-1234567890abcdef"

// Good
apiKey := os.Getenv("API_KEY")
if apiKey == "" {
    log.Fatal("API_KEY environment variable not set")
}
```

Use secret management tools in production:
- **AWS Secrets Manager**
- **HashiCorp Vault**
- **Azure Key Vault**
- **Google Secret Manager**

## Logging Security Events

Log security-relevant events:

```go
func logSecurityEvent(event string, details map[string]interface{}) {
    log.Printf("SECURITY: %s - %v", event, details)
}

// Usage
logSecurityEvent("failed_login", map[string]interface{}{
    "ip": getIP(r),
    "username": username,
    "timestamp": time.Now(),
})
```

## Security Headers

Add security headers in production:

```go
func securityHeadersMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Strict-Transport-Security", "max-age=31536000")
        w.Header().Set("Content-Security-Policy", "default-src 'self'")
        next.ServeHTTP(w, r)
    })
}
```

## Dependency Security

### Automated Scanning

GitHub Actions includes security scanning:

- **npm audit** - Frontend dependencies
- **Snyk** - Vulnerability scanning
- **gosec** - Go security checker
- **CodeQL** - Code analysis
- **Nancy** - Go dependency checker

### Manual Review

```bash
# Frontend
cd frontend
npm audit
npm audit fix

# Backend
cd backend
go list -json -deps ./... | nancy sleuth
```

## Vulnerability Reporting

If you discover a security vulnerability:

1. **DO NOT** open a public issue
2. Email security@yourdomain.com with details
3. Include steps to reproduce
4. Allow time for a fix before disclosure

## Security Updates

Stay updated with security patches:

```bash
# Frontend dependencies
cd frontend
npm update
npm audit fix

# Backend dependencies
cd backend
go get -u ./...
go mod tidy
```

## Production Deployment Checklist

- [ ] Use HTTPS/TLS certificates
- [ ] Set secure environment variables
- [ ] Enable authentication and authorization
- [ ] Implement rate limiting
- [ ] Configure CORS properly
- [ ] Set security headers
- [ ] Enable logging and monitoring
- [ ] Regular security audits
- [ ] Backup strategy in place
- [ ] Incident response plan
- [ ] Regular dependency updates
- [ ] Firewall and network security
- [ ] Use secrets management system

## Additional Resources

- [OWASP Top 10](https://owasp.org/www-project-top-ten/)
- [Go Security Best Practices](https://github.com/Checkmarx/Go-SCP)
- [React Security Best Practices](https://react.dev/learn/security)
- [Node.js Security Best Practices](https://nodejs.org/en/docs/guides/security/)

## Security Contact

For security concerns, contact: security@yourdomain.com
