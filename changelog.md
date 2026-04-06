# Changelog: Vertigo (vortex-go)

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.0] - 2026-04-06

### Rebranded
- Project renamed from `persistence_handler_go` to **Vertigo**.
- Global module rename to `vertigo`.

### Added
- [FEATURE] **REST API Demo**: Built-in HTTP endpoints for user management and SQL dispatching.
- [FEATURE] **GraphQL Interface**: Modern data access layer with GraphiQL playground.
- [TOOL] **Postman Collection**: Pre-configured `vertigo_demo.postman_collection.json` for rapid API testing.
- [CONFIG] **YAML Configuration System**: Externalized database and network settings with `config.yaml`.

### Changed
- **Testing Structure**: Refactored BDD tests into the `features/` directory for better isolation.
- **Enhanced Documentation**: Redesigned README and Architecture docs for the Vertigo transition.

## [0.2.0] - 2026-04-06

### Added
- [FEATURE] **Raw SQL Query Dispatcher**: Decoupled from predefined query abstractions. Highly flexible SQL execution from any service layer.
- [FEATURE] **Zero-Copy Streaming**: Integrated `json-iterator` streaming for 10K+ row datasets to prevent memory spikes.
- [STABILITY] **Atomic Payloads**: Wrapped DB results with SQL metadata and weights in a unified DTO.

### Changed
- **Project Restructuring**: Migrated core logic to `pkg/` directory (`pkg/broker`, `pkg/db`, `pkg/model`) for better Go module compatibility.
- **Enhanced Readme**: Comprehensive service integration guide and architecture visualization.

## [0.1.0] - 2026-04-06


### Added
- [FEATURE] **Master Facade Pattern**: `DoubleBaseBroker` for simplified subsystem access.
- [FEATURE] **Double Base Architecture**: Disjoint DB (Base 1) and Real-time (Base 2) layers.
- [FEATURE] **Streaming JSON Engine**: Memory-efficient row-by-row encoding using `json-iterator`.
- [FEATURE] **BDD Testing**: Full Gherkin/Cucumber/Godog scenario verification.
- [FEATURE] **Declarative Entities**: Struct-tag-driven data models.
- [FEATURE] **Centrifugo Integration**: Publisher with auto-reconnect logic and connection state awareness.

### Fixed
- Build errors with `json-iterator` configuration and `centrifuge-go` configuration.
- Nil-pointer panic in `Dispatch` when Centrifugo is offline (improved resilience).
- Build error in Godog test runner in `main_test.go`.

### Changed
- Standardized port to `8010` for Centrifugo.
- Updated `DoubleBaseBroker` to use `Publisher` interface for mockability.
