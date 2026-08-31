# Contributing and release process

Every implementation is tracked as a small, reviewable unit. Do not merge
unrelated changes together.

## Implementation workflow

1. Create a normal branch such as `feature/refresh-tokens`, `fix/auth-header`,
   or `chore/ci-cache`.
2. Update the progress section in the pull request description as work advances.
3. Add or update tests for the behavior being changed.
4. Add an entry under `Unreleased` in `CHANGELOG.md`.
5. Open a pull request and wait for CI and review.
6. Merge only after the checklist is complete.

Migration files use the format `NNN_description.sql`. Never edit a migration
after it has been applied; add a new migration instead. Applied migrations are
checksum-verified at startup. Each migration must include matching
`NNN_description.up.sql` and `NNN_description.down.sql` files. Rollbacks are
manual, one migration at a time, and must be reviewed before production use.

## Semantic versioning

Use `MAJOR.MINOR.PATCH` tags prefixed with `v`:

- `PATCH`: backward-compatible bug or security fix (`v0.1.1`).
- `MINOR`: backward-compatible feature (`v0.2.0`).
- `MAJOR`: breaking API, schema, or operational contract change (`v1.0.0`).

After a PR is merged into `master`, the release workflow analyzes Conventional
Commit messages with `go-semantic-release`, calculates the next semantic
version, creates the Git tag, builds release archives with GoReleaser, and
publishes the GitHub Release automatically. Do not create release tags
manually.

Use Conventional Commit prefixes so release impact can be calculated reliably:
`fix:` for patches, `feat:` for minor releases, and `BREAKING CHANGE:` or `!`
for major releases.

## Progress tracking

Each PR must state:

- Objective and scope
- Current status and completed work
- Tests and verification performed
- Database/configuration/deployment impact
- Rollback plan
- Follow-up work, if any
