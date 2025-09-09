package vies

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	defaultEndpoint  = "https://ec.europa.eu/taxation_customs/vies/services/checkVatService"
	defaultUserAgent = "viesquery/1.0.0"
	soapNamespace    = "urn:ec.europa.eu:taxud:vies:services:checkVat:types"
)

// Client represents a VIES API client
type Client struct {
	httpClient *http.Client
	endpoint   string
	userAgent  string
	verbose    bool
	logger     *log.Logger
}

// NewClient creates a new VIES client with the given options
func NewClient(options ...ClientOption) *Client {
	opts := &ClientOptions{
		Timeout:   30 * time.Second,
		UserAgent: defaultUserAgent,
		Verbose:   false,
		Endpoint:  defaultEndpoint,
	}

	// Apply options
	for _, option := range options {
		option(opts)
	}

	// Create HTTP client with security settings
	client := &Client{
		httpClient: &http.Client{
			Timeout: opts.Timeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					MinVersion: tls.VersionTLS12,
				},
				DisableKeepAlives:     false,
				MaxIdleConns:          10,
				MaxIdleConnsPerHost:   2,
				IdleConnTimeout:       30 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
			},
		},
		endpoint:  opts.Endpoint,
		userAgent: opts.UserAgent,
		verbose:   opts.Verbose,
		logger:    log.New(os.Stderr, "[VIES] ", log.LstdFlags),
	}

	return client
}

// CheckVAT validates a VAT number using the VIES service
func (c *Client) CheckVAT(ctx context.Context, vatNumber string) (*CheckVatResult, error) {
	startTime := time.Now()

	if c.verbose {
		c.logger.Printf("Validating VAT number: %s", vatNumber)
	}

	// Parse and validate VAT number format
	countryCode, number, err := ParseVATNumber(vatNumber)
	if err != nil {
		return nil, err
	}

	if c.verbose {
		c.logger.Printf("Parsed VAT: Country=%s, Number=%s", countryCode, number)
	}

	// Create SOAP request
	soapRequest := createSOAPRequest(countryCode, number)
	
	// Marshal to XML
	requestBody, err := xml.Marshal(soapRequest)
	if err != nil {
		return nil, &ServiceError{
			Code:      ErrServiceError,
			Message:   fmt.Sprintf("Failed to create SOAP request: %v", err),
			VATNumber: vatNumber,
		}
	}

	// Add XML declaration
	fullRequest := []byte(xml.Header + string(requestBody))

	if c.verbose {
		c.logger.Printf("SOAP Request: %s", string(fullRequest))
	}

	// Send HTTP request
	result, err := c.sendSOAPRequest(ctx, fullRequest)
	if err != nil {
		return nil, err
	}

	// Set original VAT number for display
	result.VatNumber = number
	result.CountryCode = countryCode

	duration := time.Since(startTime)
	if c.verbose {
		c.logger.Printf("Validation completed in %v. Valid: %t", duration, result.Valid)
	}

	return result, nil
}

// sendSOAPRequest sends a SOAP request and parses the response
func (c *Client) sendSOAPRequest(ctx context.Context, requestBody []byte) (*CheckVatResult, error) {
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return nil, &ServiceError{
			Code:    ErrServiceError,
			Message: fmt.Sprintf("Failed to create HTTP request: %v", err),
		}
	}

	// Set headers
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "checkVat")
	req.Header.Set("User-Agent", c.userAgent)

	if c.verbose {
		c.logger.Printf("Sending request to: %s", c.endpoint)
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return nil, &ServiceError{
				Code:    ErrNetworkTimeout,
				Message: "Request timeout exceeded",
			}
		}
		return nil, &ServiceError{
			Code:    ErrServiceError,
			Message: fmt.Sprintf("HTTP request failed: %v", err),
		}
	}
	defer resp.Body.Close()

	// Read response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &ServiceError{
			Code:    ErrServiceError,
			Message: fmt.Sprintf("Failed to read response body: %v", err),
		}
	}

	if c.verbose {
		c.logger.Printf("Response Status: %s", resp.Status)
		c.logger.Printf("Response Body: %s", string(responseBody))
	}

	// Check HTTP status
	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusServiceUnavailable ||
			resp.StatusCode == http.StatusBadGateway ||
			resp.StatusCode == http.StatusGatewayTimeout {
			return nil, &ServiceError{
				Code:    ErrServiceUnavailable,
				Message: "VIES service is temporarily unavailable",
			}
		}
		return nil, &ServiceError{
			Code:    ErrServiceError,
			Message: fmt.Sprintf("HTTP error: %s", resp.Status),
		}
	}

	// Parse SOAP response
	return c.parseSOAPResponse(responseBody)
}

// parseSOAPResponse parses the SOAP response from VIES
func (c *Client) parseSOAPResponse(responseBody []byte) (*CheckVatResult, error) {
	// Use inline structs to avoid namespace conflicts
	var envelope struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    struct {
			CheckVatResponse *struct {
				XMLName     xml.Name `xml:"checkVatResponse"`
				CountryCode string   `xml:"countryCode"`
				VatNumber   string   `xml:"vatNumber"`
				RequestDate string   `xml:"requestDate"` // Parse as string first, then convert to time.Time
				Valid       bool     `xml:"valid"`
				Name        string   `xml:"name"`
				Address     string   `xml:"address"`
			} `xml:"checkVatResponse"`
			Fault *struct {
				XMLName xml.Name `xml:"Fault"`
				Code    string   `xml:"faultcode"`
				String  string   `xml:"faultstring"`
			} `xml:"Fault"`
		} `xml:"Body"`
	}

	err := xml.Unmarshal(responseBody, &envelope)
	if err != nil {
		return nil, &ServiceError{
			Code:    ErrServiceError,
			Message: fmt.Sprintf("Failed to parse SOAP response: %v", err),
		}
	}

	// Check for SOAP fault
	if envelope.Body.Fault != nil {
		return nil, &ServiceError{
			Code:    ErrSOAPFault,
			Message: fmt.Sprintf("SOAP fault: %s - %s", envelope.Body.Fault.Code, envelope.Body.Fault.String),
		}
	}

	// Check for valid response
	if envelope.Body.CheckVatResponse == nil {
		return nil, &ServiceError{
			Code:    ErrServiceError,
			Message: "Invalid SOAP response: missing checkVatResponse",
		}
	}

	resp := envelope.Body.CheckVatResponse

	// Parse request date (xsd:date format: YYYY-MM-DD)
	requestDate, err := time.Parse("2006-01-02", resp.RequestDate)
	if err != nil {
		// Try parsing with timezone suffix if present
		requestDate, err = time.Parse("2006-01-02-07:00", resp.RequestDate)
		if err != nil {
			return nil, &ServiceError{
				Code:    ErrServiceError,
				Message: fmt.Sprintf("Failed to parse request date '%s': %v", resp.RequestDate, err),
			}
		}
	}

	// Create result
	result := &CheckVatResult{
		CountryCode: resp.CountryCode,
		VatNumber:   resp.VatNumber,
		RequestDate: requestDate,
		Valid:       resp.Valid,
		Name:        strings.TrimSpace(resp.Name),
		Address:     strings.TrimSpace(resp.Address),
	}

	return result, nil
}

// createSOAPRequest creates a SOAP envelope for VAT validation
func createSOAPRequest(countryCode, vatNumber string) *SOAPEnvelope {
	return &SOAPEnvelope{
		XmlnsSoapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		XmlnsUrn:     soapNamespace,
		Body: SOAPBody{
			CheckVat: &CheckVatRequest{
				CountryCode: countryCode,
				VatNumber:   vatNumber,
			},
		},
	}
}

// Ping tests connectivity to the VIES service
func (c *Client) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "HEAD", c.endpoint, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("service returned status: %s", resp.Status)
	}

	return nil
}
