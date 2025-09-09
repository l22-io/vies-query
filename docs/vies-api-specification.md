# VIES API Complete Specification

This document contains the complete API specification for the European Commission's VAT Information Exchange System (VIES) based on the official WSDL.

## Service Information

- **Service Name**: checkVatService
- **WSDL URL**: https://ec.europa.eu/taxation_customs/vies/checkVatService.wsdl
- **Service URL**: http://ec.europa.eu/taxation_customs/vies/services/checkVatService (redirects to HTTPS)
- **Namespace**: `urn:ec.europa.eu:taxud:vies:services:checkVat:types`
- **SOAP Style**: Document/Literal
- **Transport**: HTTP

## Available Operations

The VIES service provides two operations:

### 1. checkVat (Basic VAT Validation)
- **Purpose**: Validate a VAT identification number
- **Input**: Country code and VAT number
- **Output**: Validation result with basic company information

### 2. checkVatApprox (Approximate VAT Validation with Matching)
- **Purpose**: Validate VAT number and match trader information
- **Input**: Country code, VAT number, and optional trader details
- **Output**: Validation result with detailed matching information

## Data Types and Schema

### checkVat Request
```xml
<xsd:element name="checkVat">
    <xsd:complexType>
        <xsd:sequence>
            <xsd:element name="countryCode" type="xsd:string"/>
            <xsd:element name="vatNumber" type="xsd:string"/>
        </xsd:sequence>
    </xsd:complexType>
</xsd:element>
```

### checkVat Response
```xml
<xsd:element name="checkVatResponse">
    <xsd:complexType>
        <xsd:sequence>
            <xsd:element name="countryCode" type="xsd:string"/>
            <xsd:element name="vatNumber" type="xsd:string"/>
            <xsd:element name="requestDate" type="xsd:date"/>
            <xsd:element name="valid" type="xsd:boolean"/>
            <xsd:element name="name" type="xsd:string" minOccurs="0" maxOccurs="1" nillable="true"/>
            <xsd:element name="address" type="xsd:string" minOccurs="0" maxOccurs="1" nillable="true"/>
        </xsd:sequence>
    </xsd:complexType>
</xsd:element>
```

**Important Notes:**
- `requestDate` is of type `xsd:date` (not `xsd:dateTime`)
- `name` and `address` are optional and nullable
- Date format should be YYYY-MM-DD (not ISO 8601 with time)

### checkVatApprox Request
```xml
<xsd:element name="checkVatApprox">
    <xsd:complexType>
        <xsd:sequence>
            <xsd:element name="countryCode" type="xsd:string"/>
            <xsd:element name="vatNumber" type="xsd:string"/>
            <xsd:element name="traderName" type="xsd:string" minOccurs="0" maxOccurs="1"/>
            <xsd:element name="traderCompanyType" type="companyTypeCode" minOccurs="0" maxOccurs="1"/>
            <xsd:element name="traderStreet" type="xsd:string" minOccurs="0" maxOccurs="1"/>
            <xsd:element name="traderPostcode" type="xsd:string" minOccurs="0" maxOccurs="1"/>
            <xsd:element name="traderCity" type="xsd:string" minOccurs="0" maxOccurs="1"/>
            <xsd:element name="requesterCountryCode" type="xsd:string" minOccurs="0" maxOccurs="1"/>
            <xsd:element name="requesterVatNumber" type="xsd:string" minOccurs="0" maxOccurs="1"/>
        </xsd:sequence>
    </xsd:complexType>
</xsd:element>
```

### checkVatApprox Response
```xml
<xsd:element name="checkVatApproxResponse">
    <xsd:complexType>
        <xsd:sequence>
            <xsd:element name="countryCode" type="xsd:string"/>
            <xsd:element name="vatNumber" type="xsd:string"/>
            <xsd:element name="requestDate" type="xsd:date"/>
            <xsd:element name="valid" type="xsd:boolean"/>
            <xsd:element name="traderName" type="xsd:string" minOccurs="0" maxOccurs="1" nillable="true"/>
            <xsd:element name="traderCompanyType" type="companyTypeCode" minOccurs="0" maxOccurs="1" nillable="true"/>
            <xsd:element name="traderAddress" type="xsd:string" minOccurs="0" maxOccurs="1"/>
            <xsd:element name="traderStreet" type="xsd:string" minOccurs="0" maxOccurs="1"/>
            <xsd:element name="traderPostcode" type="xsd:string" minOccurs="0" maxOccurs="1"/>
            <xsd:element name="traderCity" type="xsd:string" minOccurs="0" maxOccurs="1"/>
            <xsd:element name="traderNameMatch" type="matchCode" minOccurs="0" maxOccurs="1"/>
            <xsd:element name="traderCompanyTypeMatch" type="matchCode" minOccurs="0" maxOccurs="1"/>
            <xsd:element name="traderStreetMatch" type="matchCode" minOccurs="0" maxOccurs="1"/>
            <xsd:element name="traderPostcodeMatch" type="matchCode" minOccurs="0" maxOccurs="1"/>
            <xsd:element name="traderCityMatch" type="matchCode" minOccurs="0" maxOccurs="1"/>
            <xsd:element name="requestIdentifier" type="xsd:string"/>
        </xsd:sequence>
    </xsd:complexType>
</xsd:element>
```

## Custom Types

### companyTypeCode
```xml
<xsd:simpleType name="companyTypeCode">
    <xsd:restriction base="xsd:string">
        <xsd:pattern value="[A-Z]{2}\-[1-9][0-9]?"/>
    </xsd:restriction>
</xsd:simpleType>
```
Pattern: Two uppercase letters, hyphen, then 1-2 digits (e.g., "DE-1", "FR-12")

### matchCode
```xml
<xsd:simpleType name="matchCode">
    <xsd:restriction base="xsd:string">
        <xsd:enumeration value="1"> <!-- VALID -->
        <xsd:enumeration value="2"> <!-- INVALID -->
        <xsd:enumeration value="3"> <!-- NOT_PROCESSED -->
    </xsd:restriction>
</xsd:simpleType>
```

## Input Validation Rules

According to the WSDL documentation:

### Country Code
- **Pattern**: `[A-Z]{2}`
- **Description**: Must be exactly 2 uppercase letters
- **Examples**: `DE`, `FR`, `NL`, `IT`

### VAT Number
- **Pattern**: `[0-9A-Za-z\+\*\.]{2,12}`
- **Description**: 2-12 characters, alphanumeric plus `+`, `*`, `.`
- **Examples**: `123456789`, `123456789B01`, `A12345678`

## Error Handling

### SOAP Fault Codes

The service returns SOAP faults with specific fault strings:

1. **INVALID_INPUT**
   - Cause: Invalid country code or empty VAT number
   - Action: Validate input format before sending request

2. **GLOBAL_MAX_CONCURRENT_REQ**
   - Cause: Maximum concurrent requests reached globally
   - Action: Implement exponential backoff retry

3. **MS_MAX_CONCURRENT_REQ**
   - Cause: Maximum concurrent requests for specific Member State
   - Action: Implement exponential backoff retry for that country

4. **SERVICE_UNAVAILABLE**
   - Cause: Network or application level error
   - Action: Retry after delay

5. **MS_UNAVAILABLE**
   - Cause: Member State application not responding
   - Action: Check service status, retry later

6. **TIMEOUT**
   - Cause: Request timeout exceeded
   - Action: Retry with potentially longer timeout

## SOAP Message Structure

### Request Example
```xml
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" 
                  xmlns:urn="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
   <soapenv:Header/>
   <soapenv:Body>
      <urn:checkVat>
         <urn:countryCode>DE</urn:countryCode>
         <urn:vatNumber>123456789</urn:vatNumber>
      </urn:checkVat>
   </soapenv:Body>
</soapenv:Envelope>
```

### Response Example (Success)
```xml
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
   <soapenv:Body>
      <ns2:checkVatResponse xmlns:ns2="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
         <ns2:countryCode>DE</ns2:countryCode>
         <ns2:vatNumber>123456789</ns2:vatNumber>
         <ns2:requestDate>2025-01-09</ns2:requestDate>
         <ns2:valid>true</ns2:valid>
         <ns2:name>Example GmbH</ns2:name>
         <ns2:address>Musterstraße 1, 12345 Berlin</ns2:address>
      </ns2:checkVatResponse>
   </soapenv:Body>
</soapenv:Envelope>
```

### Response Example (Fault)
```xml
<env:Envelope xmlns:env="http://schemas.xmlsoap.org/soap/envelope/">
   <env:Body>
      <env:Fault>
         <faultcode>env:Server</faultcode>
         <faultstring>MS_UNAVAILABLE</faultstring>
      </env:Fault>
   </env:Body>
</env:Envelope>
```

## Legal and Usage Information

From the WSDL documentation:

- **Purpose**: Confirm validity of VAT identification numbers for intra-Community supplies
- **Prohibited Uses**:
  - Commercial retransmission
  - Data extraction not conforming to the site's objective
  - Copying or reproduction of contents
- **Limitations**:
  - No liability for information accuracy
  - Data comes from Member State databases
  - Does not provide legal/professional advice
  - Does not grant VAT exemption rights
- **Data Protection**: Processing falls under EU Data Protection Notice

## Implementation Notes

1. **Date Handling**: The `requestDate` field uses `xsd:date` format (YYYY-MM-DD), not datetime
2. **Nullable Fields**: Company name and address can be null/empty even for valid VAT numbers
3. **Rate Limiting**: Implement proper retry logic for concurrent request errors
4. **Namespace Handling**: Ensure proper XML namespace declarations in requests
5. **Error Handling**: Parse and handle all documented SOAP fault types
6. **Input Validation**: Pre-validate input patterns before sending requests
