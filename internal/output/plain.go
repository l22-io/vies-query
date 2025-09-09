package output

import (
	"fmt"
	"strings"

	"l22.io/viesquery/internal/vies"
)

// PlainFormatter formats output as plain text
type PlainFormatter struct{}

// NewPlainFormatter creates a new plain text formatter
func NewPlainFormatter() *PlainFormatter {
	return &PlainFormatter{}
}

// Format formats a validation result as plain text
func (f *PlainFormatter) Format(result *vies.CheckVatResult) (string, error) {
	var b strings.Builder

	// VAT Number
	fmt.Fprintf(&b, "VAT Number: %s%s\n", result.CountryCode, result.VatNumber)

	// Status
	status := "Invalid"
	if result.Valid {
		status = "Valid"
	}
	fmt.Fprintf(&b, "Status: %s\n", status)

	// Company information (only if valid and available)
	if result.Valid {
		if result.Name != "" {
			fmt.Fprintf(&b, "Company: %s\n", result.Name)
		}
		if result.Address != "" {
			fmt.Fprintf(&b, "Address: %s\n", result.Address)
		}
	}

	// Request date
	fmt.Fprintf(&b, "Request Date: %s\n", result.RequestDate.Format("2006-01-02 15:04:05 UTC"))

	return b.String(), nil
}

// FormatError formats an error as plain text
func (f *PlainFormatter) FormatError(err error) (string, error) {
	var b strings.Builder

	switch e := err.(type) {
	case *vies.ValidationError:
		fmt.Fprintf(&b, "Error: %s\n", e.Message)
		if e.VATNumber != "" {
			fmt.Fprintf(&b, "VAT Number: %s\n", e.VATNumber)
		}
		
		// Add format hint for validation errors
		if e.Code == vies.ErrInvalidFormat {
			// Try to get country info for format hint
			if len(e.VATNumber) >= 2 {
				countryCode := e.VATNumber[:2]
				if countryInfo, err := vies.GetCountryInfo(countryCode); err == nil {
					fmt.Fprintf(&b, "Expected Format: %s\n", countryInfo.Description)
				}
			}
		}

	case *vies.ServiceError:
		fmt.Fprintf(&b, "Error: %s\n", e.Message)
		if e.VATNumber != "" {
			fmt.Fprintf(&b, "VAT Number: %s\n", e.VATNumber)
		}
		
		// Add specific suggestions for service errors
		switch e.Code {
		case vies.ErrNetworkTimeout:
			fmt.Fprintf(&b, "Try increasing timeout with --timeout flag\n")
		case vies.ErrServiceUnavailable:
			fmt.Fprintf(&b, "Please retry later or check VIES service status\n")
		}

	default:
		fmt.Fprintf(&b, "Error: %s\n", err.Error())
	}

	return b.String(), nil
}
