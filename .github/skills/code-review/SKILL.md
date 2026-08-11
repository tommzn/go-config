# Code Review Skill for go-config

## Repository purpose
`go-config` is a Go library providing a unified interface to load YAML configuration from multiple sources (local file, in-memory string, AWS S3). It wraps [Viper](https://github.com/spf13/viper) and exposes a minimal `Config` interface with typed accessors.

## Key conventions

- **Interface-first**: All public functionality is defined via interfaces in `interfaces.go`. Implementations must satisfy these interfaces; do not add methods that bypass them.
- **Pointer-based return values**: Accessors return `*T` with a default value fallback, not `(T, error)`. Maintain this pattern for consistency.
- **100% test coverage**: Every code path must be covered by tests. Use `testify` suites (`suite.Suite`). Do not merge without full coverage.
- **No global state**: Config sources are instantiated explicitly. Avoid package-level variables.
- **Dot-notation keys**: Nested config access uses dot-notation (e.g. `"namespace.key"`). Tests should cover both flat and nested key access.

## Dependencies

- Use `github.com/spf13/viper` for YAML parsing — do not introduce alternative parsers.
- AWS S3 access uses `aws-sdk-go-v2` — stick to v2; do not mix with v1.
- Test assertions use `github.com/stretchr/testify`.
- Dependency updates go on `deps/**` branches.

## CI / branch policy

- All PRs target `main`.
- CI runs on `features/**`, `dependabot/**`, `copilot/**`, and `deps/**` branches.
- Branch protection enforces PRs for all changes including repo owner.
