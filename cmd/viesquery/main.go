package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"l22.io/viesquery/internal/output"
	"l22.io/viesquery/internal/vies"
)

var (
	// Version is set during build time
	Version = "dev"
)

func main() {
	var (
		format    = flag.String("format", getEnvString("VIESQUERY_FORMAT", "plain"), "Output format (plain, json)")
		timeout   = flag.Int("timeout", getEnvInt("VIESQUERY_TIMEOUT", 30), "Request timeout in seconds")
		verbose   = flag.Bool("verbose", getEnvBool("VIESQUERY_VERBOSE", false), "Enable verbose logging")
		version   = flag.Bool("version", false, "Display version information")
		help      = flag.Bool("help", false, "Display help information")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "VIES Query - EU VAT Number Validation Tool\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] VAT_NUMBER\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Validate EU VAT numbers using the VIES API\n\n")
		fmt.Fprintf(os.Stderr, "Arguments:\n")
		fmt.Fprintf(os.Stderr, "  VAT_NUMBER    EU VAT number to validate (e.g., DE123456789)\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nSupported Countries:\n")
		fmt.Fprintf(os.Stderr, "  AT, BE, BG, HR, CY, CZ, DK, EE, FI, FR, DE, EL,\n")
		fmt.Fprintf(os.Stderr, "  HU, IE, IT, LV, LT, LU, MT, NL, PL, PT, RO, SK, SI, ES, SE\n\n")
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  %s DE123456789\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --format json AT12345678\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --timeout 60 --verbose IT12345678901\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nEnvironment Variables:\n")
		fmt.Fprintf(os.Stderr, "  VIESQUERY_FORMAT   Default output format (plain, json)\n")
		fmt.Fprintf(os.Stderr, "  VIESQUERY_TIMEOUT  Default timeout in seconds\n")
		fmt.Fprintf(os.Stderr, "  VIESQUERY_VERBOSE  Enable verbose mode (true, false)\n")
	}

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *version {
		fmt.Printf("viesquery version %s\n", Version)
		fmt.Printf("https://github.com/l22-io/vies-query\n")
		os.Exit(0)
	}

	if flag.NArg() != 1 {
		fmt.Fprintf(os.Stderr, "Error: VAT number required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	vatNumber := flag.Arg(0)

	// Validate output format
	if *format != "plain" && *format != "json" {
		fmt.Fprintf(os.Stderr, "Error: Invalid format '%s'. Supported formats: plain, json\n", *format)
		os.Exit(1)
	}

	// Validate timeout
	if *timeout < 1 {
		fmt.Fprintf(os.Stderr, "Error: Invalid timeout '%d'. Must be greater than 0\n", *timeout)
		os.Exit(1)
	}

	// Create VIES client
	client := vies.NewClient(
		vies.WithTimeout(time.Duration(*timeout)*time.Second),
		vies.WithVerbose(*verbose),
	)

	// Validate VAT number
	ctx := context.Background()
	result, err := client.CheckVAT(ctx, vatNumber)
	if err != nil {
		handleError(err, *format)
		return
	}

	// Display result
	displayResult(result, *format)
}

func handleError(err error, format string) {
	formatter := output.NewManager()
	f, fErr := formatter.GetFormatter(format)
	if fErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	output, formatErr := f.FormatError(err)
	if formatErr != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	fmt.Print(output)

	// Set appropriate exit code based on error type
	switch e := err.(type) {
	case *vies.ValidationError:
		os.Exit(3) // Invalid VAT format
	case *vies.ServiceError:
		if e.Code == vies.ErrServiceUnavailable {
			os.Exit(4) // Service unavailable
		}
		os.Exit(2) // Network/API error
	default:
		os.Exit(2) // General error
	}
}

func displayResult(result *vies.CheckVatResult, format string) {
	formatter := output.NewManager()
	f, err := formatter.GetFormatter(format)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(2)
	}

	output, err := f.Format(result)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
		os.Exit(2)
	}

	fmt.Print(output)
}

// getEnvString returns environment variable value or default
func getEnvString(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt returns environment variable as integer or default
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}

// getEnvBool returns environment variable as boolean or default
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
