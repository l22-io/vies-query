package vies

import (
	"encoding/xml"
	"time"
)

// CheckVatRequest represents the SOAP request for VAT validation
type CheckVatRequest struct {
	XMLName     xml.Name `xml:"urn:checkVat"`
	CountryCode string   `xml:"urn:countryCode"`
	VatNumber   string   `xml:"urn:vatNumber"`
}

// CheckVatResponse represents the SOAP response from VIES
type CheckVatResponse struct {
	XMLName     xml.Name  `xml:"checkVatResponse"`
	CountryCode string    `xml:"countryCode"`
	VatNumber   string    `xml:"vatNumber"`
	RequestDate time.Time `xml:"requestDate"`
	Valid       bool      `xml:"valid"`
	Name        string    `xml:"name"`
	Address     string    `xml:"address"`
}

// CheckVatResult represents the processed validation result
type CheckVatResult struct {
	CountryCode string    `json:"countryCode"`
	VatNumber   string    `json:"vatNumber"`
	RequestDate time.Time `json:"requestDate"`
	Valid       bool      `json:"valid"`
	Name        string    `json:"name,omitempty"`
	Address     string    `json:"address,omitempty"`
}

// SOAPEnvelope represents the SOAP envelope wrapper
type SOAPEnvelope struct {
	XMLName      xml.Name `xml:"soapenv:Envelope"`
	XmlnsSoapenv string   `xml:"xmlns:soapenv,attr"`
	XmlnsUrn     string   `xml:"xmlns:urn,attr"`
	Body         SOAPBody `xml:"soapenv:Body"`
}

// SOAPBody represents the SOAP body
type SOAPBody struct {
	CheckVat         *CheckVatRequest  `xml:"urn:checkVat,omitempty"`
	CheckVatResponse *CheckVatResponse `xml:"checkVatResponse,omitempty"`
	Fault            *SOAPFault        `xml:"soapenv:Fault,omitempty"`
}

// SOAPFault represents a SOAP fault response
type SOAPFault struct {
	XMLName xml.Name    `xml:"soapenv:Fault"`
	Code    string      `xml:"faultcode"`
	String  string      `xml:"faultstring"`
	Detail  FaultDetail `xml:"detail"`
}

// FaultDetail represents detailed fault information
type FaultDetail struct {
	XMLName xml.Name `xml:"detail"`
	Message string   `xml:",innerxml"`
}

// ValidationError represents VAT format validation errors
type ValidationError struct {
	Code      string
	Message   string
	VATNumber string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// ServiceError represents VIES service errors
type ServiceError struct {
	Code      string
	Message   string
	VATNumber string
}

func (e *ServiceError) Error() string {
	return e.Message
}

// Error codes for different types of failures
const (
	ErrInvalidFormat      = "INVALID_FORMAT"
	ErrUnsupportedCountry = "UNSUPPORTED_COUNTRY"
	ErrServiceError       = "SERVICE_ERROR"
	ErrNetworkTimeout     = "NETWORK_TIMEOUT"
	ErrServiceUnavailable = "SERVICE_UNAVAILABLE"
	ErrSOAPFault          = "SOAP_FAULT"
)

// ClientOptions for configuring the VIES client
type ClientOptions struct {
	Timeout   time.Duration
	UserAgent string
	Verbose   bool
	Endpoint  string
}

// ClientOption is a function type for configuring client options
type ClientOption func(*ClientOptions)

// WithTimeout sets the request timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(opts *ClientOptions) {
		opts.Timeout = timeout
	}
}

// WithUserAgent sets the User-Agent header
func WithUserAgent(userAgent string) ClientOption {
	return func(opts *ClientOptions) {
		opts.UserAgent = userAgent
	}
}

// WithVerbose enables verbose logging
func WithVerbose(verbose bool) ClientOption {
	return func(opts *ClientOptions) {
		opts.Verbose = verbose
	}
}

// WithEndpoint sets a custom endpoint (for testing)
func WithEndpoint(endpoint string) ClientOption {
	return func(opts *ClientOptions) {
		opts.Endpoint = endpoint
	}
}
