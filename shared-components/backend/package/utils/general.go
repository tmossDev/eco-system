package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"time"

	"github.com/aws/aws-lambda-go/events"
	"tmossDev.github.com/eco-system/shared-components/backend/package/constants"
	"tmossDev.github.com/eco-system/shared-components/backend/package/types"
)

func ExtractIdAtPosition(path string, position int) (int, bool) {
	parts := strings.Split(path, "/")
	if position >= 0 && position < len(parts) {
		intValue, err := strconv.Atoi(parts[position])
		if err == nil {
			return intValue, true
		}
	}
	return 0, false
}

func Getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func SafeAtoi(str string, fallback int) int {
	value, err := strconv.Atoi(str)
	if err != nil {
		return fallback
	}
	return value
}

func SafeBool(str string, fallback bool) bool {
	value, err := strconv.ParseBool(str)
	if err != nil {
		return fallback
	}
	return value
}

func FormatGatewayResponse(statusCode int, body string) *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body:            body,
		IsBase64Encoded: false,
	}
}

func GetCurrentDateFormatedForInsertingIntoDB(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func FormatErrorAPIGatewayResponse(err error) *events.APIGatewayProxyResponse {

	statusCode := constants.InternalServerErrorCode
	errorDescription := "An unexpected error has occurred"
	if err, ok := err.(*types.SocketError); ok {
		statusCode = err.StatusCode()
		errorDescription = err.Error()
	}

	body := fmt.Sprintf(`{"message":"%s"}`, errorDescription)

	return FormatGatewayResponse(statusCode, body)
}

func GetIpv4Address(sourceIp string) string {
	if sourceIp == "" || len(sourceIp) > 15 {
		return "127.0.0.1"
	}
	return sourceIp
}

func GetIpv6Address(sourceIp string) string {
	if len(sourceIp) > 15 {
		return sourceIp
	}
	return ""
}

func ReplaceQuotes(value string) string {
	return strings.Replace(value, "\"", "", -1)
}

func BoolToInt(varaible bool) int {
	if varaible {
		return 1
	}
	return 0
}

func UintToString(variable uint64) string {
	return fmt.Sprintf("%d", variable)
}

func FormatJSONString(input string) string {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return input
	}
	formattedJSON, _ := json.Marshal(data)
	return strings.ReplaceAll(string(formattedJSON), "\n", " ")
}
