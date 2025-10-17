# Togglr

[![CI](https://github.com/togglr-project/togglr/actions/workflows/ci.yml/badge.svg)](https://github.com/togglr-project/togglr/actions/workflows/ci.yml)

Togglr is a **feature flag and experimentation platform** for developers and product teams.

It allows you to **toggle features on/off**, run **A/B and multi-variant experiments**, and gradually roll out features based on user segments, time, or adaptive learning algorithms.

---

## 🚀 Experimentation and Algorithms

Togglr extends classic feature flags into a **full experimentation framework**.

Each feature can have multiple **variants** (A/B/N arms), and you can attach an **algorithm** that dynamically controls which variant is shown to users — optimizing for reward or conversion automatically.

### Supported algorithms

Currently implemented:

- **ε-Greedy** — explores randomly with probability *ε*, otherwise exploits the best variant.
- **Thompson Sampling** — Bayesian approach balancing exploration and exploitation.
- **Upper Confidence Bound (UCB)** — selects the variant with the highest confidence bound for expected reward.

Each algorithm can be configured per environment (e.g., `production`, `staging`) with custom numeric settings (for example):

```json
{
  "epsilon": 0.1
}
````

### Architecture of Experiments

* Each **feature** has one active algorithm per **environment**.
* Each algorithm tracks statistics (evaluations, successes, failures, rewards) per variant.
* Feedback events (from SDKs or API) are aggregated in TimescaleDB.
* The **BanditManager** keeps live statistics in memory and periodically syncs them to TimescaleDB.
* Algorithms learn over time and adjust rollout dynamically — no manual tweaking required.

> “One algorithm — one brain per feature.
> Variants are its hands.
> Feedback is its experience.”

---

## Features

* Feature flags organized by projects and environments (prod, stage, dev)
* Multi-variant flags and rollout experiments
* Algorithms for adaptive rollout (ε-Greedy, Thompson Sampling, UCB)
* Guarded features (pending changes with approval workflow)
* Categories and tags for organizing features
* Segments and targeting rules for user context
* Schedules for automatic activation/deactivation
* Full audit log of changes
* Health monitoring for features
* Auto-disable on runtime errors (via error reports)
* Role-based access control (RBAC)
* REST API and WebSocket events
* SDKs for multiple languages
* 2FA, LDAP, SSO/SAML authentication

---

## Architecture

* **Backend** — Go (TimescaleDB, NATS JetStream, REST API, WebSocket broadcaster)
* **Frontend** — React + TypeScript
* **SDKs** — Go, Python, Ruby, TypeScript, Elixir

Experiment data and feedback (continuous aggregates for dashboards), algorithm configuration and state are stored in TimescaleDB.

---

## Setup Development Environment

### Requirements

* Docker
* Docker Compose

### Quick Start

```bash
git clone https://github.com/togglr-project/togglr.git
cd togglr
make setup
make dev-up
```

Visit:

* Frontend: [https://togglr.local](https://togglr.local)
* API: [https://togglr.local/api/v1/](https://togglr.local/api/v1/)
* SDK: [https://togglr.local/sdk/v1/](https://togglr.local/sdk/v1/)

### Configuration Notes

* Default domain: `togglr.local`
* SSL: self-signed certs under `dev/nginx/ssl/`
* Default superuser:

    * Email: `admin@togglr.dev`
    * Password: `password543210`

---

## Usage

* REST API under `/api/v1/*`
* WebSocket events under `/api/ws`
* SDK endpoints under `/sdk/v1/*`
* Prometheus metrics at `:8081/metrics`

---

## Contributing

We welcome community contributions to Togglr!

All contributions are subject to the [Togglr Business License (TBL)](./LICENSE)
and our [Contributor License Agreement (CLA)](./CLA.md).

### How to contribute

1. Fork the repository.
2. Create a feature branch.
3. Make your changes.
4. Submit a pull request.

Please make sure to include tests and documentation updates where appropriate.
