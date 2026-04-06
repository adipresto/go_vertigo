# Architecture: Vertigo (Triple Base)

This document describes the architectural patterns used in the **Vertigo** project.

## 1. Master Facade Design Pattern
The `TripleBaseBroker` acts as the master facade. It is the only interface that the business logic (higher-level application code) interacts with. It abstracts away:
- Database connection lifecycle (SQLite WAL).
- Real-time messaging connectivity (Centrifugo).
- Industrial connectivity (MQTT).
- Query registration and SQL execution.
- Data Transformation (DTO) logic.

## 2. Triple Base Architecture
Vertigo decouples three primary data layers:
1.  **Base 1 (Persistence):** SQL Database (SQLite) with a connection pool. It serves as the durable source of truth.
2.  **Base 2 (Network):** Real-time messaging (Centrifugo). It serves as the low-latency delivery mechanism for active users.
3.  **Base 3 (Industrial Connectivity):** MQTT Broker. It serves as the standard protocol for industrial IoT and edge integration.

### Selective Activation & Resilience
If Base 2 or Base 3 is down, the system remains fully functional for persistence (Base 1). Users can selectively enable or disable Base 2 and 3 via `config.yaml`.
- **Enabled**: The broker attempts connection and logs warnings if offline.
- **Disabled**: The broker skips connection entirely, running in a clean "Base 1 only" mode.

## 3. Streaming & Zero-Copy Engine
To avoid **Out-of-Memory (OOM)** failures when handling 1,000,000+ data points:
- The system does **not** load full result sets into Go slices.
- It uses `json-iterator`'s streaming encoder to process database rows as they arrive from the driver.
- Data is transformed and written directly to the output buffer for immediate transmission to Base 2.
- **Complexity**: O(1) space relative to the number of rows.

## 4. API Gateway Subsystem
Vertigo provides a multi-interface gateway for high-level interaction:
- **GraphQL**: Type-safe, declarative data fetching using `graphql-go`.
- **REST**: Standard JSON endpoints for legacy integration.
- **Gateway Abstraction**: The `main.go` entry point wires these interfaces directly into the `VertigoBroker`.

## 5. Configuration Resilience
- **YAML Management**: Centralized `config.yaml` for environment-specific settings.
- **Separation of Concerns**: Database paths and network keys are managed outside the core logic.

## 6. Behavior Driven Development (BDD)
We use a **Gherkin-first** testing approach.
- Feature definitions in `./features/` serve as living documentation.
- `godog` maps these scenarios to Go unit tests in `features/main_test.go`.
- **Command**: `go test -v ./features/...`
