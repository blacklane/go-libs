package field

import (
	"time"
)

const (
	FieldApplication       = "application"
	FieldBody              = "body"
	FieldDurationMs        = "duration_ms"
	FieldEvent             = "event"
	FieldHTTPStatus        = "http_status"
	FieldHost              = "host"
	FieldIP                = "ip"
	FieldPath              = "path"
	FieldRequestDurationMs = "request_duration_ms"
	FieldRequestID         = "request_id"
	FieldTraceID           = "trace_id"
	FieldTrackingID        = "tracking_id"
	FieldUserAgent         = "user_agent"
	FieldHTTPMethod        = "http_method"
)

func Application(name string) Field {
	return String(FieldApplication, name)
}

func Body(body string) Field {
	return String(FieldBody, body)
}

func DurationMs(d time.Duration) Field {
	return Duration(FieldDurationMs, d)
}

func Event(name string) Field {
	return String(FieldEvent, name)
}

func HTTPStatus(status int) Field {
	return Int64(FieldHTTPStatus, int64(status))
}

func Host(host string) Field {
	return String(FieldHost, host)
}

func IP(ip string) Field {
	return String(FieldIP, ip)
}

func Path(path string) Field {
	return String(FieldPath, path)
}

func RequestDurationMs(d time.Duration) Field {
	return Duration(FieldRequestDurationMs, d)
}

func RequestID(id string) Field {
	return String(FieldRequestID, id)
}

func TraceID(id string) Field {
	return String(FieldTraceID, id)
}

func TrackingID(id string) Field {
	return String(FieldTrackingID, id)
}

func UserAgent(ua string) Field {
	return String(FieldUserAgent, ua)
}

func HTTPMethod(method string) Field {
	return String(FieldHTTPMethod, method)
}
