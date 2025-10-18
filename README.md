# Togglr

[![CI](https://github.com/togglr-project/togglr/actions/workflows/ci.yml/badge.svg)](https://github.com/togglr-project/togglr/actions/workflows/ci.yml)

Togglr is a **feature flag and experimentation platform** for developers and product teams.

It allows you to **toggle features on/off**, run **A/B and multi-variant experiments**, and gradually roll out features based on user segments, time, or adaptive learning algorithms.

## Table of Contents

- [🚀 Experimentation and Algorithms](#-experimentation-and-algorithms)
  - [Supported algorithms](#supported-algorithms)
  - [Architecture of Experiments](#architecture-of-experiments)
- [Features](#features)
- [Architecture](#architecture)
- [Installation](#installation)
  - [System Requirements](#system-requirements)
  - [Quick Installation](#quick-installation)
  - [Manual Installation](#manual-installation)
  - [Environment Variables](#environment-variables)
  - [Post-Installation](#post-installation)
- [Setup Development Environment](#setup-development-environment)
- [Usage](#usage)
- [Contributing](#contributing)
- [Troubleshooting](#troubleshooting)
  - [Installation Issues](#installation-issues)
  - [Platform Management](#platform-management)
  - [Getting Help](#getting-help)

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
* **SDKs** — Go, Python, Ruby, TypeScript, Elixir, Rust

Experiment data and feedback (continuous aggregates for dashboards), algorithm configuration and state are stored in TimescaleDB.

---

## Installation

### System Requirements

**Minimum requirements:**
- Linux or macOS
- Docker 20.10+
- Docker Compose 2.0+
- 2GB RAM
- 10GB free disk space
- Root/sudo access

**Recommended for production:**
- 4GB+ RAM
- 50GB+ free disk space
- SSL certificate (Let's Encrypt recommended)
- SMTP server for email notifications

### Quick Installation

The easiest way to install Togglr is using our automated installer script:

```bash
# Download and run the quick installer
curl -fsSL https://raw.githubusercontent.com/togglr-project/togglr/main/prod/install.sh | sudo bash
```

Or if you have the repository cloned locally:

```bash
cd togglr
sudo ./prod/quick_install.sh
```

The installer will:
- Check system requirements (Docker, Docker Compose, OpenSSL, Make)
- Guide you through the configuration process
- Create necessary directories and configuration files
- Generate SSL certificates (self-signed or use your existing ones)
- Set up the complete Togglr platform

### Manual Installation

If you prefer to install manually or need custom configuration:

1. **Clone the repository:**
   ```bash
   git clone https://github.com/togglr-project/togglr.git
   cd togglr
   ```

2. **Run the installer:**
   ```bash
   sudo ./prod/quick_install.sh
   ```

3. **Follow the interactive prompts** to configure:
   - Administrator email
   - Domain name
   - SSL certificate settings
   - SMTP server configuration

### Environment Variables

For non-interactive installation, you can set environment variables:

```bash
sudo TOGGLR_ADMIN_EMAIL=admin@example.com \
     TOGGLR_DOMAIN=example.com \
     TOGGLR_MAILER_ADDR=smtp.example.com:587 \
     TOGGLR_MAILER_USER=admin@example.com \
     TOGGLR_MAILER_PASSWORD=your_password \
     TOGGLR_MAILER_FROM=noreply@example.com \
     TOGGLR_HAS_SSL_CERT=false \
     ./prod/install.sh
```

**Available environment variables:**
- `TOGGLR_ADMIN_EMAIL` - Administrator email (required)
- `TOGGLR_DOMAIN` - Platform domain (required)
- `TOGGLR_MAILER_ADDR` - SMTP server address (required)
- `TOGGLR_MAILER_USER` - SMTP username (default: admin@domain)
- `TOGGLR_MAILER_PASSWORD` - SMTP password (default: password)
- `TOGGLR_MAILER_FROM` - Email sender address (default: noreply@domain)
- `TOGGLR_HAS_SSL_CERT` - Whether you have existing SSL certificate (default: false)

### Post-Installation

After installation, the platform will be available at:
- **Frontend**: `https://your-domain`
- **API**: `https://your-domain/api/v1/`
- **SDK**: `https://your-domain/sdk/v1/`

**Default login credentials:**
- Email: The administrator email you provided during installation
- Password: A temporary password generated during installation (you'll be prompted to change it on first login)

**Managing the platform:**
```bash
cd /opt/togglr
make up    # Start the platform
make down  # Stop the platform
make pull  # Update Docker images
```

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

---

## Troubleshooting

### Installation Issues

**Installer hangs or freezes:**
- Ensure you're running in an interactive terminal
- Check that you have sudo/root privileges
- Try running with environment variables for non-interactive mode

**Docker permission issues:**
```bash
sudo usermod -aG docker $USER
# Log out and log back in
```

**SSL certificate issues:**
- For production, use Let's Encrypt or a trusted CA certificate
- Place your certificate files at `/opt/togglr/nginx/ssl/nginx_cert.pem` and `/opt/togglr/nginx/ssl/nginx_key.pem`

**Database connection issues:**
- Check if PostgreSQL container is running: `docker ps`
- Verify database credentials in `/opt/togglr/config.env`
- Check logs: `docker logs togglr-postgresql`

### Platform Management

**Start/stop the platform:**
```bash
cd /opt/togglr
make up    # Start all services
make down  # Stop all services
```

**View logs:**
```bash
docker logs togglr-backend
docker logs togglr-frontend
docker logs togglr-postgresql
```

**Update the platform:**
```bash
cd /opt/togglr
make pull  # Pull latest images
make up    # Restart with new images
```

**Reset admin password:**
```bash
# Edit /opt/togglr/config.env and change ADMIN_TMP_PASSWORD
# Then restart: make down && make up
```

### Getting Help

- **Documentation**: Check this README and inline help
- **Issues**: Report bugs on [GitHub Issues](https://github.com/togglr-project/togglr/issues)
- **Discussions**: Join community discussions on [GitHub Discussions](https://github.com/togglr-project/togglr/discussions)
