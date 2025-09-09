# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

VIES Query is a production-ready command-line tool for validating European Union VAT identification numbers using the VIES (VAT Information Exchange System) API. The tool provides fast, reliable validation of EU VAT numbers with support for all 27 EU member states through client-side format validation and real-time verification against the official European Commission VIES service.

## Development Commands

The project uses a comprehensive Makefile for all development tasks:

```bash
# Core development workflow
make build          # Build the binary to bin/viesquery
make test           # Run all tests with race detection and coverage
make check          # Run format, lint, and test (full quality check)
make clean          # Remove build artifacts

# Code quality
make fmt            # Format Go code with gofmt
make lint           # Run golangci-lint
make test-coverage  # Generate HTML coverage report (coverage.html)

# Module management  
make mod-tidy       # Tidy go modules
make mod-verify     # Verify go modules

# Release and distribution
make release        # Build cross-platform binaries (Linux, macOS, Windows)
make install        # Install binary to GOPATH/bin

# Development setup
make dev-setup      # Install required development tools (golangci-lint)
```

### Running and Testing

```bash
# Run built binary with example
make run            # Shows help output

# Run specific tests
go test ./internal/vies -v                    # VIES client tests
go test ./internal/output -v                  # Output formatter tests
go test -run TestValidateFormat ./internal/vies  # Specific test

# Test with different verbosity
./bin/viesquery --verbose --timeout 10 DE123456789
```

## Architecture

The project follows a clean architecture pattern with clear separation of concerns:

### Package Structure

```
cmd/viesquery/          # CLI application entry point
├── main.go            # Command-line parsing, error handling, orchestration

internal/vies/          # Core VIES API client and validation
├── client.go          # SOAP client implementation
├── types.go           # Data structures and client options
├── validation.go      # VAT format validation for all EU countries
└── types_test.go      # SOAP marshaling and namespace tests

internal/output/        # Output formatting
├── formatter.go       # Formatter interface and manager
├── plain.go           # Plain text formatter
└── json.go            # JSON formatter

docs/                   # Documentation
└── [various .md files] # Requirements, implementation notes, API specs
```

### Key Components

1. **CLI Layer** (`cmd/viesquery/main.go`): Handles argument parsing, environment variables, and coordinates between validation, API calls, and output formatting. Contains comprehensive error handling with appropriate exit codes.

2. **VIES Client** (`internal/vies/client.go`): SOAP client with custom XML namespace handling to work around VIES service quirks. Features configurable timeouts, TLS security settings, and detailed response parsing.

3. **Validation Engine** (`internal/vies/validation.go`): Client-side VAT format validation for all 27 EU member states using regex patterns. Handles country-specific rules like Austria's "ATU" prefix and Greece's EL/GR code conversion.

4. **Output System** (`internal/output/`): Pluggable formatter system supporting plain text and JSON outputs with comprehensive error formatting.

## Critical Implementation Details

### SOAP Client Challenges

The VIES service requires specific XML namespace handling that caused compatibility issues. The solution uses:

```go
// Custom SOAP envelope with explicit namespace declarations
type SOAPEnvelope struct {
    XMLName      xml.Name `xml:"soapenv:Envelope"`
    XmlnsSoapenv string   `xml:"xmlns:soapenv,attr"`
    XmlnsUrn     string   `xml:"xmlns:urn,attr"`
    Body         SOAPBody `xml:"soapenv:Body"`
}

// Namespace constant matching VIES expectations
const soapNamespace = "urn:ec.europa.eu:taxud:vies:services:checkVat:types"
```

### VAT Number Processing

The validation system handles country-specific quirks:

- **Austria**: Removes "U" prefix from "ATU" format for API calls
- **Greece**: Converts "GR" codes to "EL" for VIES compatibility  
- **All countries**: Regex patterns validate format before expensive API calls

### Error Handling Pattern

The codebase uses typed errors for precise error handling:

```go
type ValidationError struct {    // Client-side format errors
    Code, Message, VATNumber string
}

type ServiceError struct {       // VIES API and network errors
    Code, Message, VATNumber string
}
```

Exit codes map to error types: 1 (args), 2 (network), 3 (format), 4 (service unavailable).

## Testing Strategy

The project uses table-driven tests and focuses on:

1. **SOAP XML Generation**: Tests in `types_test.go` validate correct namespace handling and XML structure
2. **VAT Format Validation**: Comprehensive coverage of all 27 EU country formats
3. **Error Scenarios**: Network timeouts, service faults, malformed responses
4. **Output Formatting**: Both JSON and plain text output validation

Key test commands:
```bash
go test ./... -v -race -cover    # Full test suite
go test -bench=. ./...           # Benchmark performance
```

## Technical Considerations

### VAT Number Support

Supports all 27 EU member states with country-specific validation patterns:
- Format pre-validation prevents unnecessary API calls
- Handles variable-length formats (e.g., Romania: 2-10 digits, Lithuania: 9 or 12 digits)
- Special character handling (Netherlands "B" separator, Cyprus letter suffix)

### Network Configuration

- **TLS Security**: Enforces TLS 1.2+ with proper certificate validation
- **Timeouts**: Configurable request timeouts (default 30s, configurable via flag/env)
- **HTTP Transport**: Optimized with connection pooling and reasonable defaults
- **User Agent**: Proper identification for VIES service logs

### Configuration Sources

The CLI accepts configuration from multiple sources (precedence order):
1. Command-line flags (`--format`, `--timeout`, `--verbose`)
2. Environment variables (`VIESQUERY_FORMAT`, `VIESQUERY_TIMEOUT`, `VIESQUERY_VERBOSE`)
3. Built-in defaults

### Performance Characteristics

- **Binary Size**: < 10MB (no external dependencies)
- **Memory Usage**: < 50MB runtime
- **Response Time**: < 1ms format validation, < 5s API calls
- **Concurrency**: Thread-safe client implementation

## Build and Distribution

Cross-platform support via Go's native compilation:
- Linux: amd64, arm64
- macOS: amd64 (Intel), arm64 (Apple Silicon)  
- Windows: amd64

Release artifacts are created with `make release` and include version information embedded at build time using `-ldflags`.

## Development Notes

- **Go Version**: Requires Go 1.21+ (currently targeting 1.25.1)
- **Dependencies**: Uses only standard library (no external dependencies)
- **Code Quality**: Enforced via golangci-lint with comprehensive rule set
- **Documentation**: All public APIs documented with godoc comments

The codebase prioritizes production readiness with comprehensive error handling, security best practices, and performance optimization suitable for automation and CI/CD integration.

## Backlog / TODO (Release & Branch Management)

The following items are intentionally deferred and should be picked up when formalizing release management:

1. Release management
   - Adopt Semantic Versioning; use prereleases while pre-production.
   - Create annotated tags per release and publish GitHub Releases with artifacts from `make release` (Linux amd64/arm64, macOS amd64/arm64, Windows amd64).
   - Maintain a CHANGELOG.md (Keep a Changelog format) and link releases to tags.
   - Document release steps in `docs/RELEASING.md` (build, sign, tag, release notes, attach binaries).

2. Branch protection for `main`
   - Require pull request reviews (≥1 approval); dismiss stale approvals on new commits.
   - Require status checks (fmt, lint, test) to pass before merging; consider “Require branches to be up to date”.
   - Enforce linear history (squash or rebase merges) and restrict direct pushes to `main`.
   - Optionally enforce signed commits/tags.

3. CI checks (non-blocking until enabled)
   - Provide CI jobs for `make fmt`, `make lint`, `make test` and wire them as required status checks.

4. Repository governance
   - Add `CODEOWNERS` for critical paths (e.g., `internal/vies/**`, `internal/output/**`, `cmd/viesquery/**`).
   - Define contribution guidelines in `CONTRIBUTING.md` (PR process, coding standards, DCO/sign-offs if adopted).

Note: Do not enable or run any cost-incurring pipelines without explicit human approval.
