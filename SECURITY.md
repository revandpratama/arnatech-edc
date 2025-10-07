# Security Model

This document outlines the security mechanisms implemented in the **BRI EDC Backend Service** to ensure data protection, device authentication, and secure inter-service communication.

---

## 1. Device Authentication: JWT

All sensitive endpoints (`/sale` and `/settlement`) are protected using **JWT-HMAC with SHA-256** signatures for **stateless authentication**.

### **Flow**
1. Each registered EDC terminal is assigned a **unique secret key** (`secret_key`) stored securely in the database.  
2. When the terminal boots up or needs to perform a transaction, it first requests an **authentication token** using its credentials (e.g., `terminal_id` and `shared_secret`).  
3. The server validates the credentials and generates a **JWT**, signed using the terminal’s individual secret (`jwt_secret`).  
4. The terminal attaches this JWT in the `Authorization` header for every subsequent request (e.g., `Authorization: Bearer <token>`).  
5. On each incoming request:
   - The middleware extracts the token from the header.
   - It retrieves the terminal’s corresponding secret from the database.
   - The token is verified using that secret to ensure authenticity and integrity.
6. If verification succeeds, the request proceeds to the protected endpoint; otherwise, it is rejected with an **Unauthorized** response.

### **Token Structure**
A standard JWT includes:
- **Header:** Algorithm (`HS256`) and token type.  
- **Payload:** Terminal-specific claims such as:
  ```json
  {
    "terminal_id": "T01",
    "iat": 1733640000,
    "exp": 1733643600
  }

### **Benefits**
- **Authenticity** — Ensures the request comes from a legitimate, registered terminal.  
- **Integrity** — Guarantees that the request body has not been tampered with during transmission.

---

## 2. Secure Inter-Service Communication

Communication between the public-facing `edc-service` and the internal `core-service` is secured to prevent unauthorized internal API calls.

### **Method**
- **gRPC with a Pre-Shared Secret Token**

### **Implementation**
1. A static secret token is shared between both services via **environment variables**.  
2. The `edc-service` (gRPC client) attaches this token in the metadata (headers) of each outgoing RPC call.  
3. The `core-service` (gRPC server) employs a **gRPC interceptor** to inspect every incoming request.  
4. If the token is missing or invalid, the interceptor immediately rejects the request with an **Unauthenticated** status.

---

## 3. Data Security

### **Card Number Masking**
- The full Primary Account Number (PAN) of a customer's card is **never stored** in the database.  
- The service temporarily receives the full card number for transaction processing, then masks it (e.g., `411111******1111`) before saving the record.  
- This reduces exposure risk in the event of a data breach.

### **Secret Management**
- Sensitive configuration data (e.g., database credentials, HMAC secrets, and internal tokens) are stored as **environment variables**.  
- No secrets are hardcoded in the source code, following the **principle of least privilege** and **secure configuration management** best practices.

---
