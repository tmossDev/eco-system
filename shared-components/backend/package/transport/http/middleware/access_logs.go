package middleware

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/accesslog"
	"tmossDev.github.com/eco-system/shared-components/backend/package/constants"
)

const redactedLogValue = "[redacted]"

type noCloseWriter struct {
	io.Writer
}

type compactJSONFormatter struct {
	output io.Writer
}

func (writer noCloseWriter) Close() error {
	return nil
}

func (formatter *compactJSONFormatter) SetOutput(dest io.Writer) {
	formatter.output = dest
}

func (formatter *compactJSONFormatter) Format(log *accesslog.Log) (bool, error) {
	payload := accessLogPayload{
		Timestamp:     log.Now.Format(log.TimeFormat),
		Latency:       int64(log.Latency),
		Code:          log.Code,
		Method:        log.Method,
		Path:          log.Path,
		IP:            log.IP,
		RequestID:     getAccessLogField(log, "request_id"),
		Request:       redactAccessLogBody(log.Request),
		Response:      redactAccessLogBody(log.Response),
		BytesReceived: log.BytesReceived,
		BytesSent:     log.BytesSent,
	}

	line, err := json.Marshal(payload)
	if err != nil {
		return true, err
	}

	line = append(line, '\n')
	_, err = formatter.output.Write(line)

	return true, err
}

type accessLogPayload struct {
	Timestamp     string `json:"timestamp"`
	Latency       int64  `json:"latency"`
	Code          int    `json:"code"`
	Method        string `json:"method"`
	Path          string `json:"path"`
	IP            string `json:"ip,omitempty"`
	RequestID     string `json:"request_id,omitempty"`
	Request       string `json:"request,omitempty"`
	Response      string `json:"response,omitempty"`
	BytesReceived int    `json:"bytes_received,omitempty"`
	BytesSent     int    `json:"bytes_sent,omitempty"`
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
	ac.AddFields(func(ctx *context.Context, fields *accesslog.Fields) {
		fields.Set("request_id", ctx.Values().GetString(constants.CTXRequestIdKey))
	})

	// Default line format if formatter is missing:
	// Time|Latency|Code|Method|Path|IP|Path Params Query Fields|Bytes Received|Bytes Sent|Request|Response|
	//
	// Set Custom Formatter:
	ac.SetFormatter(&compactJSONFormatter{})
	// ac.SetFormatter(&accesslog.CSV{})
	// ac.SetFormatter(&accesslog.Template{Text: "{{.Code}}"})

	return ac
}

func getAccessLogField(log *accesslog.Log, key string) string {
	for _, field := range log.Fields {
		if field.Key == key {
			return fmt.Sprint(field.ValueRaw)
		}
	}

	return ""
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
