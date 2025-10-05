# Togglr

Togglr is a feature flag and experimentation platform for developers.
It allows developers and product teams to toggle features on/off, run A/B tests, roll out features to specific user segments, and track feature stability.

## Features

* Feature flags organized by projects and environments (prod, stage, dev).
* Variants and targeting rules.
* Guarded features (pending changes with approval workflow).
* Categories and tags for organizing features.
* Schedules for automatic feature activation/deactivation.
* Full audit log of changes.
* SLA and health monitoring for features.
* Auto-disable on runtime errors (via error reports).
* Role-based access control (RBAC) with roles and permissions.
* REST API and WebSocket events.

## Architecture

* **Backend** — Go (PostgreSQL/TimescaleDB, NATS, REST API, WebSocket broadcaster).
* **Frontend** — React + TypeScript.
* **SDKs** available for:

    * Go
    * Ruby
    * PHP
    * Python
    * TypeScript (Node.js and browser)
    * Elixir

## Setup Development Environment

### Requirements

- Docker
- Docker Compose

### Quick Start

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd togglr
   ```

2. Setup environment files:
   ```bash
   make setup
   ```

3. Add hostname togglr.local to /etc/hosts

4. Start the development environment:
   ```bash
   make dev-up
   ```

5. Access the application:
   - Frontend: https://togglr.local
   - API: https://togglr.local/api/v1/
   - SDK: https://togglr.local/sdk/v1/

### Configuration Notes

- **Domain**: By default, the application is configured for `togglr.local` in `dev/config.env` and `dev/platform.env`
- **SSL Certificates**: Self-signed certificates are required in `dev/nginx/ssl/` directory. Pre-generated certificates are included but may be expired
- **Superuser**: On first startup, a superuser is created with:
  - Email: `ADMIN_EMAIL` from `dev/config.env` (default: `admin@togglr.dev`)
  - Password: `ADMIN_TMP_PASSWORD` from `dev/config.env` (default: `password543210`)
  - You can change these credentials after first login

### Development Commands

- `make dev-up` - Start all services
- `make dev-down` - Stop all services
- `make dev-clean` - Stop services and clean up volumes/images
- `make build` - Build the application (requires Go 1.25+)
- `make test` - Run tests (requires Go 1.25+)

## Usage

The server exposes REST API under `/api/v1/*`.

Prometheus metrics are available at `/metrics`.

WebSocket events are available at `/api/ws`.

SDK interface is available under `/sdk/v1/*`.

## Contributing

We welcome community contributions to Togglr!

Before submitting a pull request, please note:

- All contributions are subject to the [Togglr Business License (TBL)](./LICENSE).
- By contributing, you agree to the terms of our [Contributor License Agreement (CLA)](./CLA.md).
- This means:
    - You confirm that you have the legal right to submit the code.
    - You agree that your contribution will be licensed under TBL.
    - The project owner may also include your contribution in commercial licenses without any obligation to provide royalties, equity, or other compensation.
    - You retain authorship and copyright of your contribution, visible in Git history.

### How to contribute

1. Fork the repository.
2. Create a feature branch.
3. Make your changes.
4. Submit a pull request.

Please make sure to include tests and documentation updates where appropriate.
