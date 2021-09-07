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
	// TODO: change it to something shorter like "blacklane-otel" or "blacklane-otel-go"
	TracerName = "github.com/blacklane/go-libs/otel"
)
