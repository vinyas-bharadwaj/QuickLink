# URLShortener Command Reference

This document provides a comprehensive operational reference for managing, scaling, testing, and monitoring the URLShortener application stack.

---

## Docker Compose Management

### Start Services
Build and start the application, Nginx load balancer, MariaDB, and Redis containers in detached mode:

```bash
docker-compose up --build -d
```

### Scale Application Replicas
Horizontally scale the application service to 3 instances behind the Nginx load balancer:

```bash
docker-compose up --build --scale app=3 -d
```

### Stop Services
Gracefully stop all running containers:

```bash
docker-compose stop
```

### Shut Down and Clean Up
Stop and remove all containers, networks, and persistent volume data:

```bash
docker-compose down -v
```

---

## Verification & Testing

### Verify Load Balancer Functionality
Execute the Python verification suite to test endpoint routing, URL creation, and HTTP 302 redirects across scaled instances:

```bash
python3 tests/test_load_balancer.py
```

### Custom Test Execution
Run the verification test with custom request counts and target URL:

```bash
python3 tests/test_load_balancer.py --url http://localhost:8080 -n 10
```

---

## Monitoring & Diagnostic Logging

### Check Running Containers
Display active containers and mapped ports:

```bash
docker-compose ps
```

### Stream Application Logs
Stream live log outputs from all scaled application container replicas (`app-1`, `app-2`, `app-3`):

```bash
docker-compose logs -f app
```

### Stream Load Balancer Logs
Stream access and routing logs from the Nginx reverse proxy (`lb`):

```bash
docker-compose logs -f lb
```

### View Recent Container Logs
View the last 50 log lines for a specific container service:

```bash
docker-compose logs --tail 50 app
```

---

## Database & Cache CLI Operations

### Connect to Redis CLI
Interact with the Redis caching instance using `redis-cli`:

```bash
docker-compose exec redis redis-cli
```

### Connect to MariaDB CLI
Interact with the MariaDB database using the MariaDB command-line client:

```bash
docker-compose exec mariadb mariadb -u appuser -p urlshortener
```

### Run Standalone Redis Instance
Run a standalone Redis container for isolated development:

```bash
docker run -d --name redis -p 6379:6379 redis:latest
```
