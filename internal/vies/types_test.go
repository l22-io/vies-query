package vies

import (
	"encoding/xml"
	"strings"
	"testing"
)

func TestSOAPRequestMarshaling(t *testing.T) {
	// Expected XML structure based on VIES documentation
	_ = `<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:urn="urn:ec.europa.eu:taxud:vies:services:checkVat:types">
   <soapenv:Body>
      <urn:checkVat>
         <urn:countryCode>DE</urn:countryCode>
         <urn:vatNumber>266201128</urn:vatNumber>
      </urn:checkVat>
   </soapenv:Body>
</soapenv:Envelope>`

	// Create SOAP request
	soapRequest := createSOAPRequest("DE", "266201128")
	
	// Marshal to XML
	requestBody, err := xml.Marshal(soapRequest)
	if err != nil {
		t.Fatalf("Failed to marshal SOAP request: %v", err)
	}

	actualXML := string(requestBody)
	t.Logf("Actual XML:\n%s", actualXML)
	
	// Basic validation - check for key elements
	if !strings.Contains(actualXML, "checkVat") {
		t.Error("Missing checkVat element")
	}
	if !strings.Contains(actualXML, "countryCode") {
		t.Error("Missing countryCode element")
	}
	if !strings.Contains(actualXML, "vatNumber") {
		t.Error("Missing vatNumber element")
	}
	if !strings.Contains(actualXML, "DE") {
		t.Error("Missing country code value")
	}
	if !strings.Contains(actualXML, "266201128") {
		t.Error("Missing VAT number value")
	}
}

func TestSimpleStructMarshaling(t *testing.T) {
	// Test a simple structure to understand the namespace issue
	type TestRequest struct {
		XMLName xml.Name `xml:"urn:ec.europa.eu:taxud:vies:services:checkVat:types checkVat"`
		Country string   `xml:"urn:ec.europa.eu:taxud:vies:services:checkVat:types countryCode"`
		VAT     string   `xml:"urn:ec.europa.eu:taxud:vies:services:checkVat:types vatNumber"`
	}

	req := TestRequest{
		Country: "DE",
		VAT:     "266201128",
	}

	data, err := xml.Marshal(req)
	if err != nil {
		t.Logf("Simple struct marshaling failed: %v", err)
	} else {
		t.Logf("Simple struct XML: %s", string(data))
	}
}

func TestNamespaceApproach(t *testing.T) {
	// Test different namespace approaches
	type TestEnvelope struct {
		XMLName      xml.Name `xml:"soapenv:Envelope"`
		XmlnsSoapenv string   `xml:"xmlns:soapenv,attr"`
		XmlnsUrn     string   `xml:"xmlns:urn,attr"`
		Body         struct {
			CheckVat struct {
				XMLName     xml.Name `xml:"urn:checkVat"`
				CountryCode string   `xml:"urn:countryCode"`
				VatNumber   string   `xml:"urn:vatNumber"`
			} `xml:",omitempty"`
		} `xml:"soapenv:Body"`
	}

	envelope := TestEnvelope{
		XmlnsSoapenv: "http://schemas.xmlsoap.org/soap/envelope/",
		XmlnsUrn:     "urn:ec.europa.eu:taxud:vies:services:checkVat:types",
	}
	envelope.Body.CheckVat.CountryCode = "DE"
	envelope.Body.CheckVat.VatNumber = "266201128"

	data, err := xml.Marshal(envelope)
	if err != nil {
		t.Logf("Namespace approach failed: %v", err)
	} else {
		t.Logf("Namespace approach XML: %s", string(data))
	}
}
