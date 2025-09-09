package output

import (
	"encoding/json"

	"l22.io/viesquery/internal/vies"
)

// JSONFormatter formats output as JSON
type JSONFormatter struct{}

// NewJSONFormatter creates a new JSON formatter
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// Format formats a validation result as JSON
func (f *JSONFormatter) Format(result *vies.CheckVatResult) (string, error) {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}

// ErrorResponse represents an error in JSON format
type ErrorResponse struct {
	Error     bool   `json:"error"`
	Message   string `json:"message"`
	Code      string `json:"code,omitempty"`
	VATNumber string `json:"vatNumber,omitempty"`
}

// FormatError formats an error as JSON
func (f *JSONFormatter) FormatError(err error) (string, error) {
	var errorResponse ErrorResponse

	errorResponse.Error = true
	errorResponse.Message = err.Error()

	switch e := err.(type) {
	case *vies.ValidationError:
		errorResponse.Code = e.Code
		errorResponse.VATNumber = e.VATNumber
	case *vies.ServiceError:
		errorResponse.Code = e.Code
		errorResponse.VATNumber = e.VATNumber
	}

	data, err := json.MarshalIndent(errorResponse, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data) + "\n", nil
}
