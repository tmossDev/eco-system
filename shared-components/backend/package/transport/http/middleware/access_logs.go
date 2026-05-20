package middleware

import (
	"encoding/json"
	"io"
	"os"
	"strings"

	"github.com/kataras/iris/v12/middleware/accesslog"
)

const redactedLogValue = "[redacted]"

type noCloseWriter struct {
	io.Writer
}

type redactingJSONFormatter struct {
	formatter accesslog.JSON
}

func (writer noCloseWriter) Close() error {
	return nil
}

func (formatter *redactingJSONFormatter) SetOutput(dest io.Writer) {
	formatter.formatter.SetOutput(dest)
}

func (formatter *redactingJSONFormatter) Format(log *accesslog.Log) (bool, error) {
	redactedLog := *log
	redactedLog.Request = redactAccessLogBody(log.Request)
	redactedLog.Response = redactAccessLogBody(log.Response)

	return formatter.formatter.Format(&redactedLog)
}

func MakeAccessLog() *accesslog.AccessLog {
	// Initialize a new access log middleware.
	var ac = accesslog.File("./access.log")
	ac.AddOutput(noCloseWriter{Writer: os.Stdout})

	// The default configuration:
	ac.Delim = '|'
	ac.TimeFormat = "2006-01-02 15:04:05"
	ac.Async = false
	ac.IP = true
	ac.BytesReceivedBody = true
	ac.BytesSentBody = true
	ac.BytesReceived = false
	ac.BytesSent = false
	ac.BodyMinify = true
	ac.RequestBody = true
	ac.ResponseBody = true
	ac.KeepMultiLineError = true
	ac.PanicLog = accesslog.LogHandler

	// Default line format if formatter is missing:
	// Time|Latency|Code|Method|Path|IP|Path Params Query Fields|Bytes Received|Bytes Sent|Request|Response|
	//
	// Set Custom Formatter:
	ac.SetFormatter(&redactingJSONFormatter{
		formatter: accesslog.JSON{
			Indent:    "  ",
			HumanTime: true,
		},
	})
	// ac.SetFormatter(&accesslog.CSV{})
	// ac.SetFormatter(&accesslog.Template{Text: "{{.Code}}"})

	return ac
}

func redactAccessLogBody(body string) string {
	if body == "" {
		return body
	}

	var payload any
	if err := json.Unmarshal([]byte(body), &payload); err != nil {
		return body
	}

	redactAccessLogValue(payload)

	redactedBody, err := json.Marshal(payload)
	if err != nil {
		return body
	}

	return string(redactedBody)
}

func redactAccessLogValue(value any) {
	switch typedValue := value.(type) {
	case map[string]any:
		for key, childValue := range typedValue {
			if isSensitiveAccessLogKey(key) {
				typedValue[key] = redactedLogValue
				continue
			}

			redactAccessLogValue(childValue)
		}
	case []any:
		for _, childValue := range typedValue {
			redactAccessLogValue(childValue)
		}
	}
}

func isSensitiveAccessLogKey(key string) bool {
	normalizedKey := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(key, "_", ""), "-", ""))

	return normalizedKey == "password" ||
		normalizedKey == "confirmpassword" ||
		normalizedKey == "jwt" ||
		normalizedKey == "accesstoken" ||
		normalizedKey == "refreshtoken" ||
		normalizedKey == "token" ||
		normalizedKey == "authorization"
}
