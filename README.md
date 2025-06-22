
# 📡 Device Monitoring PoC

A production-ready backend PoC to monitor network devices (routers, switches, cameras, doors) with REST & gRPC APIs. Built for trade show demonstration — scalable, resilient, and extensible.

---

## 🚀 Features

- ✅ REST API to register, monitor and fetch device status
- ✅ gRPC server scaffolded (handlers ready)
- ✅ PostgreSQL for persistent device and status storage
- ✅ Periodic device health checks with retries
- ✅ External checksum generator (via binary)
- ✅ Structured logging with internal logger
- ✅ Fully dockerized with `docker-compose`

---

## 🧩 Services

| Port | Protocol | Description         |
|------|----------|---------------------|
| 8080 | HTTP     | REST API            |
| 50051| gRPC     | gRPC API (optional) |
| —    | Binary   | Checksum generator  |

---

## 🧪 How to Run

```bash
# Build & start
make build
make up
```

OR manually:

```bash
docker-compose up --build
```

---

## 📂 Project Structure

```
.
├── cmd/
│   └── checksum-generator/
├── internal/
│   ├── grpc/
│   ├── handlers/
│   ├── logger/
│   ├── models/
│   ├── repositories/
│   ├── services/
│   └── proto/
├── Dockerfile(s)
├── docker-compose.yml
└── README.md
```

