package tracking

import (
	"context"
	"testing"
)

func TestContextWithID(t *testing.T) {
	ctx := ContextWithID(context.Background())

	got, ok := ctx.Value(ctxKey).(string)
	if !ok {
		t.Errorf("could not convert %v to string", ctx.Value(ctxKey))
	}

	if got == "" {
		t.Errorf("want an UUID, got empty string")
	}
}

func TestSetContextID(t *testing.T) {
	want := "42"
	ctx := SetContextID(context.Background(), want)

	got, ok := ctx.Value(ctxKey).(string)
	if !ok {
		t.Errorf("could not convert %v to string", ctx.Value(ctxKey))
	}

	if got == "" {
		t.Errorf("want an UUID, got empty string")
	}
}

func TestSetContextIDOverridesID(t *testing.T) {
	want := "42"

	ctx := ContextWithID(context.Background())
	ctx = SetContextID(ctx, want)

	got, ok := ctx.Value(ctxKey).(string)
	if !ok {
		t.Errorf("could not convert %v to string", ctx.Value(ctxKey))
	}

	if got == "" {
		t.Errorf("want an UUID, got empty string")
	}
}

func TestIdFromContext(t *testing.T) {
	want := "uuid"
	ctx := context.WithValue(context.Background(), ctxKey, want)

	got := IDFromContext(ctx)

	if got != want {
		t.Errorf("got: %s, wnat %s", got, want)
	}
}
