# QuickLink — High-Throughput URL Shortener

QuickLink is a horizontally scalable, production-ready URL Shortening service engineered with Go, Nginx, MariaDB, and Redis.

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Nginx](https://img.shields.io/badge/Nginx-009639?style=for-the-badge&logo=nginx&logoColor=white)
![MariaDB](https://img.shields.io/badge/MariaDB-003545?style=for-the-badge&logo=mariadb&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-DC382D?style=for-the-badge&logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![JavaScript](https://img.shields.io/badge/JavaScript-F7DF1E?style=for-the-badge&logo=javascript&logoColor=black)

---

## Key Features

- **Horizontal Scaling**: Nginx reverse proxy load-balances traffic across multiple Go application container replicas.
- **Sub-Millisecond Redirects**: In-memory Redis caching offloads read traffic from database queries.
- **Reliable Persistence**: MariaDB relational store with $O(1)$ primary key indexing (`VARCHAR(255) PRIMARY KEY`).
- **Collision Resilience**: Cryptographically secure 6-character Base64 generation with a 5-attempt retry loop.
- **Built-in Rate Limiting**: Thread-safe in-memory middleware protects endpoints from abuse.
- **Ultra-Slim Container Footprint**: Multi-stage Docker build producing a ~12MB static runtime image.

---

## Tech Stack Overview

| Component | Technology | Purpose |
| :--- | :--- | :--- |
| **Backend API** | **Go (Golang)** | High-concurrency HTTP handlers & business logic |
| **Load Balancer** | **Nginx** | Round-robin reverse proxy load distribution |
| **Primary Database** | **MariaDB** | Persistent storage for short-to-long URL mappings |
| **In-Memory Cache** | **Redis** | Ultra-fast caching layer for redirect lookups |
| **Containerization** | **Docker / Compose** | Infrastructure orchestration & scaling |

---

## Environment Variables (`.env`)

Create a `.env` file in the project root directory with the following variables:

```env
# Redis Configuration
REDIS_ADDR="localhost:6379"
REDIS_PASSWORD=""
REDIS_PORT=6379

# MariaDB Configuration
MARIADB_ROOT_PASSWORD=rootpassword
MARIADB_DATABASE=urlshortener
MARIADB_USER=appuser
MARIADB_PASSWORD=apppassword
MARIADB_PORT=3306
```

---

## Quick Start Guide

### Prerequisites
- [Docker](https://www.docker.com/) & Docker Compose installed.
- Python 3 (for running automated verification tests).

### 1. Launch Stack with 3 Server Replicas
Run Docker Compose with `--scale app=3` to launch the multi-node load balanced stack:

```bash
docker-compose up --build --scale app=3 -d
```

### 2. Access the Application
Open your browser or HTTP client:
- **Web Interface**: `http://localhost:8080/app`
- **API Endpoint**: `POST http://localhost:8080/shorten`

### 3. Run Load Balancer Verification Test
Execute the automated test script to verify routing across all scaled app instances:

```bash
python3 tests/test_load_balancer.py
```

---

## Documentation

Detailed technical design and command references are available under the `docs/` directory:

- [System Architecture Specification](docs/architecture.md) — Capacity planning, collision analysis, and component specifications.
- [Operational Command Reference](docs/commands.md) — Complete CLI cheatsheet for Docker, logs, and database management.

---

## Service Operations Summary

```bash
# Check status of running containers
docker-compose ps

# Stream logs from app instances
docker-compose logs -f app

# Stop the stack
docker-compose stop

# Clean up containers and volumes
docker-compose down -v
```
