# 🏦 BRI EDC Backend

This project is a backend system designed to simulate the core functionalities of an **Electronic Data Capture (EDC)** service
It provides secure RESTful APIs for processing sale transactions and handling end-of-day settlements.

The system is designed using a **microservice architecture**, consisting of:
- **edc-service** — the public-facing REST API service.
- **core-service** — a simulated core banking authorization server.

The entire stack is containerized using **Docker** for consistent deployment, scalability, and easy local setup.

---

## 🚀 Features

- **Sale Transaction Processing**  
  Securely process individual sale transactions via REST API.

- **Batch Settlement**  
  Handle end-of-day settlements that aggregate multiple transactions into one batch.

- **Device Authentication (HMAC-SHA256)**  
  Protect sensitive endpoints (`/sale`, `/settlement`) using per-terminal HMAC signatures.

- **Secure Internal Communication (gRPC)**  
  The `edc-service` communicates with the internal `core-service` using **gRPC** secured by a pre-shared secret token.

- **API Documentation (Swagger)**  
  Interactive API documentation provided via Swagger UI for testing and inspection.

---

## 🧱 Tech Stack

| Component | Technology |
|------------|-------------|
| **Language** | Go |
| **Framework** | Fiber (for REST API) |
| **Database** | PostgreSQL |
| **Inter-service Communication** | gRPC |
| **ORM** | GORM |
| **Containerization** | Docker & Docker Compose |

---

## ⚙️ Setup and Deployment

### **Prerequisites**
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

---

### **1. Clone the Repository**
```bash
git clone <your-repository-url>
cd <your-project-directory>
```

### **2. Build and Run with Docker Compose**
From the root directory, run the following command:
```bash
docker-compose up -d --build
```

This command will build the Docker images for both services, start the containers, and run the application in the background.

---

## Accessing the Services

- **EDC REST API**: The main service will be available on port `8100`.  
  URL: `https://revand.test.bri-edc.arnatechnology.com`

- **Swagger Documentation**:  
  URL: `https://revand.test.bri-edc.arnatechnology.com/swagger`

---

## API Endpoints

- `POST /api/v1/auth` — Get bearer token for authentication.
- `POST /api/v1/transactions/sale` — Processes a single sale transaction.  
- `POST /api/v1/transactions/settlement` — Processes a batch of transactions for settlement.  