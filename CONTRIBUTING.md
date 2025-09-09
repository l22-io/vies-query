# Contributing to VIES Query

Thank you for your interest in contributing to VIES Query! This document provides guidelines and information for contributors.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing Requirements](#testing-requirements)
- [Submitting Changes](#submitting-changes)
- [Release Process](#release-process)
- [Getting Help](#getting-help)

## Code of Conduct

This project adheres to a Code of Conduct that we expect all participants to follow. Please read [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md) before contributing.

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- Make
- golangci-lint (installed via `make dev-setup`)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone git@github.com:YOUR_USERNAME/vies-query.git
   cd vies-query
   ```
3. Add the upstream remote:
   ```bash
   git remote add upstream git@github.com:l22-io/vies-query.git
   ```

### Development Setup

1. Install development tools:
   ```bash
   make dev-setup
   ```

2. Verify your setup:
   ```bash
   make check
   ```

3. Build the project:
   ```bash
   make build
   ```

## Development Workflow

### Branch Strategy

- `main`: Production-ready code, protected branch
- `develop`: Integration branch for features (if using Gitflow)
- `feature/issue-number-description`: Feature development
- `bugfix/issue-number-description`: Bug fixes
- `hotfix/issue-number-description`: Critical production fixes

### Working on Issues

1. **Find or Create an Issue**: All changes should be associated with a GitHub issue
2. **Assign Yourself**: Comment on the issue to claim it
3. **Create a Branch**: 
   ```bash
   git checkout -b feature/123-add-new-validator
   ```
4. **Keep Updated**: Regularly sync with upstream:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

### Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
type(scope): description

[optional body]

[optional footer(s)]
```

**Types:**
- `feat`: New features
- `fix`: Bug fixes
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `test`: Test additions or modifications
- `chore`: Build process or auxiliary tool changes
- `ci`: CI/CD configuration changes

**Examples:**
```
feat(validator): add support for new EU country format
fix(client): handle SOAP fault with missing details
docs: update README with installation instructions
test(validation): add comprehensive VAT format tests
```

## Coding Standards

### Go Code Style

1. **Follow Go conventions**: Use `gofmt`, `goimports`, and `go vet`
2. **Run linter**: All code must pass `golangci-lint`
3. **Documentation**: All public functions must have godoc comments
4. **Error handling**: Use typed errors, handle all errors appropriately
5. **Naming**: Use clear, descriptive names following Go conventions

### Code Organization

```
internal/
├── vies/           # Core VIES API client
├── output/         # Output formatting
└── [package]/      # Other internal packages

cmd/
└── viesquery/      # CLI application

docs/               # Documentation
testdata/           # Test fixtures
```

### Security Guidelines

- **Input validation**: Validate all external inputs
- **No secrets in code**: Use environment variables for sensitive data
- **TLS requirements**: Enforce TLS 1.2+ for all HTTP requests
- **Error messages**: Don't expose sensitive information in error messages

## Testing Requirements

### Test Coverage

- Minimum 80% test coverage for new code
- All public APIs must have tests
- Critical paths require comprehensive test coverage

### Test Types

1. **Unit Tests**: Test individual functions/methods
   ```bash
   go test ./internal/vies -v
   ```

2. **Integration Tests**: Test component interactions
   ```bash
   go test ./internal/vies -tags=integration
   ```

3. **End-to-End Tests**: Full application testing
   ```bash
   make test
   ```

### Test Guidelines

- Use table-driven tests for multiple test cases
- Test both success and error scenarios
- Mock external dependencies (HTTP calls, etc.)
- Use testdata/ directory for test fixtures
- Test names should describe what they test

**Example Test Structure:**
```go
func TestValidateVATNumber(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    bool
        wantErr bool
    }{
        {
            name:    "valid German VAT",
            input:   "DE123456789",
            want:    true,
            wantErr: false,
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ValidateVATNumber(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateVATNumber() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("ValidateVATNumber() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Running Tests

```bash
# All tests
make test

# Specific package
go test ./internal/vies -v

# With coverage
make test-coverage

# Race detection
go test -race ./...

# Benchmarks
make bench
```

## Submitting Changes

### Before Submitting

1. **Rebase on main**:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run full test suite**:
   ```bash
   make check
   ```

3. **Update documentation** if needed

4. **Add/update tests** for your changes

### Pull Request Process

1. **Create Pull Request**: Open a PR against the `main` branch
2. **Fill out PR template**: Provide clear description of changes
3. **Link issues**: Reference related issues using "Closes #123"
4. **Wait for review**: Address reviewer feedback promptly
5. **CI checks**: Ensure all automated checks pass
6. **Squash and merge**: Maintainers will handle the merge

### Pull Request Template

Your PR should include:

- **Description**: What does this change accomplish?
- **Type of Change**: Bug fix, new feature, breaking change, etc.
- **Testing**: How was this tested?
- **Checklist**: Confirm all requirements are met
- **Breaking Changes**: Document any breaking changes
- **Related Issues**: Link to relevant issues

### Review Process

- **Code Review**: At least one maintainer review required
- **Automated Checks**: All CI checks must pass
- **Documentation**: Updates reviewed for accuracy
- **Testing**: Test coverage and quality assessed

## Release Process

### Semantic Versioning

This project follows [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backwards compatible)
- **PATCH**: Bug fixes (backwards compatible)

### Release Workflow

1. **Version Bump**: Update version in relevant files
2. **Changelog**: Update CHANGELOG.md with new version
3. **Tag Creation**: Create annotated git tag
4. **Release Notes**: Document changes and improvements
5. **Binary Distribution**: Automated via GitHub Actions

### Release Checklist

- [ ] All tests pass
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Version numbers updated
- [ ] Release notes prepared
- [ ] Binaries built and tested

## Getting Help

### Communication Channels

- **Issues**: GitHub Issues for bugs and feature requests
- **Discussions**: GitHub Discussions for general questions
- **Security**: See SECURITY.md for security-related issues

### Documentation

- **README.md**: Project overview and basic usage
- **docs/**: Detailed documentation
- **Godoc**: API documentation (`godoc -http=:6060`)

### Common Issues

1. **Build Failures**: Check Go version and run `make dev-setup`
2. **Test Failures**: Ensure dependencies are installed
3. **Linting Errors**: Run `make fmt` and `make lint`
4. **Import Issues**: Verify module path and dependencies

## Recognition

Contributors are recognized in:
- GitHub contributors page
- CHANGELOG.md for significant contributions
- Release notes for major features

## License

By contributing, you agree that your contributions will be licensed under the same license as the project (MIT License).

---

**Thank you for contributing to VIES Query!**

For questions about contributing, please open a GitHub Discussion or reach out to the maintainers.
