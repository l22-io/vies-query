# VIES Query Tool - User Guide

## Overview

VIES Query is a command-line tool for validating European Union VAT identification numbers using the European Commission's VIES (VAT Information Exchange System) API. The tool provides fast, reliable VAT validation with support for all EU member states.

**Latest Updates:**
- Fixed XML namespace conflicts in SOAP requests
- Added comprehensive VIES API documentation based on official WSDL
- Improved error handling with proper SOAP fault detection
- Added support for xsd:date format parsing
- Enhanced reliability with proper namespace handling

## Installation

### Prerequisites

- Go 1.21 or later
- Internet connection for API queries
- Terminal/command prompt access

### Installation Methods

#### Method 1: Install from Source
```bash
go install l22.io/viesquery/cmd/viesquery@latest
```

#### Method 2: Build from Repository
```bash
git clone https://github.com/l22-io/vies-query.git
cd vies-query
go build -o viesquery cmd/viesquery/main.go
```

#### Method 3: Download Pre-built Binary
Download the appropriate binary for your platform from the [releases page](https://github.com/l22-io/vies-query/releases).

### Verify Installation
```bash
viesquery --version
```

## Basic Usage

### Command Syntax
```bash
viesquery [flags] VAT_NUMBER
```

### Required Arguments
- `VAT_NUMBER`: Complete VAT number including country code (e.g., `DE123456789`)

### Available Flags
| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--format` | `-f` | `plain` | Output format (`plain`, `json`) |
| `--timeout` | `-t` | `30` | Request timeout in seconds |
| `--verbose` | `-v` | `false` | Enable verbose logging |
| `--help` | `-h` | - | Display help information |
| `--version` | - | - | Display version information |

## Examples

### Basic Validation
Validate a German VAT number with plain text output:
```bash
viesquery DE123456789
```

Output:
```
VAT Number: DE123456789
Status: Valid
Company: Beispiel GmbH
Address: Musterstraße 1, 12345 Berlin, Germany
Request Date: 2025-01-09 12:00:00 UTC
```

### JSON Output Format
Get validation results in JSON format:
```bash
viesquery --format json DE123456789
```

Output:
```json
{
  "countryCode": "DE",
  "vatNumber": "123456789",
  "requestDate": "2025-01-09T12:00:00Z",
  "valid": true,
  "name": "Beispiel GmbH",
  "address": "Musterstraße 1, 12345 Berlin, Germany"
}
```

### Custom Timeout
Set a custom timeout for slow network connections:
```bash
viesquery --timeout 60 IT12345678901
```

### Verbose Logging
Enable detailed logging for debugging:
```bash
viesquery --verbose --format json FR12123456789
```

### Invalid VAT Number
Example of validation failure:
```bash
viesquery DE12345
```

Output:
```
Error: Invalid VAT number format for country DE
VAT Number: DE12345
Expected Format: DE + 9 digits
```

## Supported Countries and Formats

### EU Member States

| Country | Code | Format | Example |
|---------|------|--------|---------|
| Austria | AT | ATU + 8 digits | `ATU12345678` |
| Belgium | BE | BE0/BE1 + 9 digits | `BE0123456789` |
| Bulgaria | BG | BG + 9-10 digits | `BG123456789` |
| Croatia | HR | HR + 11 digits | `HR12345678901` |
| Cyprus | CY | CY + 8 digits + letter | `CY12345678L` |
| Czech Republic | CZ | CZ + 8-10 digits | `CZ12345678` |
| Denmark | DK | DK + 8 digits | `DK12345678` |
| Estonia | EE | EE + 9 digits | `EE123456789` |
| Finland | FI | FI + 8 digits | `FI12345678` |
| France | FR | FR + 2 chars + 9 digits | `FR12123456789` |
| Germany | DE | DE + 9 digits | `DE123456789` |
| Greece | EL | EL + 9 digits | `EL123456789` |
| Hungary | HU | HU + 8 digits | `HU12345678` |
| Ireland | IE | IE + 8 characters | `IE1234567L` |
| Italy | IT | IT + 11 digits | `IT12345678901` |
| Latvia | LV | LV + 11 digits | `LV12345678901` |
| Lithuania | LT | LT + 9/12 digits | `LT123456789` |
| Luxembourg | LU | LU + 8 digits | `LU12345678` |
| Malta | MT | MT + 8 digits | `MT12345678` |
| Netherlands | NL | NL + 9 digits + B + 2 digits | `NL123456789B01` |
| Poland | PL | PL + 10 digits | `PL1234567890` |
| Portugal | PT | PT + 9 digits | `PT123456789` |
| Romania | RO | RO + 2-10 digits | `RO12345678` |
| Slovakia | SK | SK + 10 digits | `SK1234567890` |
| Slovenia | SI | SI + 8 digits | `SI12345678` |
| Spain | ES | ES + char + 7 digits + char | `ESA1234567L` |
| Sweden | SE | SE + 12 digits | `SE123456789012` |

## Error Handling

### Exit Codes
- `0`: Successful validation
- `1`: Invalid command arguments
- `2`: Network or API error
- `3`: Invalid VAT number format
- `4`: VIES service unavailable

### Common Error Messages

#### Invalid Format
```
Error: Invalid VAT number format for country DE
VAT Number: DE12345
Expected Format: DE + 9 digits
```

#### Network Timeout
```
Error: Request timeout after 30 seconds
VAT Number: DE123456789
Try increasing timeout with --timeout flag
```

#### SOAP Fault Errors
The tool now properly handles VIES service-specific errors:

```
Error: SOAP fault: env:Server - MS_UNAVAILABLE
```

Common SOAP fault types:
- `MS_UNAVAILABLE`: Member State service not available
- `MS_MAX_CONCURRENT_REQ`: Too many concurrent requests for this country
- `GLOBAL_MAX_CONCURRENT_REQ`: Global rate limit exceeded
- `SERVICE_UNAVAILABLE`: General service error
- `TIMEOUT`: Request timeout
- `INVALID_INPUT`: Invalid country code or VAT number format

#### Service Unavailable
```
Error: VIES service temporarily unavailable
VAT Number: DE123456789
Please retry later or check VIES service status
```

#### Unsupported Country
```
Error: Unsupported country code: XX
VAT Number: XX123456789
Supported countries: AT, BE, BG, HR, CY, CZ, DK, EE, FI, FR, DE, EL, HU, IE, IT, LV, LT, LU, MT, NL, PL, PT, RO, SK, SI, ES, SE
```

## Advanced Usage

### Batch Processing with Shell Scripts
Process multiple VAT numbers using shell scripting:

```bash
#!/bin/bash
# validate_vats.sh

vat_numbers=(
    "DE123456789"
    "AT12345678"
    "FR12123456789"
)

for vat in "${vat_numbers[@]}"; do
    echo "Validating $vat..."
    viesquery --format json "$vat"
    echo "---"
done
```

### Integration with CI/CD
Use VIES Query in automated workflows:

```bash
# Validate VAT number and exit on failure
if ! viesquery "$VAT_NUMBER" > /dev/null 2>&1; then
    echo "Invalid VAT number: $VAT_NUMBER"
    exit 1
fi
echo "VAT validation successful"
```

### JSON Processing with jq
Extract specific fields from JSON output:

```bash
# Get only the company name
viesquery --format json DE123456789 | jq -r '.name'

# Check if VAT is valid (returns true/false)
viesquery --format json DE123456789 | jq -r '.valid'

# Get formatted output
viesquery --format json DE123456789 | jq -r '"Company: \(.name), Valid: \(.valid)"'
```

## Performance Considerations

### Network Optimization
- Default timeout is 30 seconds (configurable)
- Single HTTP connection per request
- Automatic retry with exponential backoff
- TLS 1.2+ enforced for security

### Rate Limiting
The VIES service has rate limiting in place:
- Recommended: Max 1 request per second
- Burst limit: ~10 requests before throttling
- Daily limit: Varies by IP address

### Caching Results
For applications processing the same VAT numbers repeatedly, consider implementing client-side caching as the VIES service does not change VAT validity frequently.

## Troubleshooting

### Common Issues

#### "Command not found: viesquery"
**Solution**: Ensure the binary is in your PATH or use the full path to the executable.

```bash
# Add to PATH (bash/zsh)
export PATH=$PATH:/path/to/viesquery

# Or use full path
/path/to/viesquery DE123456789
```

#### "Connection refused" or "Network unreachable"
**Solution**: Check your internet connection and firewall settings.

```bash
# Test connectivity to VIES service
curl -I https://ec.europa.eu/taxation_customs/vies/

# Use verbose mode for detailed error information
viesquery --verbose DE123456789
```

#### "Invalid format" for correct-looking VAT numbers
**Solution**: Verify the VAT number format against the country-specific rules. Some countries have strict formatting requirements.

```bash
# Check exact format requirements
viesquery --help | grep -A 30 "Supported Countries"
```

#### Timeout errors on slow connections
**Solution**: Increase the timeout value.

```bash
# Increase timeout to 60 seconds
viesquery --timeout 60 DE123456789
```

### VIES Service Status

The VIES service has scheduled maintenance periods and occasional outages:

- **Scheduled Maintenance**: Usually weekends and EU holidays
- **Service Status**: Check [VIES availability](https://ec.europa.eu/taxation_customs/vies/)
- **Alternative**: Some countries provide their own VAT validation services

### Debug Mode

Use verbose logging to diagnose issues:

```bash
viesquery --verbose --format json DE123456789
```

This will show:
- Request details and headers
- Response parsing information
- Timing information
- Error details

## Environment Variables

The tool respects the following environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `VIESQUERY_TIMEOUT` | Default timeout in seconds | `30` |
| `VIESQUERY_FORMAT` | Default output format | `plain` |
| `VIESQUERY_VERBOSE` | Enable verbose mode | `false` |

Example:
```bash
export VIESQUERY_TIMEOUT=60
export VIESQUERY_FORMAT=json
viesquery DE123456789
```

## API Rate Limits and Best Practices

### VIES Service Limitations
- The VIES service is provided for occasional verification of VAT numbers
- Not intended for bulk verification or systematic downloading
- Rate limits are enforced but not officially documented
- Service availability varies (typically 99%+ uptime)

### Best Practices
1. **Cache Results**: VAT numbers don't change frequently
2. **Implement Retries**: Handle temporary service unavailability
3. **Validate Format First**: Use local validation before API calls
4. **Monitor Usage**: Respect rate limits to avoid IP blocking
5. **Handle Errors Gracefully**: Provide fallback mechanisms

## Support and Contributing

### Getting Help
- Check the [troubleshooting section](#troubleshooting) first
- Review the [requirements document](requirements.md) for technical details
- Check the [VIES API specification](vies-api-specification.md) for detailed API documentation
- Review the [namespace fix summary](namespace-fix-summary.md) for recent technical improvements
- Check [GitHub issues](https://github.com/l22-io/vies-query/issues) for known problems

### Reporting Issues
When reporting issues, include:
- Tool version (`viesquery --version`)
- Operating system and architecture
- Complete command used
- Error messages (with `--verbose` if applicable)
- Expected vs actual behavior

### Contributing
See the [implementation plan](implementation_plan.md) for development guidelines and contribution instructions.

## Legal and Privacy Notes

### Data Privacy
- The tool only sends VAT numbers to the official EU VIES service
- Company names and addresses are returned from the VIES service
- No data is stored locally or transmitted to third parties
- Verbose mode may log request/response data locally

### Service Availability
- The VIES service is operated by the European Commission
- Service availability and accuracy are not guaranteed
- Always verify critical VAT information through multiple sources
- Some member states may have different data availability levels

### Compliance
- Use this tool in compliance with applicable laws and regulations
- Respect the VIES service terms of use
- Implement appropriate data protection measures for any stored results
