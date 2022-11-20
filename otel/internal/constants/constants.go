package constants

//revive:disable:exported obvious internal constants.
const (
	LogKeyTraceID   = "trace_id"
	LogKeyEventName = "event_name"
)

const (
	// TracerName is the name for the OpenTelemetry Tracer
	// TODO: change it to something shorter like "blacklane-otel" or "blacklane-otel-go"?
	TracerName = "github.com/blacklane/go-libs/otel"
)
