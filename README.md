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
