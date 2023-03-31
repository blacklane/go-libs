package internal

// Internal constants
const (
	FieldApplication     = "application"
	FieldBody            = "body"
	FieldDuration        = "duration_ms"
	FieldEvent           = "event"
	FieldEventKey        = "event_key"
	FieldTopic           = "topic"
	FieldPartition       = "partition"
	FieldOffset          = "offset"
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
	FieldEnv             = "env"

	TimeFieldFormat = "2006-01-02T15:04:05.000Z07:00"

	// Headers constants

	HeaderForwardedFor = "X-Forwarded-For"
	HeaderRequestID    = "X-Request-Id"
	HeaderTrackingID   = "X-Tracking-Id"
)

type key string

const (
	FilterTag      = "[FILTERED]"
	FilterKeys key = "filterKeys"
)

// DefaultKeys Please don't mutate it!
var DefaultKeys = []string{"new_payer_email", "new_booker_email", "email", "new_booker_mobile_phone",
	"old_booker_new_booker_mobile_phone", "new_payer_mobile_phone", "mobile_phone", "phone"}
