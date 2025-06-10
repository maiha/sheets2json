# sheets2json

![Go](https://img.shields.io/badge/go-1.21-blue.svg)
![License](https://img.shields.io/badge/license-MIT-green.svg)

A CLI tool to fetch Google Sheets data as JSON

```bash
./sheets2json SPREADSHEET_ID                    # Get all sheet data
./sheets2json SPREADSHEET_ID Sheet2             # Specify a worksheet (Sheet2)
./sheets2json SPREADSHEET_ID Sheet2 A1:C10      # Specify worksheet range
```

## Usage

```
sheets2json [OPTIONS] SPREADSHEET_ID [WORKSHEET] [RANGE]
```

### Options
- `-c string` - Path to credential JSON file
- `-o string` - Output file (default: stdout)
- `--version` - Show version information

### Authentication

Enable Google Sheets API in Google Cloud Console and obtain service account credential JSON file.

```bash
# Use either method
./sheets2json -c credentials.json SPREADSHEET_ID
export GOOGLE_SHEETS_CREDENTIAL=credentials.json
```

### Other Examples

```bash
# Output to file
./sheets2json -o output.json SPREADSHEET_ID Sheet2
```

## Build

Requires Docker to build. Outputs `sheets2json` executable binary (Linux x86_64, statically linked).

```bash
make build      # Build
make clean      # Clean up
```
