package constants

// For a definition of each constant check the logging standards on http://handbook.int.blacklane.io/monitoring/kiev.html
const (
	FieldEvent = "event"

	HeaderForwardedFor = "X-Forwarded-For"
	HeaderRequestID    = "X-Request-Id"
	HeaderTrackingID   = "X-Tracking-Id"

	TraceID       = "trace_id"
	TagTrackingID = "tracking_id"
)

const (
	// TracerName is the name for the OpenTelemetry Tracer
	TracerName = "github.com/blacklane/go-libs/tracing/events"
)
