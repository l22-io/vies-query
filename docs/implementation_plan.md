# VIES Query Tool - Implementation Plan

## Project Architecture

### Package Structure
```
l22.io/viesquery/
├── cmd/viesquery/           # CLI entry point
│   └── main.go             # Application bootstrap and flag parsing
├── internal/vies/          # VIES client implementation
│   ├── client.go           # SOAP client wrapper
│   ├── types.go            # Request/response types
│   ├── validation.go       # VAT number format validation
│   └── countries.go        # EU country code mappings
├── internal/output/        # Output formatting
│   ├── formatter.go        # Output format interface
│   ├── plain.go           # Plain text formatter
│   └── json.go            # JSON formatter
├── pkg/vies/              # Public API (if needed for library usage)
├── docs/                  # Documentation
├── testdata/              # Test fixtures and mock responses
└── scripts/               # Build and development scripts
```

### Module Dependencies
```go
module l22.io/viesquery

go 1.22

require (
    // Standard library only for initial version
    // Potential future dependencies:
    // github.com/spf13/cobra v1.8.0  // CLI framework (optional)
)
```

## SOAP Client Implementation Strategy

### Code Generation Approach
Since Go doesn't have native WSDL-to-Go generation like Java or .NET, we'll implement a lightweight SOAP client manually rather than using heavy code generation tools.

### Manual SOAP Client Implementation
```go
// internal/vies/types.go
type CheckVatRequest struct {
    CountryCode string
    VatNumber   string
}

type CheckVatResponse struct {
    CountryCode string    `xml:"countryCode"`
    VatNumber   string    `xml:"vatNumber"`
    RequestDate time.Time `xml:"requestDate"`
    Valid       bool      `xml:"valid"`
    Name        string    `xml:"name"`
    Address     string    `xml:"address"`
}

type SOAPEnvelope struct {
    XMLName xml.Name `xml:"soap:Envelope"`
    Body    SOAPBody `xml:"soap:Body"`
}

type SOAPBody struct {
    CheckVat         *CheckVatRequest  `xml:"checkVat,omitempty"`
    CheckVatResponse *CheckVatResponse `xml:"checkVatResponse,omitempty"`
    Fault            *SOAPFault        `xml:"soap:Fault,omitempty"`
}
```

### HTTP Client Configuration
```go
// internal/vies/client.go
type Client struct {
    httpClient *http.Client
    endpoint   string
    timeout    time.Duration
    userAgent  string
    logger     Logger
}

func NewClient(options ...ClientOption) *Client {
    return &Client{
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
            Transport: &http.Transport{
                TLSClientConfig: &tls.Config{
                    MinVersion: tls.VersionTLS12,
                },
            },
        },
        endpoint:  "https://ec.europa.eu/taxation_customs/vies/services/checkVatService",
        userAgent: "viesquery/1.0.0",
    }
}
```

## VAT Number Validation Logic

### Pre-validation Strategy
Implement client-side format validation before making API calls to reduce unnecessary network requests.

```go
// internal/vies/validation.go
type CountryValidator struct {
    Code    string
    Pattern *regexp.Regexp
    Length  []int
}

var countryValidators = map[string]CountryValidator{
    "DE": {
        Code:    "DE",
        Pattern: regexp.MustCompile(`^DE\d{9}$`),
        Length:  []int{11}, // DE + 9 digits
    },
    "AT": {
        Code:    "AT",
        Pattern: regexp.MustCompile(`^ATU\d{8}$`),
        Length:  []int{11}, // ATU + 8 digits
    },
    // ... more country validators
}

func ValidateFormat(vatNumber string) error {
    if len(vatNumber) < 3 {
        return errors.New("VAT number too short")
    }
    
    countryCode := vatNumber[:2]
    validator, exists := countryValidators[countryCode]
    if !exists {
        return fmt.Errorf("unsupported country code: %s", countryCode)
    }
    
    if !validator.Pattern.MatchString(vatNumber) {
        return fmt.Errorf("invalid format for country %s", countryCode)
    }
    
    return nil
}
```

### Request Processing Workflow
1. Parse command-line arguments
2. Pre-validate VAT number format
3. Extract country code and VAT number
4. Create SOAP request
5. Send HTTP request with retry logic
6. Parse SOAP response
7. Format and display results

## CLI Application Design

### Command-Line Interface
```go
// cmd/viesquery/main.go
func main() {
    var (
        format    = flag.String("format", "plain", "Output format (plain, json)")
        timeout   = flag.Int("timeout", 30, "Request timeout in seconds")
        verbose   = flag.Bool("verbose", false, "Enable verbose logging")
        version   = flag.Bool("version", false, "Display version information")
    )
    
    flag.Usage = func() {
        fmt.Fprintf(os.Stderr, "Usage: %s [flags] VAT_NUMBER\n", os.Args[0])
        fmt.Fprintf(os.Stderr, "\nValidate EU VAT numbers using VIES API\n\n")
        fmt.Fprintf(os.Stderr, "Arguments:\n")
        fmt.Fprintf(os.Stderr, "  VAT_NUMBER    EU VAT number to validate (e.g., DE123456789)\n\n")
        fmt.Fprintf(os.Stderr, "Flags:\n")
        flag.PrintDefaults()
    }
    
    flag.Parse()
    
    if *version {
        fmt.Printf("viesquery version %s\n", Version)
        return
    }
    
    if flag.NArg() != 1 {
        flag.Usage()
        os.Exit(1)
    }
    
    vatNumber := flag.Arg(0)
    
    client := vies.NewClient(
        vies.WithTimeout(time.Duration(*timeout) * time.Second),
        vies.WithVerbose(*verbose),
    )
    
    result, err := client.CheckVAT(context.Background(), vatNumber)
    if err != nil {
        handleError(err, *format)
        return
    }
    
    displayResult(result, *format)
}
```

### Error Handling Strategy
```go
type ViesError struct {
    Code    string
    Message string
    VATNumber string
}

func (e *ViesError) Error() string {
    return fmt.Sprintf("%s: %s (VAT: %s)", e.Code, e.Message, e.VATNumber)
}

const (
    ErrInvalidFormat    = "INVALID_FORMAT"
    ErrServiceError     = "SERVICE_ERROR" 
    ErrNetworkTimeout   = "NETWORK_TIMEOUT"
    ErrServiceUnavailable = "SERVICE_UNAVAILABLE"
)
```

## Output Formatting System

### Formatter Interface
```go
// internal/output/formatter.go
type Formatter interface {
    Format(result *vies.CheckVatResult) (string, error)
    FormatError(err error) (string, error)
}

type Manager struct {
    formatters map[string]Formatter
}

func NewManager() *Manager {
    return &Manager{
        formatters: map[string]Formatter{
            "plain": NewPlainFormatter(),
            "json":  NewJSONFormatter(),
        },
    }
}
```

### Plain Text Formatter
```go
// internal/output/plain.go
type PlainFormatter struct{}

func (f *PlainFormatter) Format(result *vies.CheckVatResult) (string, error) {
    var b strings.Builder
    
    fmt.Fprintf(&b, "VAT Number: %s%s\n", result.CountryCode, result.VatNumber)
    fmt.Fprintf(&b, "Status: %s\n", formatStatus(result.Valid))
    
    if result.Valid {
        if result.Name != "" {
            fmt.Fprintf(&b, "Company: %s\n", result.Name)
        }
        if result.Address != "" {
            fmt.Fprintf(&b, "Address: %s\n", result.Address)
        }
    }
    
    fmt.Fprintf(&b, "Request Date: %s\n", result.RequestDate.Format("2006-01-02 15:04:05 UTC"))
    
    return b.String(), nil
}
```

## Testing Strategy

### Unit Testing Approach
```go
// internal/vies/validation_test.go
func TestValidateFormat(t *testing.T) {
    tests := []struct {
        name      string
        vatNumber string
        wantErr   bool
        errMsg    string
    }{
        {
            name:      "Valid German VAT",
            vatNumber: "DE123456789",
            wantErr:   false,
        },
        {
            name:      "Invalid German VAT - too short",
            vatNumber: "DE12345",
            wantErr:   true,
            errMsg:    "invalid format for country DE",
        },
        {
            name:      "Unsupported country",
            vatNumber: "XX123456789",
            wantErr:   true,
            errMsg:    "unsupported country code: XX",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateFormat(tt.vatNumber)
            if tt.wantErr {
                if err == nil {
                    t.Error("expected error, got nil")
                }
                if !strings.Contains(err.Error(), tt.errMsg) {
                    t.Errorf("expected error containing %q, got %q", tt.errMsg, err.Error())
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
            }
        })
    }
}
```

### Integration Testing with Mock Server
```go
// internal/vies/client_test.go
func TestClient_CheckVAT_Integration(t *testing.T) {
    // Create mock SOAP server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Validate request
        body, _ := io.ReadAll(r.Body)
        if !strings.Contains(string(body), "DE123456789") {
            t.Error("request doesn't contain expected VAT number")
        }
        
        // Return mock response
        response := `<?xml version="1.0" encoding="UTF-8"?>
        <soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
            <soap:Body>
                <ns2:checkVatResponse xmlns:ns2="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
                    <ns2:countryCode>DE</ns2:countryCode>
                    <ns2:vatNumber>123456789</ns2:vatNumber>
                    <ns2:requestDate>2025-01-09T12:00:00Z</ns2:requestDate>
                    <ns2:valid>true</ns2:valid>
                    <ns2:name>Test Company</ns2:name>
                    <ns2:address>Test Address</ns2:address>
                </ns2:checkVatResponse>
            </soap:Body>
        </soap:Envelope>`
        w.Header().Set("Content-Type", "text/xml")
        w.Write([]byte(response))
    }))
    defer server.Close()
    
    client := NewClient(WithEndpoint(server.URL))
    result, err := client.CheckVAT(context.Background(), "DE123456789")
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if !result.Valid {
        t.Error("expected valid result")
    }
    
    if result.Name != "Test Company" {
        t.Errorf("expected company name 'Test Company', got %q", result.Name)
    }
}
```

### Test Data Organization
```
testdata/
├── soap_responses/
│   ├── valid_response.xml
│   ├── invalid_response.xml
│   ├── soap_fault.xml
│   └── service_unavailable.xml
├── vat_numbers/
│   ├── valid_numbers.json
│   └── invalid_numbers.json
└── mock_server/
    └── responses.go
```

## Build and Release Process

### Makefile Configuration
```makefile
# Makefile
.PHONY: build test clean lint fmt install

VERSION ?= $(shell git describe --tags --always --dirty)
LDFLAGS = -ldflags "-X main.Version=$(VERSION) -s -w"

build:
	go build $(LDFLAGS) -o bin/viesquery ./cmd/viesquery

test:
	go test -v -race -cover ./...

bench:
	go test -bench=. -benchmem ./...

lint:
	golangci-lint run

fmt:
	go fmt ./...
	gofmt -s -w .

clean:
	rm -rf bin/

install: build
	cp bin/viesquery $(GOPATH)/bin/

release: clean test lint
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/viesquery-linux-amd64 ./cmd/viesquery
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/viesquery-linux-arm64 ./cmd/viesquery
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/viesquery-darwin-amd64 ./cmd/viesquery
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/viesquery-darwin-arm64 ./cmd/viesquery
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/viesquery-windows-amd64.exe ./cmd/viesquery
```

### Git Hooks Setup
```bash
# .git/hooks/pre-commit
#!/bin/bash
set -e

echo "Running pre-commit checks..."

# Format check
if ! gofmt -l . | grep -q '^$'; then
    echo "Code is not formatted. Run 'make fmt'"
    exit 1
fi

# Lint check
if ! golangci-lint run; then
    echo "Linting failed"
    exit 1
fi

# Tests
if ! go test -short ./...; then
    echo "Tests failed"
    exit 1
fi

echo "Pre-commit checks passed!"
```

## Performance Optimization

### HTTP Client Optimizations
- Connection pooling with Keep-Alive
- Timeout configuration per request type
- Retry with exponential backoff
- Request/response compression support

### Memory Management
- Minimize allocations in hot paths
- Reuse buffers for SOAP encoding/decoding
- Proper resource cleanup in defer statements

### Caching Strategy
Consider implementing optional caching for repeated queries:
```go
type Cache interface {
    Get(vatNumber string) (*CheckVatResult, bool)
    Set(vatNumber string, result *CheckVatResult, ttl time.Duration)
}

// In-memory cache for development/testing
type MemoryCache struct {
    data   map[string]cacheEntry
    mutex  sync.RWMutex
}
```

## Development Workflow

### Phase 1: Core Implementation
1. Implement VAT format validation
2. Create basic SOAP client
3. Add plain text output formatting
4. Write comprehensive unit tests

### Phase 2: Enhanced Features  
1. Add JSON output formatting
2. Implement retry logic with exponential backoff
3. Add timeout configuration
4. Create integration tests with mock server

### Phase 3: Production Readiness
1. Add proper logging with configurable levels
2. Implement error handling and exit codes
3. Add build scripts and release automation
4. Write comprehensive documentation

### Phase 4: Optional Enhancements
1. Add caching support
2. Implement batch processing
3. Add configuration file support
4. Performance optimizations

## Versioning and Release Strategy

### Semantic Versioning
- `v0.x.x`: Development versions
- `v1.0.0`: First stable release
- Patch versions for bug fixes
- Minor versions for new features
- Major versions for breaking changes

### Release Process
1. Update version in code
2. Run full test suite
3. Update CHANGELOG.md
4. Create git tag
5. Build multi-platform binaries
6. Create GitHub release with artifacts
7. Update documentation

This implementation plan provides a solid foundation for building a robust, maintainable VIES query tool that follows Go best practices and meets all the technical requirements.
