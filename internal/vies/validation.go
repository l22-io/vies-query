package vies

import (
	"fmt"
	"regexp"
	"strings"
)

// CountryValidator contains validation rules for a specific EU country
type CountryValidator struct {
	Code        string
	Name        string
	Pattern     *regexp.Regexp
	MinLength   int
	MaxLength   int
	Description string
}

// EU member state VAT validation patterns
var countryValidators = map[string]CountryValidator{
	"AT": {
		Code:        "AT",
		Name:        "Austria",
		Pattern:     regexp.MustCompile(`^ATU\d{8}$`),
		MinLength:   11,
		MaxLength:   11,
		Description: "ATU + 8 digits",
	},
	"BE": {
		Code:        "BE",
		Name:        "Belgium",
		Pattern:     regexp.MustCompile(`^BE[01]\d{9}$`),
		MinLength:   12,
		MaxLength:   12,
		Description: "BE0 or BE1 + 9 digits",
	},
	"BG": {
		Code:        "BG",
		Name:        "Bulgaria",
		Pattern:     regexp.MustCompile(`^BG\d{9,10}$`),
		MinLength:   11,
		MaxLength:   12,
		Description: "BG + 9 or 10 digits",
	},
	"HR": {
		Code:        "HR",
		Name:        "Croatia",
		Pattern:     regexp.MustCompile(`^HR\d{11}$`),
		MinLength:   13,
		MaxLength:   13,
		Description: "HR + 11 digits",
	},
	"CY": {
		Code:        "CY",
		Name:        "Cyprus",
		Pattern:     regexp.MustCompile(`^CY\d{8}[A-Z]$`),
		MinLength:   11,
		MaxLength:   11,
		Description: "CY + 8 digits + 1 letter",
	},
	"CZ": {
		Code:        "CZ",
		Name:        "Czech Republic",
		Pattern:     regexp.MustCompile(`^CZ\d{8,10}$`),
		MinLength:   10,
		MaxLength:   12,
		Description: "CZ + 8, 9, or 10 digits",
	},
	"DK": {
		Code:        "DK",
		Name:        "Denmark",
		Pattern:     regexp.MustCompile(`^DK\d{8}$`),
		MinLength:   10,
		MaxLength:   10,
		Description: "DK + 8 digits",
	},
	"EE": {
		Code:        "EE",
		Name:        "Estonia",
		Pattern:     regexp.MustCompile(`^EE\d{9}$`),
		MinLength:   11,
		MaxLength:   11,
		Description: "EE + 9 digits",
	},
	"FI": {
		Code:        "FI",
		Name:        "Finland",
		Pattern:     regexp.MustCompile(`^FI\d{8}$`),
		MinLength:   10,
		MaxLength:   10,
		Description: "FI + 8 digits",
	},
	"FR": {
		Code:        "FR",
		Name:        "France",
		Pattern:     regexp.MustCompile(`^FR[A-HJ-NP-Z0-9]{2}\d{9}$`),
		MinLength:   13,
		MaxLength:   13,
		Description: "FR + 2 characters + 9 digits",
	},
	"DE": {
		Code:        "DE",
		Name:        "Germany",
		Pattern:     regexp.MustCompile(`^DE\d{9}$`),
		MinLength:   11,
		MaxLength:   11,
		Description: "DE + 9 digits",
	},
	"EL": {
		Code:        "EL",
		Name:        "Greece",
		Pattern:     regexp.MustCompile(`^EL\d{9}$`),
		MinLength:   11,
		MaxLength:   11,
		Description: "EL + 9 digits",
	},
	"GR": { // Alternative code for Greece
		Code:        "GR",
		Name:        "Greece",
		Pattern:     regexp.MustCompile(`^GR\d{9}$`),
		MinLength:   11,
		MaxLength:   11,
		Description: "GR + 9 digits (alternative for EL)",
	},
	"HU": {
		Code:        "HU",
		Name:        "Hungary",
		Pattern:     regexp.MustCompile(`^HU\d{8}$`),
		MinLength:   10,
		MaxLength:   10,
		Description: "HU + 8 digits",
	},
	"IE": {
		Code:        "IE",
		Name:        "Ireland",
		Pattern:     regexp.MustCompile(`^IE[A-Z0-9]{8}$`),
		MinLength:   10,
		MaxLength:   10,
		Description: "IE + 8 alphanumeric characters",
	},
	"IT": {
		Code:        "IT",
		Name:        "Italy",
		Pattern:     regexp.MustCompile(`^IT\d{11}$`),
		MinLength:   13,
		MaxLength:   13,
		Description: "IT + 11 digits",
	},
	"LV": {
		Code:        "LV",
		Name:        "Latvia",
		Pattern:     regexp.MustCompile(`^LV\d{11}$`),
		MinLength:   13,
		MaxLength:   13,
		Description: "LV + 11 digits",
	},
	"LT": {
		Code:        "LT",
		Name:        "Lithuania",
		Pattern:     regexp.MustCompile(`^LT(\d{9}|\d{12})$`),
		MinLength:   11,
		MaxLength:   14,
		Description: "LT + 9 or 12 digits",
	},
	"LU": {
		Code:        "LU",
		Name:        "Luxembourg",
		Pattern:     regexp.MustCompile(`^LU\d{8}$`),
		MinLength:   10,
		MaxLength:   10,
		Description: "LU + 8 digits",
	},
	"MT": {
		Code:        "MT",
		Name:        "Malta",
		Pattern:     regexp.MustCompile(`^MT\d{8}$`),
		MinLength:   10,
		MaxLength:   10,
		Description: "MT + 8 digits",
	},
	"NL": {
		Code:        "NL",
		Name:        "Netherlands",
		Pattern:     regexp.MustCompile(`^NL\d{9}B\d{2}$`),
		MinLength:   14,
		MaxLength:   14,
		Description: "NL + 9 digits + B + 2 digits",
	},
	"PL": {
		Code:        "PL",
		Name:        "Poland",
		Pattern:     regexp.MustCompile(`^PL\d{10}$`),
		MinLength:   12,
		MaxLength:   12,
		Description: "PL + 10 digits",
	},
	"PT": {
		Code:        "PT",
		Name:        "Portugal",
		Pattern:     regexp.MustCompile(`^PT\d{9}$`),
		MinLength:   11,
		MaxLength:   11,
		Description: "PT + 9 digits",
	},
	"RO": {
		Code:        "RO",
		Name:        "Romania",
		Pattern:     regexp.MustCompile(`^RO\d{2,10}$`),
		MinLength:   4,
		MaxLength:   12,
		Description: "RO + 2 to 10 digits",
	},
	"SK": {
		Code:        "SK",
		Name:        "Slovakia",
		Pattern:     regexp.MustCompile(`^SK\d{10}$`),
		MinLength:   12,
		MaxLength:   12,
		Description: "SK + 10 digits",
	},
	"SI": {
		Code:        "SI",
		Name:        "Slovenia",
		Pattern:     regexp.MustCompile(`^SI\d{8}$`),
		MinLength:   10,
		MaxLength:   10,
		Description: "SI + 8 digits",
	},
	"ES": {
		Code:        "ES",
		Name:        "Spain",
		Pattern:     regexp.MustCompile(`^ES[A-Z0-9]\d{7}[A-Z0-9]$`),
		MinLength:   11,
		MaxLength:   11,
		Description: "ES + character + 7 digits + character",
	},
	"SE": {
		Code:        "SE",
		Name:        "Sweden",
		Pattern:     regexp.MustCompile(`^SE\d{12}$`),
		MinLength:   14,
		MaxLength:   14,
		Description: "SE + 12 digits",
	},
}

// ValidateFormat validates VAT number format according to EU country rules
func ValidateFormat(vatNumber string) error {
	// Remove spaces and convert to uppercase
	vatNumber = strings.ToUpper(strings.ReplaceAll(vatNumber, " ", ""))

	if len(vatNumber) < 3 {
		return &ValidationError{
			Code:      ErrInvalidFormat,
			Message:   "VAT number too short (minimum 3 characters)",
			VATNumber: vatNumber,
		}
	}

	// Extract country code (first 2 characters)
	countryCode := vatNumber[:2]
	
	// Special case: Some systems use GR instead of EL for Greece
	if countryCode == "GR" {
		countryCode = "EL"
		vatNumber = "EL" + vatNumber[2:]
	}

	validator, exists := countryValidators[countryCode]
	if !exists {
		return &ValidationError{
			Code:      ErrUnsupportedCountry,
			Message:   fmt.Sprintf("Unsupported country code: %s", countryCode),
			VATNumber: vatNumber,
		}
	}

	// Check length
	if len(vatNumber) < validator.MinLength || len(vatNumber) > validator.MaxLength {
		return &ValidationError{
			Code:    ErrInvalidFormat,
			Message: fmt.Sprintf("Invalid length for %s VAT number. Expected: %s", validator.Name, validator.Description),
			VATNumber: vatNumber,
		}
	}

	// Check pattern
	if !validator.Pattern.MatchString(vatNumber) {
		return &ValidationError{
			Code:    ErrInvalidFormat,
			Message: fmt.Sprintf("Invalid format for %s VAT number. Expected: %s", validator.Name, validator.Description),
			VATNumber: vatNumber,
		}
	}

	return nil
}

// ParseVATNumber extracts country code and VAT number from a full VAT number
func ParseVATNumber(vatNumber string) (string, string, error) {
	// Clean and validate format
	if err := ValidateFormat(vatNumber); err != nil {
		return "", "", err
	}

	// Remove spaces and convert to uppercase
	vatNumber = strings.ToUpper(strings.ReplaceAll(vatNumber, " ", ""))

	// Special case: Convert GR to EL for Greece
	if strings.HasPrefix(vatNumber, "GR") {
		vatNumber = "EL" + vatNumber[2:]
	}

	countryCode := vatNumber[:2]
	number := vatNumber[2:]

	// Some countries have prefix letters (like ATU for Austria)
	// Remove them for the API call
	switch countryCode {
	case "AT":
		if strings.HasPrefix(number, "U") {
			number = number[1:] // Remove the 'U' prefix
		}
	}

	return countryCode, number, nil
}

// GetSupportedCountries returns a list of all supported country codes
func GetSupportedCountries() []string {
	countries := make([]string, 0, len(countryValidators))
	for code := range countryValidators {
		if code != "GR" { // Skip GR as it's an alias for EL
			countries = append(countries, code)
		}
	}
	return countries
}

// GetCountryInfo returns information about a specific country's VAT format
func GetCountryInfo(countryCode string) (*CountryValidator, error) {
	validator, exists := countryValidators[countryCode]
	if !exists {
		return nil, fmt.Errorf("unsupported country code: %s", countryCode)
	}
	return &validator, nil
}
