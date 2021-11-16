package tracking

import (
	"context"

	guuid "github.com/google/uuid"
)

type key struct{}

var ctxKey = key{}

// SetContextID returns a copy of the parent context with the id set to id.
// If an ID was already set, it will be overridden.
func SetContextID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey, id)
}

// ContextWithID returns a copy of the parent context with a Version 1 UUID on it.
// If it fails when generating the uuid an empty string is used instead.
func ContextWithID(ctx context.Context) context.Context {
	var id string

	if uuid, err := guuid.NewUUID(); err == nil {
		id = uuid.String()
	}

	return context.WithValue(ctx, ctxKey, id)
}

// IDFromContext returns the ID in the context, is none is found an empty string is returned.
// Be aware an empty string might mean there was a failure when generating the uuid, see ContextWithID
// for more details
func IDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(ctxKey).(string)
	return id
}
