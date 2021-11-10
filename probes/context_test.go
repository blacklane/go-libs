package probes

import (
	"context"
	"reflect"
	"testing"
)

func TestFromCtx(t *testing.T) {
	p := New("")

	ctx := p.WithContext(context.Background())

	got := FromCtx(ctx)

	if !reflect.DeepEqual(got, p) {
		t.Errorf("want: %T, got: %T", p, got)
	}
}

func TestFromCtx_returnsNoop(t *testing.T) {
	want := NewNoop()

	got := FromCtx(context.Background())

	if !reflect.DeepEqual(got, want) {
		t.Errorf("want: %T, got: %T", want, got)
	}
}
