# Repository Guidelines

## Project Structure & Module Organization
- `cmd/viesquery/`: CLI entrypoint and main package.
- `internal/vies/`: VIES client, types, and validation logic.
- `internal/output/`: Plain and JSON formatters, date rendering.
- `docs/`: API spec, implementation notes, and WSDL reference.
- `bin/`: Built binaries (created by `make build`).
- `testdata/`: Test fixtures and sample payloads.

## Build, Test, and Development Commands
- `make build`: Compile binary to `bin/viesquery` (injects version from git).
- `make test`: Run unit tests with `-race` and coverage summary.
- `make test-coverage`: Produce `coverage.html` report.
- `make fmt`: Format Go code (`go fmt` + `gofmt -s`).
- `make lint`: Run `golangci-lint` (install via `make dev-setup`).
- `make check`: Format, lint, then test.
- `make run`: Build and show CLI help.
- `make release`: Cross-compile release binaries to `bin/`.

## Coding Style & Naming Conventions
- Use standard Go style; always run `make fmt` and `make lint`.
- Package names are short, lowercase, no underscores (e.g., `vies`, `output`).
- Exported identifiers use `CamelCase`; errors follow `ErrXxx` where useful.
- Tests live next to code in `*_test.go`; keep files cohesive and small.
- JSON field names use lowerCamelCase struct tags.

## Testing Guidelines
- Framework: Go’s built-in `testing` (`go test ./...`).
- Name tests `TestXxx`, benchmarks `BenchmarkXxx`.
- Place fixtures in `testdata/`; keep tests deterministic and network-stable.
- Aim to cover parsing, validation, and formatter behavior; include race checks.

## Commit & Pull Request Guidelines
- Use Conventional Commits: `feat:`, `fix:`, `docs:`, `refactor:`, etc.
- PRs should include: clear description, related issues (e.g., `Closes #123`), tests, and docs updates (`README.md`/`docs/`).
- Verify locally: `make check` and, if relevant, `make test-coverage`.
- Follow the PR template in `.github/pull_request_template.md`.

## Security & Configuration Tips
- Do not commit secrets; validate inputs and avoid leaking error details.
- Config defaults to `$XDG_CONFIG_HOME/viesquery/config.json` (see README).
- See `SECURITY.md` for reporting and hardening guidance.
