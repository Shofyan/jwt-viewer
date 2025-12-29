# 🐳 Docker Deployment Guide

Complete guide for deploying JWT Debugger using Docker and Docker Compose.

## 📋 Table of Contents

- [Quick Start](#quick-start)
- [Docker Compose](#docker-compose)
- [Plain Docker](#plain-docker)
- [Configuration](#configuration)
- [Production Deployment](#production-deployment)
- [Troubleshooting](#troubleshooting)

## 🚀 Quick Start

The fastest way to get started:

```bash
# Clone or navigate to the project
cd jwt-viewer

# Start with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f

# Access the application
open http://localhost:8080
```

## 📦 Docker Compose

### Basic Commands

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs -f jwt-viewer

# Restart services
docker-compose restart

# Rebuild after changes
docker-compose up -d --build

# View status
docker-compose ps
```

### Configuration Options

Edit `docker-compose.yml` to customize:

```yaml
services:
  jwt-viewer:
    environment:
      - GIN_MODE=release        # Set to 'release' for production
    ports:
      - "8080:8080"             # Change host port if needed
    restart: unless-stopped     # Restart policy
```

### Multiple Instances

Run multiple instances with different ports:

```bash
# Modify docker-compose.yml or use override
docker-compose -f docker-compose.yml -f docker-compose.override.yml up -d
```

## 🐋 Plain Docker

### Build Image

```bash
# Build with tag
docker build -t jwt-viewer:latest .

# Build with specific version
docker build -t jwt-viewer:1.0.0 .

# Build with custom Dockerfile
docker build -f Dockerfile.prod -t jwt-viewer:prod .
```

### Run Container

```bash
# Basic run
docker run -d -p 8080:8080 --name jwt-viewer jwt-viewer:latest

# Run with environment variables
docker run -d -p 8080:8080 \
  -e GIN_MODE=release \
  --name jwt-viewer \
  jwt-viewer:latest

# Run with custom network
docker network create jwt-network
docker run -d -p 8080:8080 \
  --network jwt-network \
  --name jwt-viewer \
  jwt-viewer:latest

# Run with volume mount (for development)
docker run -d -p 8080:8080 \
  -v $(pwd)/static:/root/static \
  --name jwt-viewer \
  jwt-viewer:latest
```

### Container Management

```bash
# View logs
docker logs -f jwt-viewer

# Execute command inside container
docker exec -it jwt-viewer sh

# Inspect container
docker inspect jwt-viewer

# Check health status
docker inspect --format='{{json .State.Health}}' jwt-viewer

# Stop container
docker stop jwt-viewer

# Remove container
docker rm jwt-viewer

# Remove image
docker rmi jwt-viewer:latest
```

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default | Values |
|----------|-------------|---------|--------|
| `GIN_MODE` | Application mode | `debug` | `debug`, `release` |

### Port Configuration

To change the default port, modify `main.go`:

```go
port := ":9000"  // Change from :8080 to :9000
```

Then rebuild the image:

```bash
docker-compose up -d --build
```

### Health Checks

The Docker Compose configuration includes health checks:

```yaml
healthcheck:
  test: ["CMD", "wget", "--quiet", "--tries=1", "--spider", "http://localhost:8080/"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```

View health status:

```bash
docker-compose ps
docker inspect jwt-viewer | grep -A 10 Health
```

## 🚀 Production Deployment

### Best Practices

1. **Use Release Mode:**
   ```bash
   GIN_MODE=release docker-compose up -d
   ```

2. **Use Specific Version Tags:**
   ```bash
   docker build -t jwt-viewer:1.0.0 .
   ```

3. **Resource Limits:**
   ```yaml
   services:
     jwt-viewer:
       deploy:
         resources:
           limits:
             cpus: '0.5'
             memory: 256M
           reservations:
             cpus: '0.25'
             memory: 128M
   ```

4. **Logging Configuration:**
   ```yaml
   services:
     jwt-viewer:
       logging:
         driver: "json-file"
         options:
           max-size: "10m"
           max-file: "3"
   ```

### Docker Swarm

Deploy to Docker Swarm:

```bash
# Initialize swarm (if not already done)
docker swarm init

# Deploy stack
docker stack deploy -c docker-compose.yml jwt-stack

# List services
docker stack services jwt-stack

# View logs
docker service logs -f jwt-stack_jwt-viewer

# Scale service
docker service scale jwt-stack_jwt-viewer=3

# Remove stack
docker stack rm jwt-stack
```

### Kubernetes

Create Kubernetes manifests:

**deployment.yaml:**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: jwt-viewer
spec:
  replicas: 3
  selector:
    matchLabels:
      app: jwt-viewer
  template:
    metadata:
      labels:
        app: jwt-viewer
    spec:
      containers:
      - name: jwt-viewer
        image: jwt-viewer:latest
        ports:
        - containerPort: 8080
        env:
        - name: GIN_MODE
          value: "release"
        resources:
          limits:
            memory: "256Mi"
            cpu: "500m"
          requests:
            memory: "128Mi"
            cpu: "250m"
```

**service.yaml:**
```yaml
apiVersion: v1
kind: Service
metadata:
  name: jwt-viewer-service
spec:
  type: LoadBalancer
  selector:
    app: jwt-viewer
  ports:
  - port: 80
    targetPort: 8080
```

Deploy to Kubernetes:

```bash
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl get pods
kubectl get services
```

### Reverse Proxy (Nginx)

Example nginx configuration:

```nginx
server {
    listen 80;
    server_name jwt.example.com;

    location / {
        proxy_pass http://jwt-viewer:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
```

Add to docker-compose.yml:

```yaml
services:
  jwt-viewer:
    # ... existing config

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
    volumes:
      - ./nginx.conf:/etc/nginx/conf.d/default.conf
    depends_on:
      - jwt-viewer
    networks:
      - jwt-network
```

## 🔧 Troubleshooting

### Common Issues

**1. Port Already in Use:**
```bash
# Check what's using the port
docker ps
netstat -tulpn | grep 8080

# Change port in docker-compose.yml
ports:
  - "8081:8080"  # Use different host port
```

**2. Container Won't Start:**
```bash
# Check logs
docker-compose logs jwt-viewer

# Check container status
docker-compose ps

# Inspect container
docker inspect jwt-viewer
```

**3. Build Failures:**
```bash
# Clean build cache
docker builder prune

# Rebuild without cache
docker-compose build --no-cache

# Check Dockerfile syntax
docker build --dry-run -t jwt-viewer:test .
```

**4. Static Files Not Loading:**
```bash
# Verify files are copied
docker exec jwt-viewer ls -la /root/static

# Check permissions
docker exec jwt-viewer ls -l /root/
```

**5. Health Check Failing:**
```bash
# Test health check manually
docker exec jwt-viewer wget -q -O- http://localhost:8080/

# Check health status
docker inspect --format='{{json .State.Health}}' jwt-viewer | jq
```

### Performance Issues

**Check resource usage:**
```bash
# Container stats
docker stats jwt-viewer

# Detailed info
docker inspect jwt-viewer | grep -A 20 HostConfig
```

**Optimize image size:**
```bash
# Check image size
docker images jwt-viewer

# Remove unused layers
docker image prune

# Use multi-stage builds (already implemented)
```

### Debugging

**Access container shell:**
```bash
docker exec -it jwt-viewer sh
```

**View real-time logs:**
```bash
docker-compose logs -f --tail=100 jwt-viewer
```

**Network debugging:**
```bash
# Check network
docker network inspect jwt-network

# Test connectivity
docker exec jwt-viewer wget -q -O- http://localhost:8080/
```

## 🔄 Updates and Maintenance

### Update Application

```bash
# Pull latest changes (if from registry)
docker-compose pull

# Rebuild and restart
docker-compose up -d --build

# Or for plain Docker
docker build -t jwt-viewer:latest .
docker stop jwt-viewer
docker rm jwt-viewer
docker run -d -p 8080:8080 --name jwt-viewer jwt-viewer:latest
```

### Cleanup

```bash
# Remove stopped containers
docker-compose down

# Remove with volumes
docker-compose down -v

# Clean up unused resources
docker system prune -a

# Remove specific images
docker rmi jwt-viewer:latest
```

## 📊 Monitoring

### Logging

View logs in different formats:

```bash
# Follow logs
docker-compose logs -f

# Last 100 lines
docker-compose logs --tail=100

# Specific service
docker-compose logs jwt-viewer

# Since specific time
docker-compose logs --since 2024-01-01T00:00:00
```

### Metrics

For production monitoring, consider:
- Prometheus + Grafana
- ELK Stack (Elasticsearch, Logstash, Kibana)
- Datadog
- New Relic

## 🔐 Security

### Best Practices

1. **Run as non-root user** (TODO: add to Dockerfile)
2. **Scan images for vulnerabilities:**
   ```bash
   docker scan jwt-viewer:latest
   ```
3. **Use secrets management:**
   ```bash
   docker secret create jwt_secret secret.txt
   ```
4. **Limit container capabilities:**
   ```yaml
   cap_drop:
     - ALL
   cap_add:
     - NET_BIND_SERVICE
   ```

## 📚 Additional Resources

- [Docker Documentation](https://docs.docker.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [Dockerfile Best Practices](https://docs.docker.com/develop/develop-images/dockerfile_best-practices/)
- [Kubernetes Documentation](https://kubernetes.io/docs/)

---

For more information, see the main [README.md](README.md)
