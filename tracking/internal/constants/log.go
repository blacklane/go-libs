package constants

const (
	// Fields constants
	// For a definition of each constant check the logging standards on http://handbook.int.blacklane.io/monitoring/kiev.html

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

	FieldErrorClass   = "error_class"
	FieldErrorMessage = "error_message"
	FieldBody         = "body"

	// Event constants

	// EventRequestFinished is the 'event' logged when the request finishes
	EventRequestFinished = "request_finished"
)
