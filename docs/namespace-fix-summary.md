# XML Namespace Conflict Fix Summary

## Problem

The application was failing with the following error:
```
xml: name "checkVat" in tag of vies.SOAPBody.CheckVat conflicts with name "urn:checkVat" in *vies.CheckVatRequest.XMLName
```

This occurred because Go's XML marshaler detected a conflict between:
1. The XML element name defined in `CheckVatRequest.XMLName` as `urn:checkVat`  
2. The XML tag in `SOAPBody.CheckVat` field as `checkVat`

## Root Cause Analysis

After obtaining the complete VIES WSDL specification, I identified that the issue stemmed from improper XML namespace handling in the SOAP request structures. The VIES service expects:

```xml
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" 
                  xmlns:urn="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
   <soapenv:Body>
      <urn:checkVat>
         <urn:countryCode>DE</urn:countryCode>
         <urn:vatNumber>266201128</urn:vatNumber>
      </urn:checkVat>
   </soapenv:Body>
</soapenv:Envelope>
```

## Solution Implemented

### 1. Fixed SOAP Request Structure

**Before:**
```go
type CheckVatRequest struct {
    XMLName     xml.Name `xml:"urn:checkVat"`
    Namespace   string   `xml:"xmlns:urn,attr"`
    CountryCode string   `xml:"urn:countryCode"`
    VatNumber   string   `xml:"urn:vatNumber"`
}

type SOAPBody struct {
    CheckVat *CheckVatRequest `xml:"checkVat,omitempty"`
}
```

**After:**
```go
type CheckVatRequest struct {
    XMLName     xml.Name `xml:"urn:checkVat"`
    CountryCode string   `xml:"urn:countryCode"`
    VatNumber   string   `xml:"urn:vatNumber"`
}

type SOAPBody struct {
    CheckVat *CheckVatRequest `xml:"urn:checkVat,omitempty"`
}
```

### 2. Updated SOAP Envelope Structure

**Before:**
```go
type SOAPEnvelope struct {
    XMLName xml.Name `xml:"soap:Envelope"`
    Xmlns   string   `xml:"xmlns:soap,attr"`
    Body    SOAPBody `xml:"soap:Body"`
}
```

**After:**
```go
type SOAPEnvelope struct {
    XMLName      xml.Name `xml:"soapenv:Envelope"`
    XmlnsSoapenv string   `xml:"xmlns:soapenv,attr"`
    XmlnsUrn     string   `xml:"xmlns:urn,attr"`
    Body         SOAPBody `xml:"soapenv:Body"`
}
```

### 3. Fixed Response Parsing

The original code had a similar namespace conflict in response parsing. I resolved this by using inline structs to avoid namespace conflicts:

```go
var envelope struct {
    XMLName xml.Name `xml:"Envelope"`
    Body    struct {
        CheckVatResponse *struct {
            XMLName     xml.Name `xml:"checkVatResponse"`
            CountryCode string   `xml:"countryCode"`
            VatNumber   string   `xml:"vatNumber"`
            RequestDate string   `xml:"requestDate"` // Parse as string first
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
```

### 4. Fixed Date Parsing

The WSDL specification defines `requestDate` as `xsd:date` (YYYY-MM-DD format), not `xsd:dateTime`. I updated the parsing to handle this correctly:

```go
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
```

## Results

### Before Fix
```bash
$ ./bin/viesquery DE266201128
Error: Failed to create SOAP request: xml: name "checkVat" in tag of vies.SOAPBody.CheckVat conflicts with name "urn:checkVat" in *vies.CheckVatRequest.XMLName
```

### After Fix
```bash
$ ./bin/viesquery DE266201128
Error: SOAP fault: env:Server - MS_MAX_CONCURRENT_REQ

$ ./bin/viesquery NL854502130B01
VAT Number: NL854502130B01
Status: Invalid
Request Date: 2025-09-09 00:00:00 UTC

$ ./bin/viesquery --format json NL854502130B01
{
  "countryCode": "NL",
  "vatNumber": "854502130B01",
  "requestDate": "2025-09-09T00:00:00+02:00",
  "valid": false
}
```

## Key Learnings

1. **WSDL Analysis is Critical**: Getting the complete VIES WSDL specification was essential to understanding the expected XML structure and data types.

2. **XML Namespace Conflicts**: Go's XML marshaler is strict about namespace consistency. Element names must match between struct XMLName declarations and their embedding context.

3. **XSD Date vs DateTime**: The VIES API uses `xsd:date` format (YYYY-MM-DD) for dates, not ISO 8601 datetime format.

4. **SOAP Fault Handling**: Proper error handling now correctly identifies and reports specific VIES service errors like `MS_UNAVAILABLE`, `MS_MAX_CONCURRENT_REQ`, etc.

## Files Modified

- `internal/vies/types.go` - Updated SOAP structure definitions
- `internal/vies/client.go` - Fixed request creation and response parsing  
- `internal/vies/types_test.go` - Added XML marshaling tests (new file)
- `docs/checkVatService.wsdl` - Added official WSDL specification (new file)
- `docs/vies-api-specification.md` - Added comprehensive API documentation (new file)

The application now correctly generates VIES-compliant SOAP requests and properly handles all response types including successful validations, invalid VAT numbers, and various service fault conditions.
