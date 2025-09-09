# VIES Query

Command-line tool for validating European Union VAT identification numbers using the VIES (VAT Information Exchange System) API.

## Overview

VIES Query provides fast, reliable validation of EU VAT numbers with support for all 27 EU member states. The tool performs client-side format validation and queries the official European Commission VIES service for real-time verification.

## Features

- **Complete EU Coverage**: Supports all 27 EU member states
- **Format Validation**: Client-side validation before API calls  
- **Multiple Output Formats**: Plain text and JSON output
- **Robust Error Handling**: Detailed error messages and appropriate exit codes
- **Network Resilience**: Configurable timeouts and proper TLS security
- **Professional Output**: Clean, parseable results for automation

## Installation

### Prerequisites

- Go 1.21 or later
- Internet connection for API queries

### Install from Source

```bash
go install l22.io/viesquery/cmd/viesquery@latest
```

### Build from Repository

```bash
git clone https://github.com/l22-io/vies-query.git
cd vies-query
make build
```

### Download Binary

Download pre-built binaries from the [releases page](https://github.com/l22-io/vies-query/releases).

## Quick Start

### Basic Usage

```bash
# Validate a German VAT number
viesquery DE123456789

# JSON output format
viesquery --format json DE123456789

# Custom timeout and verbose logging
viesquery --timeout 60 --verbose IT12345678901
```

### Example Output

**Plain Text:**
```
VAT Number: DE123456789
Status: Valid
Company: Example GmbH
Address: Musterstraße 1, 12345 Berlin, Germany
Request Date: 2025-01-09 12:00:00 UTC
```

**JSON:**
```json
{
  "countryCode": "DE",
  "vatNumber": "123456789",
  "requestDate": "2025-01-09T12:00:00Z",
  "valid": true,
  "name": "Example GmbH",
  "address": "Musterstraße 1, 12345 Berlin, Germany"
}
```

## Supported Countries

| Country | Code | Format Example |
|---------|------|--------------  |
| Austria | AT | ATU12345678 |
| Belgium | BE | BE0123456789 |
| Bulgaria | BG | BG123456789 |
| Croatia | HR | HR12345678901 |
| Cyprus | CY | CY12345678L |
| Czech Republic | CZ | CZ12345678 |
| Denmark | DK | DK12345678 |
| Estonia | EE | EE123456789 |
| Finland | FI | FI12345678 |
| France | FR | FR12123456789 |
| Germany | DE | DE123456789 |
| Greece | EL | EL123456789 |
| Hungary | HU | HU12345678 |
| Ireland | IE | IE1234567L |
| Italy | IT | IT12345678901 |
| Latvia | LV | LV12345678901 |
| Lithuania | LT | LT123456789 |
| Luxembourg | LU | LU12345678 |
| Malta | MT | MT12345678 |
| Netherlands | NL | NL123456789B01 |
| Poland | PL | PL1234567890 |
| Portugal | PT | PT123456789 |
| Romania | RO | RO12345678 |
| Slovakia | SK | SK1234567890 |
| Slovenia | SI | SI12345678 |
| Spain | ES | ESA1234567L |
| Sweden | SE | SE123456789012 |

## Command-Line Options

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--format` | `-f` | `plain` | Output format (plain, json) |
| `--timeout` | `-t` | `30` | Request timeout in seconds |
| `--verbose` | `-v` | `false` | Enable verbose logging |
| `--help` | `-h` | - | Display help information |
| `--version` | - | - | Display version information |

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `VIESQUERY_FORMAT` | Default output format | `plain` |
| `VIESQUERY_TIMEOUT` | Default timeout in seconds | `30` |
| `VIESQUERY_VERBOSE` | Enable verbose mode | `false` |

## Exit Codes

- `0`: Successful validation
- `1`: Invalid command arguments
- `2`: Network or API error  
- `3`: Invalid VAT number format
- `4`: VIES service unavailable

## Advanced Usage

### Batch Processing

```bash
#!/bin/bash
vat_numbers=("DE123456789" "AT12345678" "FR12123456789")

for vat in "${vat_numbers[@]}"; do
    echo "Validating $vat..."
    viesquery --format json "$vat"
done
```

### CI/CD Integration

```bash
# Validate VAT and exit on failure
if ! viesquery "$VAT_NUMBER" > /dev/null 2>&1; then
    echo "Invalid VAT number: $VAT_NUMBER"
    exit 1
fi
```

### JSON Processing

```bash
# Extract company name
viesquery --format json DE123456789 | jq -r '.name'

# Check validity
viesquery --format json DE123456789 | jq -r '.valid'
```

## Development

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make release

# Run tests
make test

# Format and lint code
make check
```

### Project Structure

```
l22.io/viesquery/
├── cmd/viesquery/           # CLI application
├── internal/vies/           # VIES client and validation
├── internal/output/         # Output formatting
├── docs/                    # Documentation
└── testdata/               # Test fixtures
```

## Performance and Rate Limits

### VIES Service Limitations

- Intended for occasional verification, not bulk processing
- Rate limits enforced (approximately 1 request/second recommended)
- Service availability typically 99%+ uptime
- Scheduled maintenance on weekends and EU holidays

### Best Practices

1. **Cache Results**: VAT numbers change infrequently
2. **Validate Format First**: Avoid unnecessary API calls
3. **Implement Retries**: Handle temporary service issues
4. **Respect Rate Limits**: Avoid IP blocking
5. **Monitor Usage**: Track validation volumes

## Documentation

- [Requirements Specification](docs/requirements.md)
- [Implementation Plan](docs/implementation_plan.md)  
- [User Guide](docs/user_guide.md)

## Contributing

Contributions are welcome! Please see the [implementation plan](docs/implementation_plan.md) for development guidelines.

### Reporting Issues

When reporting issues, include:

- Tool version (`viesquery --version`)
- Operating system and architecture
- Complete command used
- Error messages (with `--verbose` if applicable)
- Expected vs actual behavior

## Legal Notice

### Data Privacy

- Only sends VAT numbers to the official EU VIES service
- No data stored locally or sent to third parties
- Company information returned directly from VIES
- Verbose mode may log request/response data locally

### Service Disclaimer  

- VIES service operated by the European Commission
- Service availability and accuracy not guaranteed
- Always verify critical VAT information through multiple sources
- Use in compliance with applicable laws and regulations

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- Check the [User Guide](docs/user_guide.md) for detailed usage instructions
- Review [GitHub Issues](https://github.com/l22-io/vies-query/issues) for known problems
- Consult the [requirements document](docs/requirements.md) for technical details

---

**VIES Query** - Professional VAT validation for the command line  
© 2025 l22.io - Made with Go
