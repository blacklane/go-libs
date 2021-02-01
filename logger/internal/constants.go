package internal

const (
	// Field log constants
	FieldTrackingID = "tracking_id"
	FieldDuration   = "duration_ms"

	// Fields, Event and Headers constants following http://handbook.int.blacklane.io/monitoring/kiev.html

	FieldApplication  = "application"
	FieldTimestamp    = "timestamp"
	FieldEntryPoint   = "entry_point"
	FieldRequestID    = "request_id"
	FieldRequestDepth = "request_depth"
	FieldTreePath     = "tree_path"
	FieldRoute        = "route"
	FieldHost         = "host"
	FieldVerb         = "verb"
	FieldPath         = "path"

	FieldEvent = "event"

	FieldParams          = "params"
	FieldIP              = "ip"
	FieldUserAgent       = "user_agent"
	FieldStatus          = "status"
	FieldRequestDuration = "request_duration"

	TimeFieldFormat = "2006-01-02T15:04:05.000Z07:00"

	// Event constants

	// EventRequestFinished should be used as 'event' when logging a finished request/job
	EventRequestFinished = "request_finished"

	// Headers constants
	HeaderForwardedFor = "X-Forwarded-For"
	HeaderRequestID    = "X-Request-Id"
	HeaderRequestDepth = "X-Request-Depth"
	HeaderTrackingID   = "X-Tracking-Id"
	HeaderTreePath     = "X-Tree-Path"
)
