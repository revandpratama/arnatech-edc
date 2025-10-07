# Architecture Decision Record (ADR) - EDC Backend Service

This document records the key architectural decisions made during the development of the **EDC Backend Service** project.

---

## System Architecture

### Context
The project requires building a backend system to handle EDC transactions — including sales, settlements, and secure communication with a simulated core banking system.  
The system needs to be scalable, maintainable, and have a clear separation of concerns.

adopting a **microservice architecture** composed of two main services:

- **edc-service**: The primary, public-facing REST API responsible for handling requests from EDC terminals, managing business logic, and persisting data.  
- **core-service**: An internal gRPC service that simulates a core banking server, responsible for authorizing or declining transactions.

I managed them both in the same docker-compose

### Rationale
- **Scalability**: Each service can be scaled independently based on load.  
- **Requirement Alignment**: This architecture meets the test requirement to create a separate internal authorization service.

---

## Device Authentication Method

### Context
Each EDC terminal must be authenticated to ensure that only legitimate devices can submit transactions. The system design allows for either JWT (session-based) or HMAC (stateless) authentication.

### Decision
I chose to use **JWT-HS256 signatures** for terminal authentication.  
Each terminal has a pre-shared secret key stored in the database. For every request, the terminal signs the request body and includes the signature in the `X-Signature` header.

### Rationale
- **Statelessness**: HMAC is stateless and fits well for EDC-like devices that don't need session management.  
- **Integrity**: HMAC ensures both authenticity of the sender and integrity of the message payload.  
- **Simplicity**: Avoids the complexity of login endpoints and token refresh logic, keeping the system design clean and lightweight.

---

## Inter-Service Communication Protocol

### Context
The `edc-service` must communicate securely and efficiently with the internal `core-service` for transaction authorization. Both REST and gRPC were considered.

### Decision
I decided to use **gRPC** for communication between the `edc-service` and `core-service`.  
Security is enforced through a **pre-shared secret token** included in the gRPC metadata on every call.

### Rationale
- **Performance**: gRPC’s binary protocol (Protobuf) offers better performance and lower latency than JSON/HTTP.  
- **Type Safety**: Protocol Buffers enforce a strict API contract, reducing integration issues.  
- **Bonus Alignment**: This approach fulfills one of the optional bonus criteria outlined in the project’s technical requirements.

---

## Database Primary Key Strategy

### Context
A consistent method is needed to uniquely identify records and define relationships (e.g., linking Terminals to Merchants). Options include using business identifiers or surrogate primary keys.

### Decision
I decided to use **auto-incrementing integers (BIGSERIAL)** as surrogate primary keys (`id`) for all tables.  
Business identifiers like `merchant_id` and `terminal_id` are stored as unique columns but are not used as primary keys.

### Rationale
- **Stability**: Surrogate keys remain unchanged even if business identifiers are updated, preserving referential integrity.  
- **Performance**: Integer-based joins are faster and more efficient than string-based joins.  
- **Decoupling**: Keeps the internal database schema independent from external API identifiers, allowing flexibility for future updates.

---
