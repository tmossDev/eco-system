package logger

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"runtime"
	"strings"
	"time"

	"tmossDev.github.com/eco-system/shared-components/backend/lib/env"
)

var LogLevel = map[string]int8{
	"TRACE": 0,
	"DEBUG": 1,
	"INFO":  2,
	"WARN":  3,
	"ERROR": 4,
}

const (
	// TRACE Log level
	TRACE = "TRACE"
	// DEBUG Log level
	DEBUG = "DEBUG"
	// INFO Log level
	INFO = "INFO"
	// WARN Log level
	WARN = "WARN"
	// ERROR Log level
	ERROR = "ERROR"
	// NOSUMO level
	NOSUMO = "NOSUMO"
	// format to log
	logFormat = "%[1]s [%[2]s] [%[3]s] (%[5]s) %[4]s"
	// format to log info
	logInfoFormat = "%[1]s [%[2]s] [%[3]s] %[4]s"
	// tag name for redacting
	tagName = "redact"
	// MaskTypeLast4 mask type
	MaskTypeLast4 = "last4"
	// MaskTypeFirst6 mask type
	MaskTypeFirst6 = "first6"
	// MaskTypeFirst6Last4 mask type
	MaskTypeFirst6Last4 = "first6last4"
	// MaskTypeComplete mask type
	MaskTypeComplete = "complete"
	// SystemFailure log indicator
	SystemFailure = "System_Failure"
	// SystemInfo log indicator
	SystemInfo = "System_Info"
)

func First6MaskFunc(str string) string {
	if len(str) <= 6 {
		return strings.Repeat("*", len(str))
	}
	return str[:6] + strings.Repeat("*", len(str)-6)
}

func Last4MaskFunc(str string) string {
	if len(str) <= 4 {
		return strings.Repeat("*", len(str))
	}
	return strings.Repeat("*", len(str)-4) + str[len(str)-4:]
}

func First6Last4MaskFunc(str string) string {
	if len(str) <= 10 {
		return strings.Repeat("*", len(str))
	}
	return str[:6] + strings.Repeat("*", len(str)-10) + str[len(str)-4:]
}

func CompleteMaskFunc(str string) string {
	return strings.Repeat("*", len(str))
}

// GetMerchantWithClientIdAndMerchantRef new Value based on Tag
func redactGetTagValue(i int, value reflect.Value, ptr reflect.Value) string {
	var newValue string

	tag := value.Type().Field(i).Tag.Get(tagName)
	if tag == "" || tag == "-" {
		return ""
	}

	switch tag {
	case MaskTypeFirst6:
		newValue = First6MaskFunc(ptr.String())
	case MaskTypeLast4:
		newValue = Last4MaskFunc(ptr.String())
	case MaskTypeFirst6Last4:
		newValue = First6Last4MaskFunc(ptr.String())
	case MaskTypeComplete:
		newValue = CompleteMaskFunc(ptr.String())
	default:
		newValue = ptr.String()
	}

	return newValue
}

func src() string {
	// Determine caller func
	pc, file, lineno, ok := runtime.Caller(3)
	src := ""
	if ok {
		slice := strings.Split(runtime.FuncForPC(pc).Name(), "/")
		src = slice[len(slice)-1]
		slice = strings.Split(file, "/")
		file := slice[len(slice)-1]
		src = fmt.Sprintf("%s at %s:%d", src, file, lineno)
	}
	return src
}

func now() string {
	return time.Now().Format("2006-01-02T15:04:05")
}

func Init() {
	// clear all current logging flags
	log.SetFlags(0)
}

// func logIt(level string, requestId string, a ...interface{}) {
//	src := src()
//	now := now()
//
//	logItNow(level, requestId, nil, a, now, src)
// }

// Log Manual log using supplied level.
// If level is not known, then log as info
// Arguments are handled in the manner of fmt.Println.
func Log(level string, requestId string, a ...interface{}) {
	switch strings.ToUpper(level) {
	case TRACE:
		logItNow(TRACE, requestId, nil, a...)
	case DEBUG:
		logItNow(DEBUG, requestId, nil, a...)
	case WARN:
		logItNow(WARN, requestId, nil, a...)
	case ERROR:
		logItNow(ERROR, requestId, nil, a...)
	default:
		logItNow(INFO, requestId, nil, a...)
	}
}

// Println Used to log a debug level message
// Arguments are handled in the manner of fmt.Println.
func Println(requestId string, a ...interface{}) {
	logItNow(INFO, requestId, nil, a...)
}

// Trace Used to log a debug level message
// Arguments are handled in the manner of fmt.Println.
func Trace(requestId string, a ...interface{}) {
	logItNow(TRACE, requestId, nil, a...)
}

// Debug Used to log a debug level message
// Arguments are handled in the manner of fmt.Println.
func Debug(requestId string, a ...interface{}) {
	logItNow(DEBUG, requestId, nil, a...)
}

// Info Used to log a debug level message
// Arguments are handled in the manner of fmt.Println.
func Info(requestId string, a ...interface{}) {
	logItNow(INFO, requestId, nil, a...)
}

// Warn Used to log a debug level message
// Arguments are handled in the manner of fmt.Println.
func Warn(requestId string, a ...interface{}) {
	logItNow(WARN, requestId, nil, a...)
}

// Error Used to log a debug level message
// Arguments are handled in the manner of fmt.Println.
func Error(requestId string, a ...interface{}) {
	logItNow(ERROR, requestId, nil, a...)
}

// Fatal Used to log a debug level message using a format provided
func Fatal(requestId string, a ...interface{}) {
	logItNow("FATAL", requestId, nil, a...)
	panic(fmt.Sprint(a...))
}

// func logItf(level string, requestId string, format *string, a ...interface{}) {
//
//	logItNow(level, requestId, format, a...)
// }

func logItNow(level string, requestId string, format *string, a ...interface{}) {
	var msg string
	src := src()
	if ValidateAgainstConfiguredLogLevel(level) {
		// src := src()
		now := now()

		if len(requestId) == 0 {
			requestId = "N/A"
		}

		if format == nil {
			msg = fmt.Sprint(a...)
		} else {
			msg = fmt.Sprintf(*format, a...)
		}

		var finalLogFormat = logFormat
		if level == INFO {
			finalLogFormat = logInfoFormat
		}
		msg = fmt.Sprintf(finalLogFormat, now, requestId, level, msg, src)
		log.Print(msg)
	}
}

// Logf Manual log using supplied level.
// If level is not known, then log as info
// Arguments are handled in the manner of fmt.Println.
func Logf(level string, requestId string, format string, a ...interface{}) {
	switch strings.ToUpper(level) {
	case TRACE:
		logItNow(TRACE, requestId, &format, a...)
	case DEBUG:
		logItNow(DEBUG, requestId, &format, a...)
	case WARN:
		logItNow(WARN, requestId, &format, a...)
	case ERROR:
		logItNow(ERROR, requestId, &format, a...)
	default:
		logItNow(INFO, requestId, &format, a...)
	}
}

// Printf Used to log a debug level message
// Arguments are handled in the manner of fmt.Println.
func Printf(requestId string, format string, a ...interface{}) {
	logItNow(INFO, requestId, &format, a...)
}

// Tracef Used to log a debug level message using a format provided
func Tracef(requestId string, format string, a ...interface{}) {
	logItNow(TRACE, requestId, &format, a...)
}

// Debugf Used to log a debug level message using a format provided
func Debugf(requestId string, format string, a ...interface{}) {
	logItNow(DEBUG, requestId, &format, a...)
}

// Infof Used to log a debug level message using a format provided
func Infof(requestId string, format string, a ...interface{}) {
	logItNow(INFO, requestId, &format, a...)
}

// Warnf Used to log a debug level message using a format provided
func Warnf(requestId string, format string, a ...interface{}) {
	logItNow(WARN, requestId, &format, a...)
}

// Errorf Used to log a debug level message using a format provided
func Errorf(requestId string, format string, a ...interface{}) {
	logItNow(ERROR, requestId, &format, a...)
}

// Fatalf Used to log a debug level message using a format provided
func Fatalf(requestId string, format string, a ...interface{}) {
	logItNow("FATAL", requestId, &format, a...)
	panic(fmt.Sprintf(format, a...))
}

// ValidateAgainstConfiguredLogLevel a log level against the log level configured in the Environment Variable
// Levels in order: `TRACE`, `DEBUG`, `INFO`, `WARN`, `ERROR`
func ValidateAgainstConfiguredLogLevel(level string) bool {
	logLevelEnvVariable := env.Getenv("LOG_LEVEL", INFO)
	return LogLevel[level] >= LogLevel[logLevelEnvVariable]
}

func HttpLogger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			inner.ServeHTTP(w, r)

			Info(
				r.Context().Value("__request_id__").(string),
				map[string]interface{}{
					"method":     r.Method,
					"path":       r.RequestURI,
					"route_name": name,
					"time_taken": fmt.Sprint(time.Since(start)),
				},
			)
		},
	)
}
