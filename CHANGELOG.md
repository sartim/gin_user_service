# Changelog

All notable changes to this project are documented here. Entries follow
[Semantic Versioning](https://semver.org/).

## [Unreleased]

Changes planned for the next release are recorded here before implementation.

## [0.2.0] - 2026-08-31

- Added versioned PostgreSQL migrations with checksums and controlled rollback.
- Added PostgreSQL migration integration coverage.
- Added HTTP integration coverage for authentication and administrator access.
- Fixed CI migration validation and semantic-version release triggers.
- Added deployment and migration assets to release archives.

## [0.1.0] - 2026-08-30

Initial tracked production-hardening release candidate.

- Hardened JWT validation and administrator authorization.
- Replaced reflection-based CRUD with typed generic controllers.
- Added validated configuration and graceful HTTP shutdown.
- Added liveness/readiness endpoints and database pool settings.
- Added race-enabled tests and CI quality gates.
- Added non-root container, PostgreSQL Compose environment, and deployment service files.
- Removed tracked production environment secrets.

## Release links

[Unreleased]: https://github.com/sartim/gin_user_service/compare/v0.2.0...HEAD
[0.2.0]: https://github.com/sartim/gin_user_service/releases/tag/v0.2.0
[0.1.0]: https://github.com/sartim/gin_user_service/releases/tag/v0.1.0
