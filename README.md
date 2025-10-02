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
