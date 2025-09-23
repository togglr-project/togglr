# Togglr

Togglr is a feature flag and experimentation platform for modern applications.
It allows developers and product teams to toggle features on/off, run A/B tests, and roll out changes progressively without redeploying code.

✨ Features

 - Feature flags management
   - Create, update, and remove flags dynamically.
   - Organize flags by projects and environments (dev, staging, prod).
   - Versioning and audit log for all changes.
 - Progressive rollout
   - Gradual traffic shifting (e.g., 10% → 50% → 100%).
   - Percentage-based and segment-based rollouts.
 - A/B testing
   - Multiple flag variants with weighted distribution.
   - User segmentation (by ID, region, custom attributes).
   - Experiment tracking and integration with metrics.
 - SDKs
   - Go SDK (first-class support).
   - Planned: JavaScript/TypeScript, Python, Rust.
   - Client-side caching and periodic refresh.
   - Context-aware evaluation (userID, attributes, etc.).
 - Service architecture
   - Central API service for flag storage and evaluation.
   - PostgreSQL as primary storage (with migrations).
   - CLI tool (ffctl) for DevOps and automation.
   - REST/gRPC APIs for integration.
 - Admin Dashboard
   - Web-based UI for managing flags, rollouts, and experiments.
   - Real-time flag state overview.
   - Rule editor with JSON + form-based configuration.
 - Observability & Ops
   - Prometheus metrics for flag usage and evaluations.
   - Export logs for audit and compliance.
   - Docker & Helm charts for deployment.
 - Future roadmap
   - Multi-tenant mode.
   - RBAC and SSO (SAML/OAuth2).
   - Webhooks & integrations with CI/CD pipelines.
   - Support for alternative backends (etcd, Redis, Consul).
   - Hot reload & real-time streaming of flag updates.
 