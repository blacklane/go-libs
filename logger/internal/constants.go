package internal

const (
	// Constants follow http://handbook.int.blacklane.io/monitoring/logging-standard.html

	// Field log constants

	FieldApplication     = "application"
	FieldDuration        = "duration_ms"
	FieldEvent           = "event"
	FieldHTTPStatus      = "http_status"
	FieldHost            = "host"
	FieldIP              = "ip"
	FieldParams          = "params"
	FieldPath            = "path"
	FieldRequestDuration = "duration_ms"
	FieldRequestID       = "request_id"
	FieldTimestamp       = "timestamp"
	FieldTraceID         = "trace_id"
	FieldTrackingID      = "tracking_id"
	FieldUserAgent       = "user_agent"
	FieldVerb            = "verb"

	TimeFieldFormat = "2006-01-02T15:04:05.000Z07:00"

	// Headers constants

	HeaderForwardedFor = "X-Forwarded-For"
	HeaderRequestID    = "X-Request-Id"
	HeaderTrackingID   = "X-Tracking-Id"
)
