package utils_test

import (
	"fmt"
	"os"
	"testing"

	"time"

	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
	"tmossDev.github.com/eco-system/shared-components/backend/package/utils"
)

func TestGetenv_withRetrievingEnvVaraiableWithDefault_shouldReturnDefaultIfNotPresent(t *testing.T) {
	_ = os.Setenv("test01", "test01")
	result := utils.Getenv("test00", "test00")
	if result != "test00" {
		t.Errorf("Getenv test failed, expected[test00], got[%s]", result)
	}
	result = utils.Getenv("test01", "test00")
	if result != "test01" {
		t.Errorf("Getenv test failed, expected[test01], got[%s]", result)
	}
}
func TestBoolToInt_withTrue_shouldReturnOne(t *testing.T) {
	result := utils.BoolToInt(true)
	expect := 1
	if result != expect {
		t.Errorf("BoolToInt test failed, expected[%d], got[%d]", expect, result)
	}
}

func TestBoolToInt_withFalse_shouldReturnZero(t *testing.T) {
	result := utils.BoolToInt(false)
	expect := 0
	if result != expect {
		t.Errorf("BoolToInt test failed, expected[%d], got[%d]", expect, result)
	}
}
func TestSafeAtoi_withIncorectFormatAndCorrectFormat_shouldUseDefaultIncaseOfIncorrectFormat(t *testing.T) {
	result := utils.SafeAtoi("123weq", 321)
	if result != 321 {
		t.Errorf("SafeAtoi test failed, expected[321], got[%d]", result)
	}
	result = utils.SafeAtoi("123", 321)
	if result != 123 {
		t.Errorf("SafeAtoi test failed, expected[123], got[%d]", result)
	}
}

func TestSafeBool_withIncorectFormatAndCorrectFormat_shouldUseDefaultIncaseOfIncorrectFormat(t *testing.T) {
	str := "true"
	expected := true

	result := utils.SafeBool(str, false)

	if result != expected {
		t.Errorf("Expected %t, but got %t", expected, result)
	}

	str = "invalid"
	expected = false

	result = utils.SafeBool(str, expected)

	if result != expected {
		t.Errorf("Expected %t, but got %t", expected, result)
	}
}

func TestFormatGatewayResponse_withExpectedSuccesfulOutput_shouldShowStatus200AndEmptyBody(t *testing.T) {
	statusCode := 200
	body := "{}"

	response := utils.FormatGatewayResponse(statusCode, body)

	if response.StatusCode != statusCode {
		t.Errorf("Expected StatusCode %d, but got %d", statusCode, response.StatusCode)
	}

	if response.Body != body {
		t.Errorf("Expected Body %s, but got %s", body, response.Body)
	}

	if response.Headers["Content-Type"] != "application/json" {
		t.Errorf("Expected Content-Type application/json, but got %s", response.Headers["Content-Type"])
	}

	if response.IsBase64Encoded != false {
		t.Errorf("Expected IsBase64Encoded false, but got true")
	}
}
func TestExtractIdAtPosition_withValidPosition_expectFound(t *testing.T) {
	id, found := utils.ExtractIdAtPosition("/providers/overrides/123", 3)
	if !found {
		t.Error("Expected id to be found")
	}
	if id != 123 {
		t.Errorf("Expected id to be 123 but got %d", id)
	}
}

func TestExtractIdAtPosition_withOutOfRangePosition_expectNotFound(t *testing.T) {
	id, found := utils.ExtractIdAtPosition("/providers/overrides/123", 5)
	if found {
		t.Error("Expected id to not be found")
	}
	if id != 0 {
		t.Errorf("Expected id to be 0 but got %d", id)
	}
}

func TestExtractIdAtPosition_withNegativePosition_expectNotFound(t *testing.T) {
	id, found := utils.ExtractIdAtPosition("/providers/overrides/123", -1)
	if found {
		t.Error("Expected id to not be found")
	}
	if id != 0 {
		t.Errorf("Expected id to be 0 but got %d", id)
	}
}

func TestExtractIdAtPosition_withEmptyPath_expectNotFound(t *testing.T) {
	id, found := utils.ExtractIdAtPosition("", -1)
	if found {
		t.Error("Expected id to not be found")
	}
	if id != 0 {
		t.Errorf("Expected id to be 0 but got %d", id)
	}
}

func TestExtractIdAtPosition_withNonNumericPath_expectNotFound(t *testing.T) {
	id, found := utils.ExtractIdAtPosition("/providers/overrides/abc", -1)
	if found {
		t.Error("Expected id to not be found")
	}
	if id != 0 {
		t.Errorf("Expected id to be 0 but got %d", id)
	}
}

func TestFormatErrorAPIGatewayResponse_withErrorParsingIntoFormatWrapper_shouldUseErrorsStatusCodeAndDescriptionNotDefault(t *testing.T) {
	personalErrorCode := 404
	personalErrorMessage := "Not Found"
	err := types.NewSocketError(personalErrorCode, personalErrorMessage)

	response := utils.FormatErrorAPIGatewayResponse(err)

	expectedStatusCode := personalErrorCode
	if response.StatusCode != expectedStatusCode {
		t.Errorf("Expected StatusCode %d, but got %d", expectedStatusCode, response.StatusCode)
	}

	expectedBody := fmt.Sprintf(`{"message":"%s"}`, personalErrorMessage)
	if response.Body != expectedBody {
		t.Errorf("Expected Body %s, but got %s", expectedBody, response.Body)
	}

	expectedContentType := "application/json"
	if response.Headers["Content-Type"] != expectedContentType {
		t.Errorf("Expected Content-Type %s, but got %s", expectedContentType, response.Headers["Content-Type"])
	}

	if response.IsBase64Encoded != false {
		t.Errorf("Expected IsBase64Encoded false, but got true")
	}
}

func TestGetCurrentDateFormatedForInsertingIntoDB(t *testing.T) {
	// Set a specific date and time for testing purposes
	testTime := time.Date(2023, time.November, 12, 15, 30, 45, 0, time.UTC)

	// Call the function with the test time
	result := utils.GetCurrentDateFormatedForInsertingIntoDB(testTime)

	// Define the expected result based on the test time
	expectedResult := "2023-11-12 15:30:45"

	// Check if the result matches the expected result
	if result != expectedResult {
		t.Errorf("Expected: %s, Got: %s", expectedResult, result)
	}
}

func TestFormatJSONString_ErrorEncoding(t *testing.T) {
	// Provide input that cannot be JSON encoded due to a function in the map
	input := `{"key": "value", "invalidType": ["UnsupportedField": "this is a channel"]}`
	result := utils.FormatJSONString(input)

	// The result should be the same as the input because encoding will fail
	if result != input {
		t.Errorf("Expected input '%s' to be returned as is, but got '%s'", input, result)
	}
}
func TestFormatJSONString_withValidJsonWithEnters_shouldReturnSingleLineJsonString(t *testing.T) {
	jsonPayload := `{
		"priority": 1,
		"enabled": true,
		"merchantId": 267,
		"acquirerId": 3,
		"providerId": 89,
		"currencyId": 13,
		"parameters": [
		  {
				"parameterName": "MyParameter",
				"parameterValue": "46856ujnjj",
				"enabled": true
		  }
		]
  }`
	expected := `{"acquirerId":3,"currencyId":13,"enabled":true,"merchantId":267,"parameters":[{"enabled":true,"parameterName":"MyParameter","parameterValue":"46856ujnjj"}],"priority":1,"providerId":89}`

	result := utils.FormatJSONString(jsonPayload)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
func TestFormatJSONString_withInvalidJsonText_shouldReturnSingleLineJsonString(t *testing.T) {
	jsonPayload := `invalid-json`

	result := utils.FormatJSONString(jsonPayload)

	if result != jsonPayload {
		t.Errorf("Expected %s, but got %s", jsonPayload, result)
	}
}
func TestFormatJSONString_withInvalidJsonStructure_shouldReturnSingleLineJsonString(t *testing.T) {
	jsonPayload := `{"key": "value", "missing_quote:}`

	result := utils.FormatJSONString(jsonPayload)

	if result != jsonPayload {
		t.Errorf("Expected %s, but got %s", jsonPayload, result)
	}
}
func TestGetIpv4Address_withValidAndEmptyCase_shouldReturnLocalForEmptyCaseAndStandIPforNormalCase(t *testing.T) {
	sourceIp := "192.168.1.1"
	expected := "192.168.1.1"

	result := utils.GetIpv4Address(sourceIp)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	sourceIp = ""
	expected = "127.0.0.1"

	result = utils.GetIpv4Address(sourceIp)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	sourceIp = "1234.5678.9012.3456"
	expected = "127.0.0.1"

	result = utils.GetIpv4Address(sourceIp)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestGetIpv6Address_withValidCaseAndInvalidCae_shouldReturnEmptyIfInvalid(t *testing.T) {
	sourceIp := "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
	expected := "2001:0db8:85a3:0000:0000:8a2e:0370:7334"

	result := utils.GetIpv6Address(sourceIp)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	sourceIp = "192.168.1.1"
	expected = ""

	result = utils.GetIpv6Address(sourceIp)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestReplaceQuotes_withCaseWithQuotesInAndCaseWithoutQuotes_shouldAlwaysRemoveQuotesIfPresent(t *testing.T) {
	value := `"quoted value"`
	expected := `quoted value`

	result := utils.ReplaceQuotes(value)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	value = `no quotes`
	expected = `no quotes`

	result = utils.ReplaceQuotes(value)

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
