package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Build time variable
var versionString = "sheets2json"

func main() {

	var (
		credentialFile = pflag.StringP("credential", "c", "", "Path to credential JSON file")
		outputFile     = pflag.StringP("output", "o", "", "Output file (default: stdout)")
		showVersion    = pflag.Bool("version", false, "Show version information")
	)
	pflag.Parse()

	// Handle version flag
	if *showVersion {
		fmt.Println(versionString)
		os.Exit(0)
	}

	if pflag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Usage: sheets2json [-c CREDENTIAL_FILE] [-o OUTPUT_FILE] [--version] SPREADSHEET_ID [WORKSHEET] [RANGE]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -c, --credential string   Path to credential JSON file\n")
		fmt.Fprintf(os.Stderr, "  -o, --output string       Output file (default: stdout)\n")
		fmt.Fprintf(os.Stderr, "      --version             Show version information\n")
		os.Exit(1)
	}

	spreadsheetID := pflag.Arg(0)
	worksheetArg := ""
	rangeArg := ""

	if pflag.NArg() >= 2 {
		worksheetArg = pflag.Arg(1)
	}
	if pflag.NArg() >= 3 {
		rangeArg = pflag.Arg(2)
	}

	// Get credentials
	var credData []byte
	var err error

	// Priority: 1. Command line file, 2. Environment variable
	if *credentialFile != "" {
		credData, err = os.ReadFile(*credentialFile)
		if err != nil {
			log.Fatalf("Error reading credential file: %v", err)
		}
	} else if envCredFile := os.Getenv("GOOGLE_SHEETS_CREDENTIAL"); envCredFile != "" {
		credData, err = os.ReadFile(envCredFile)
		if err != nil {
			log.Fatalf("Error reading credential file from GOOGLE_SHEETS_CREDENTIAL: %v", err)
		}
	} else {
		log.Fatal("No credentials provided. Use -c or GOOGLE_SHEETS_CREDENTIAL environment variable")
	}

	// Create Google Sheets service
	ctx := context.Background()
	config, err := google.JWTConfigFromJSON(credData, sheets.SpreadsheetsReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse credential JSON: %v", err)
	}

	client := config.Client(ctx)
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to create Sheets service: %v", err)
	}

	// Determine range from positional arguments
	rangeToRead := "A:ZZ" // Default: read all columns

	if worksheetArg != "" {
		if rangeArg != "" {
			rangeToRead = worksheetArg + "!" + rangeArg
		} else {
			rangeToRead = worksheetArg + "!A:ZZ"
		}
	}

	// Get data
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, rangeToRead).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	if len(resp.Values) == 0 {
		log.Fatal("No data found")
	}

	// Convert to JSON array of objects
	result := convertToJSON(resp.Values)

	// Output
	var output io.Writer = os.Stdout
	if *outputFile != "" {
		file, err := os.Create(*outputFile)
		if err != nil {
			log.Fatalf("Error creating output file: %v", err)
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Printf("Error closing output file: %v", err)
			}
		}()
		output = file
	}

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
	}

	if *outputFile != "" {
		log.Printf("Data saved to %s", *outputFile)
	}
}

// OrderedRow represents a row with ordered keys
type OrderedRow struct {
	Keys   []string
	Values map[string]interface{}
}

// MarshalJSON implements json.Marshaler to maintain key order
func (r OrderedRow) MarshalJSON() ([]byte, error) {
	// Build JSON manually to preserve order
	items := make([]string, 0, len(r.Keys))
	for _, key := range r.Keys {
		value := r.Values[key]
		keyJSON, _ := json.Marshal(key)
		valueJSON, _ := json.Marshal(value)
		items = append(items, fmt.Sprintf("%s:%s", keyJSON, valueJSON))
	}
	return []byte("{" + strings.Join(items, ",") + "}"), nil
}

func convertToJSON(values [][]interface{}) []OrderedRow {
	if len(values) == 0 {
		return []OrderedRow{}
	}

	// First row as headers
	headers := make([]string, len(values[0]))
	for i, cell := range values[0] {
		headers[i] = fmt.Sprintf("%v", cell)
	}

	// Convert remaining rows
	result := make([]OrderedRow, 0, len(values)-1)
	for _, row := range values[1:] {
		obj := OrderedRow{
			Keys:   headers,
			Values: make(map[string]interface{}),
		}
		for i, header := range headers {
			if i < len(row) {
				// Keep original type if possible
				obj.Values[header] = row[i]
			} else {
				obj.Values[header] = ""
			}
		}
		result = append(result, obj)
	}

	return result
}
