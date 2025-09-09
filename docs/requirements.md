# VIES Query Tool - Requirements Specification

## Implementation Status

**Current Version**: 1.0.0  
**Status**: ALL REQUIREMENTS IMPLEMENTED AND TESTED

All functional and technical requirements listed in this document have been successfully implemented. The tool is production-ready and has been tested against the live VIES service.

**Key Achievements:**
- Complete SOAP client with proper XML namespace handling
- All 27 EU country VAT formats supported
- Comprehensive error handling including SOAP fault detection
- Both plain text and JSON output formats
- Full command-line interface with all specified flags
- Robust network handling with configurable timeouts
- Security requirements met (TLS 1.2+, input sanitization)
- Performance targets achieved (< 50MB memory, < 10MB binary)

## Overview

VIES Query is a command-line tool for validating VAT identification numbers using the European Commission's VAT Information Exchange System (VIES) API. The tool provides programmatic access to VAT validation services for all EU member states.

## VIES API Specification

### Service Endpoint
- **WSDL URL**: `https://ec.europa.eu/taxation_customs/vies/checkVatService.wsdl`
- **Service URL**: `https://ec.europa.eu/taxation_customs/vies/services/checkVatService`
- **Protocol**: SOAP 1.1/1.2
- **Method**: `checkVat`

### Request Schema
```xml
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:urn="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
   <soapenv:Header/>
   <soapenv:Body>
      <urn:checkVat>
         <urn:countryCode>?</urn:countryCode>
         <urn:vatNumber>?</urn:vatNumber>
      </urn:checkVat>
   </soapenv:Body>
</soapenv:Envelope>
```

### Response Schema
```xml
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
   <soapenv:Body>
      <ns2:checkVatResponse xmlns:ns2="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
         <ns2:countryCode>string</ns2:countryCode>
         <ns2:vatNumber>string</ns2:vatNumber>
         <ns2:requestDate>date</ns2:requestDate>
         <ns2:valid>boolean</ns2:valid>
         <ns2:name>string</ns2:name>
         <ns2:address>string</ns2:address>
      </ns2:checkVatResponse>
   </soapenv:Body>
</soapenv:Envelope>
```

## EU Member State VAT Number Formats

### Country Codes and Validation Patterns

| Country | Code | Format Pattern | Example |
|---------|------|----------------|---------|
| Austria | AT | ATU + 8 digits | ATU12345678 |
| Belgium | BE | BE0 + 9 digits or BE1 + 9 digits | BE0123456789 |
| Bulgaria | BG | BG + 9 or 10 digits | BG123456789 |
| Croatia | HR | HR + 11 digits | HR12345678901 |
| Cyprus | CY | CY + 8 digits + 1 letter | CY12345678L |
| Czech Republic | CZ | CZ + 8, 9, or 10 digits | CZ12345678 |
| Denmark | DK | DK + 8 digits | DK12345678 |
| Estonia | EE | EE + 9 digits | EE123456789 |
| Finland | FI | FI + 8 digits | FI12345678 |
| France | FR | FR + 2 characters + 9 digits | FR12123456789 |
| Germany | DE | DE + 9 digits | DE123456789 |
| Greece | EL | EL + 9 digits | EL123456789 |
| Hungary | HU | HU + 8 digits | HU12345678 |
| Ireland | IE | IE + 8 characters | IE1234567L |
| Italy | IT | IT + 11 digits | IT12345678901 |
| Latvia | LV | LV + 11 digits | LV12345678901 |
| Lithuania | LT | LT + 9 or 12 digits | LT123456789 |
| Luxembourg | LU | LU + 8 digits | LU12345678 |
| Malta | MT | MT + 8 digits | MT12345678 |
| Netherlands | NL | NL + 9 digits + B + 2 digits | NL123456789B01 |
| Poland | PL | PL + 10 digits | PL1234567890 |
| Portugal | PT | PT + 9 digits | PT123456789 |
| Romania | RO | RO + 2-10 digits | RO12345678 |
| Slovakia | SK | SK + 10 digits | SK1234567890 |
| Slovenia | SI | SI + 8 digits | SI12345678 |
| Spain | ES | ES + 1 character + 7 digits + 1 character | ESA1234567L |
| Sweden | SE | SE + 12 digits | SE123456789012 |

## Functional Requirements

### Core Features
1. **VAT Number Validation**: Query VIES API to validate EU VAT numbers
2. **Format Validation**: Pre-validate VAT number format before API call
3. **Company Information Retrieval**: Return company name and address when available
4. **Error Handling**: Graceful handling of API errors and network failures
5. **Output Formatting**: Support multiple output formats (JSON, plain text)

### Command-Line Interface
```bash
viesquery [flags] <VAT_NUMBER>
```

#### Required Arguments
- `VAT_NUMBER`: Full VAT number including country code (e.g., DE123456789)

#### Optional Flags
- `--format, -f`: Output format (json, plain) [default: plain]
- `--timeout, -t`: Request timeout in seconds [default: 30]
- `--verbose, -v`: Enable verbose logging
- `--help, -h`: Display help information
- `--version`: Display version information

#### Example Usage
```bash
# Basic validation
viesquery DE123456789

# JSON output
viesquery --format json DE123456789

# With custom timeout
viesquery --timeout 10 DE123456789
```

### Output Format Specifications

#### Plain Text Output
```
VAT Number: DE123456789
Status: Valid
Company: Example GmbH
Address: Musterstraße 1, 12345 Berlin, Germany
Request Date: 2025-01-09 12:00:00 UTC
```

#### JSON Output
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

## Error Handling Requirements

### API Error Scenarios
1. **Invalid VAT Number Format**: Pre-validation failure
2. **Service Unavailable**: VIES service downtime
3. **Network Timeout**: Request timeout exceeded
4. **SOAP Fault**: Malformed request or service error
5. **Rate Limiting**: Too many requests from IP address

### Error Response Format
```json
{
  "error": true,
  "message": "Invalid VAT number format for country DE",
  "code": "INVALID_FORMAT",
  "vatNumber": "DE12345"
}
```

### Exit Codes
- `0`: Success
- `1`: Invalid arguments
- `2`: Network/API error
- `3`: Invalid VAT number format
- `4`: Service unavailable

## Network and Security Requirements

### Network Configuration
- **HTTP Client**: Support HTTP/HTTPS protocols
- **Timeout**: Configurable request timeout (default 30 seconds)
- **Retries**: Implement exponential backoff for transient failures
- **User Agent**: Include proper User-Agent header

### Security Considerations
- **TLS Verification**: Enforce certificate validation for HTTPS
- **Input Sanitization**: Sanitize VAT number inputs to prevent injection
- **Rate Limiting**: Respect VIES service rate limits
- **Data Privacy**: Do not log sensitive company information

## Performance Requirements

### Response Time
- **Local Validation**: < 1ms for format validation
- **API Response**: < 5 seconds for VIES API calls
- **Timeout Handling**: Configurable timeout with sensible defaults

### Resource Usage
- **Memory**: < 50MB runtime memory usage
- **Binary Size**: < 10MB compiled binary
- **CPU**: Minimal CPU usage during operation

## Logging Requirements

### Log Levels
- **ERROR**: API failures, network errors, invalid formats
- **WARN**: Timeout warnings, retry attempts
- **INFO**: Successful validations (verbose mode)
- **DEBUG**: Request/response details (verbose mode)

### Log Format
```
2025-01-09T12:00:00Z [INFO] VAT validation successful: DE123456789
2025-01-09T12:00:01Z [ERROR] VIES API error: Service temporarily unavailable
```

## Compatibility Requirements

### Go Version
- **Minimum**: Go 1.21
- **Target**: Go 1.22+

### Operating Systems
- Linux (amd64, arm64)
- macOS (amd64, arm64) 
- Windows (amd64)

### Dependencies
- Standard library preferred
- Minimal external dependencies
- No CGO dependencies

## Quality Requirements

### Testing
- Unit tests for all validation logic
- Integration tests with mock VIES service
- Test coverage > 80%
- Table-driven tests for VAT format validation

### Code Quality
- Go formatting with `gofmt`
- Linting with `golangci-lint`
- Static analysis with `go vet`
- Documentation with `godoc`

### Reliability
- Graceful degradation on service failures
- Proper error propagation
- Resource cleanup on exit
