# 🔐 JWT Debugger - jwt.io Alternative in Go

A production-ready JWT (JSON Web Token) debugger and validator web application built with Go, inspired by jwt.io. This tool allows you to decode, verify, and generate JWTs with a clean, modern web interface.

[![Go Version](https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat&logo=docker)](https://www.docker.com)
[![License](https://img.shields.io/badge/License-Educational-green.svg)](LICENSE)

> 📖 **Docker Deployment Guide:** See [DOCKER.md](DOCKER.md) for comprehensive Docker and container orchestration documentation.

## ✨ Features

### 🔍 JWT Decoder
- **Real-time Decoding**: Automatically decodes JWT tokens as you type
- **Visual Breakdown**: Displays Header, Payload, and Signature in separate color-coded sections
- **Claims Detection**: Automatically identifies and formats standard JWT claims:
  - `exp` (Expiration Time)
  - `iat` (Issued At)
  - `nbf` (Not Before)
  - `iss` (Issuer)
  - `aud` (Audience)
  - `sub` (Subject)
- **Expiration Highlighting**: Visually indicates if a token is expired
- **Human-Readable Timestamps**: Converts Unix timestamps to readable dates

### 🔨 JWT Generator
- **Algorithm Support**: Generate JWTs with:
  - HS256 (HMAC SHA-256)
  - HS384 (HMAC SHA-384)
  - HS512 (HMAC SHA-512)
- **Custom Headers & Payload**: Full control over token contents
- **Copy to Clipboard**: One-click copy for generated tokens

### ✅ JWT Verifier
- **Signature Verification**: Validate JWT signatures with secret key
- **Detailed Error Messages**: Clear feedback on validation failures:
  - Signature mismatch
  - Token expired
  - Not valid yet (nbf)
  - Malformed tokens
- **Claims Display**: Shows verified claims after successful validation

## 📸 Screenshots

### JWT Decoder View
![JWT Viewer - Decode](doc/jwt-viewer.png)

### JWT Encoder View
![JWT Viewer - Encode](doc/jwt-viewer-encode.png)

## 🏗️ Architecture

Clean architecture with separated concerns:

```
jwt-viewer/
├── main.go              # Application entry point
├── go.mod               # Go module dependencies
├── go.sum               # Go module checksums
├── Dockerfile           # Docker image definition
├── docker-compose.yml   # Docker Compose configuration
├── .gitignore           # Git ignore rules
├── README.md            # Project documentation
├── handlers/            # HTTP request handlers
│   └── jwt_handler.go
├── services/            # Business logic layer
│   └── jwt_service.go
├── models/              # Data structures
│   └── jwt.go
├── utils/               # Utility functions
│   └── jwt_utils.go
└── static/              # Frontend assets
    ├── index.html
    ├── styles.css
    └── app.js
```

## 🚀 Quick Start

### Prerequisites

#### Option 1: Run with Go
- Go 1.20 or higher
- Git (optional)

#### Option 2: Run with Docker
- Docker
- Docker Compose (optional)

### Installation & Running

#### Method 1: Run Locally with Go

1. **Navigate to the project directory**:
   ```bash
   cd D:\project\jwt-viewer
   ```

2. **Install dependencies** (if not already done):
   ```bash
   go mod download
   ```

3. **Run the application**:
   ```bash
   # Default port 8080
   go run main.go
   
   # Or with custom port
   PORT=9000 go run main.go  # Linux/Mac
   $env:PORT="9000"; go run main.go  # Windows PowerShell
   
   # With release mode
   GIN_MODE=release go run main.go  # Linux/Mac
   $env:GIN_MODE="release"; go run main.go  # Windows PowerShell
   ```

4. **Open your browser**:
   ```
   http://localhost:8080
   ```

The server will start on port 8080 by default.

#### Method 2: Run with Docker Compose (Recommended)

1. **Navigate to the project directory**:
   ```bash
   cd D:\project\jwt-viewer
   ```

2. **Start the application**:
   ```bash
   docker-compose up -d
   ```

3. **Open your browser**:
   ```
   http://localhost:8080
   ```

4. **View logs** (optional):
   ```bash
   docker-compose logs -f jwt-viewer
   ```

5. **Stop the application**:
   ```bash
   docker-compose down
   ```

#### Method 3: Run with Docker (without Compose)

1. **Build the Docker image**:
   ```bash
   docker build -t jwt-viewer:latest .
   ```

2. **Run the container**:
   ```bash
   docker run -d -p 8080:8080 --name jwt-viewer jwt-viewer:latest
   ```

3. **Open your browser**:
   ```
   http://localhost:8080
   ```

4. **Stop the container**:
   ```bash
   docker stop jwt-viewer
   docker rm jwt-viewer
   ```

## 📡 API Endpoints

### POST /api/decode
Decode a JWT without signature verification.

**Request:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response:**
```json
{
  "header": {
    "alg": "HS256",
    "typ": "JWT"
  },
  "payload": {
    "sub": "1234567890",
    "name": "John Doe",
    "iat": 1516239022
  },
  "signature": "SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
  "claim_info": {
    "iat": "2018-01-18 01:30:22 UTC",
    "sub": "1234567890",
    "is_expired": false
  }
}
```

### POST /api/encode
Generate a new signed JWT.

**Request:**
```json
{
  "header": {
    "typ": "JWT"
  },
  "payload": {
    "sub": "1234567890",
    "name": "John Doe",
    "iat": 1516239022
  },
  "secret": "your-256-bit-secret",
  "algorithm": "HS256"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### POST /api/verify
Verify a JWT signature.

**Request:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "secret": "your-256-bit-secret"
}
```

**Response:**
```json
{
  "valid": true,
  "message": "Token is valid",
  "claims": {
    "sub": "1234567890",
    "name": "John Doe",
    "iat": 1516239022
  },
  "claim_info": {
    "iat": "2018-01-18 01:30:22 UTC",
    "sub": "1234567890",
    "is_expired": false
  }
}
```

## ⚙️ Configuration

### Environment Variables

The application can be configured using environment variables:

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `PORT` | Server port | `8080` | `9000` |
| `GIN_MODE` | Application mode | `debug` | `release` |
| `HOST_PORT` | Docker host port mapping | `8080` | `8181` |

**Usage Examples:**

```bash
# Linux/macOS
export PORT=9000
export GIN_MODE=release
go run main.go

# Windows PowerShell
$env:PORT="9000"
$env:GIN_MODE="release"
go run main.go

# Docker Compose with custom port
HOST_PORT=8181 docker-compose up -d

# Or create a .env file
cp .env.example .env
# Edit .env file with your values
docker-compose up -d
```

**.env File Example:**
```env
PORT=8080
HOST_PORT=8080
GIN_MODE=release
```

## 🧪 Example Usage

### Using cURL

**Decode a token:**
```bash
curl -X POST http://localhost:8080/api/decode \
  -H "Content-Type: application/json" \
  -d '{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"}'
```

**Generate a token:**
```bash
curl -X POST http://localhost:8080/api/encode \
  -H "Content-Type: application/json" \
  -d '{
    "header": {"typ": "JWT"},
    "payload": {"sub": "1234567890", "name": "John Doe", "iat": 1516239022},
    "secret": "your-256-bit-secret",
    "algorithm": "HS256"
  }'
```

**Verify a token:**
```bash
curl -X POST http://localhost:8080/api/verify \
  -H "Content-Type: application/json" \
  -d '{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
    "secret": "your-256-bit-secret"
  }'
```

## 🔒 Security Features

- **No Logging**: Secrets and tokens are never logged
- **No Storage**: All processing is stateless and in-memory
- **Client-Side Processing**: Frontend handles display logic
- **Secure Defaults**: Uses industry-standard JWT library

⚠️ **Warning**: This tool is for development and debugging purposes. Never use production secrets or expose sensitive tokens in untrusted environments.

## 🛠️ Building for Production

### Build Binary

**Build a binary:**
```bash
go build -o jwt-viewer.exe
```

**Run the binary:**
```bash
./jwt-viewer.exe
```

**Build for different platforms:**
```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o jwt-viewer-windows.exe

# Linux
GOOS=linux GOARCH=amd64 go build -o jwt-viewer-linux

# macOS
GOOS=darwin GOARCH=amd64 go build -o jwt-viewer-macos
```

### Docker Deployment

The application includes a multi-stage Dockerfile for optimized production builds:

**Features:**
- Multi-stage build (reduces image size)
- Alpine Linux base (minimal footprint ~15MB)
- Non-root user execution
- Health checks included
- Production-ready configuration

**Docker Compose Configuration:**

The included `docker-compose.yml` provides:
- Automatic restart policy
- Health checks
- Network isolation
- Environment variable support
- Easy scaling and management

**Production Deployment:**

1. **Using Docker Compose:**
   ```bash
   # Start in production mode
   GIN_MODE=release docker-compose up -d
   
   # Scale if needed (future enhancement)
   docker-compose up -d --scale jwt-viewer=3
   
   # Update the application
   docker-compose pull
   docker-compose up -d
   ```

2. **Using Docker Swarm:**
   ```bash
   docker stack deploy -c docker-compose.yml jwt-stack
   ```

3. **Using Kubernetes:**
   ```bash
   # Create deployment from Docker image
   kubectl create deployment jwt-viewer --image=jwt-viewer:latest
   kubectl expose deployment jwt-viewer --type=LoadBalancer --port=8080
   ```

**Environment Variables:**

- `GIN_MODE`: Set to `release` for production (default: `debug`)
- Custom port can be configured by modifying `main.go`

## 📦 Dependencies

- **[gin-gonic/gin](https://github.com/gin-gonic/gin)** - High-performance HTTP web framework
- **[golang-jwt/jwt](https://github.com/golang-jwt/jwt)** - JWT implementation for Go

## 🎨 Technology Stack

- **Backend**: Go 1.20+
- **Web Framework**: Gin
- **JWT Library**: golang-jwt/jwt v5
- **Frontend**: Vanilla JavaScript, HTML5, CSS3
- **UI Design**: Modern gradient design with responsive layout

## 📝 Project Structure Details

### Handlers Layer
HTTP request/response handling, input validation, and error responses.

### Services Layer
Core business logic for JWT operations (decode, encode, verify).

### Models Layer
Data structures for requests and responses.

### Utils Layer
Helper functions for base64 decoding, timestamp parsing, and claim extraction.

## � Troubleshooting

### Common Issues

**Port Already in Use:**
```bash
# Check what's using port 8080
netstat -ano | findstr :8080  # Windows
lsof -i :8080                 # Linux/macOS

# Kill the process or change port in main.go
```

**Docker Issues:**
```bash
# View container logs
docker-compose logs -f

# Restart containers
docker-compose restart

# Rebuild after code changes
docker-compose up -d --build

# Clean up everything
docker-compose down -v
docker system prune -a
```

**Module Download Issues:**
```bash
# Clear Go module cache
go clean -modcache

# Re-download dependencies
go mod download
go mod tidy
```

**Static Files Not Loading:**
- Ensure the `static/` directory is in the same location as the binary
- Check file permissions
- Verify the paths in `main.go` are correct

### Docker Health Check

Check container health:
```bash
docker ps
docker inspect jwt-viewer | grep -A 10 Health
```

### Performance Tuning

For production deployments:
1. Set `GIN_MODE=release`
2. Use a reverse proxy (nginx, Traefik)
3. Enable HTTPS/TLS
4. Implement rate limiting
5. Monitor with Prometheus/Grafana

## 🚦 Future Enhancements

Potential additions (marked as optional in the spec):
- [ ] RSA algorithm support (RS256, RS384, RS512)
- [ ] Dark mode toggle
- [ ] OpenAPI/Swagger documentation
- [ ] Unit tests and integration tests
- [ ] Rate limiting middleware
- [x] Docker support ✅
- [ ] Structured logging (zerolog/zap)
- [ ] CI/CD pipeline (GitHub Actions)
- [ ] Helm charts for Kubernetes

## 📄 License

This is a demonstration project for educational purposes.

## 👨‍💻 Author

Built as a high-quality implementation of a JWT debugging tool similar to jwt.io.

---

## 🆘 Support

**Need help?** The application includes:
- Helpful tooltips and UI guidance
- Clear error messages
- Security warnings in the footer

Enjoy debugging your JWTs! 🚀
