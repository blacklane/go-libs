package graceful

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func ai32equal(t *testing.T, value *atomic.Int32, expected int32, msg string) {
	if got := value.Load(); got != expected {
		t.Errorf("%s, got: %d, want: %d", msg, got, expected)
	}
}

func counterHook(t *testing.T, msg string, value *atomic.Int32) Hook {
	return func(ctx context.Context) error {
		t.Log(msg)
		value.Add(1)
		return nil
	}
}

func TestRun(t *testing.T) {
	beforeStart := new(atomic.Int32)
	afterStop := new(atomic.Int32)
	taskStart := new(atomic.Int32)
	taskStop := new(atomic.Int32)
	intervalTask := new(atomic.Int32)

	g := New(
		WithBeforeStartHooks(
			counterHook(t, "before start 1", beforeStart),
			counterHook(t, "before start 2", beforeStart),
		),
		WithAfterStopHooks(
			counterHook(t, "after stop", afterStop),
		),
		WithTasks(
			NewTask(
				counterHook(t, "task start", taskStart),
				counterHook(t, "task stop", taskStop),
			),
			NewIntervalTask(300*time.Millisecond, counterHook(t, "interval task", intervalTask)),
		),
	)

	time.AfterFunc(time.Second, func() {
		if err := g.Stop(); err != nil {
			t.Fatal(err)
		}
	})

	if err := g.Run(); err != nil {
		t.Fatal(err)
	}

	ai32equal(t, taskStart, 1, "unexpected calls to task start")
	ai32equal(t, taskStop, 1, "unexpected calls to task stop")
	ai32equal(t, intervalTask, 3, "unexpected calls to interval task")
	ai32equal(t, beforeStart, 2, "unexpected calls to before start")
	ai32equal(t, afterStop, 1, "unexpected calls to after stop")
}

func TestRun_FailOnBeforeStartError(t *testing.T) {
	taskStart := new(atomic.Int32)
	taskStop := new(atomic.Int32)

	errBeforeStart := errors.New("before start error")
	errAfterStop := errors.New("after stop error")

	g := New(
		WithBeforeStartHooks(func(ctx context.Context) error {
			return errBeforeStart
		}),
		WithAfterStopHooks(func(ctx context.Context) error {
			return errAfterStop
		}),
		WithTasks(
			NewTask(
				counterHook(t, "task start", taskStart),
				counterHook(t, "task stop", taskStop),
			),
		),
	)

	err := g.Run()

	if !errors.Is(err, errBeforeStart) {
		t.Errorf("invalid error result, got: %v, want: %v", err, errBeforeStart)
	}

	if !errors.Is(err, errAfterStop) {
		t.Errorf("invalid error result, got: %v, want: %v", err, errAfterStop)
	}

	ai32equal(t, taskStart, 0, "unexpected calls to task start")
	ai32equal(t, taskStop, 0, "unexpected calls to task stop")
}

func TestRun_StopOnTaskError(t *testing.T) {
	task1Stop := new(atomic.Int32)
	task2Start := new(atomic.Int32)

	errTask1 := errors.New("task-1 error")
	errTask2 := errors.New("task-2 error")

	g := New(
		WithTasks(
			NewTask(
				func(ctx context.Context) error {
					return errTask1
				},
				counterHook(t, "task-1 stop", task1Stop),
			),
			NewTask(
				counterHook(t, "task-2 start", task2Start),
				func(ctx context.Context) error {
					return errTask2
				},
			),
		),
	)

	err := g.Run()

	if !errors.Is(err, errTask1) {
		t.Errorf("invalid error result, got: %v, want: %v", err, errTask1)
	}

	if !errors.Is(err, errTask2) {
		t.Errorf("invalid error result, got: %v, want: %v", err, errTask2)
	}

	ai32equal(t, task1Stop, 1, "unexpected calls to task-1 start")
	ai32equal(t, task2Start, 1, "unexpected calls to task-2 stop")
}

func TestRun_Append(t *testing.T) {
	beforeStart := new(atomic.Int32)
	afterStop := new(atomic.Int32)
	taskStart := new(atomic.Int32)
	taskStop := new(atomic.Int32)
	intervalTask := new(atomic.Int32)

	g := New()

	failIfError(t, g.AppendBeforeStartHook(counterHook(t, "before start 1", beforeStart)))
	failIfError(t, g.AppendBeforeStartHook(counterHook(t, "before start 2", beforeStart)))
	failIfError(t, g.AppendAfterStopHook(counterHook(t, "after stop", afterStop)))

	failIfError(t, g.AppendTask(NewTask(
		counterHook(t, "task start", taskStart),
		counterHook(t, "task stop", taskStop))),
	)

	failIfError(t, g.AppendTask(NewIntervalTask(
		300*time.Millisecond,
		counterHook(t, "interval task", intervalTask))),
	)

	time.AfterFunc(time.Second, func() {
		if err := g.Stop(); err != nil {
			t.Fatal(err)
		}
	})

	if err := g.Run(); err != nil {
		t.Fatal(err)
	}

	ai32equal(t, taskStart, 1, "unexpected calls to task start")
	ai32equal(t, taskStop, 1, "unexpected calls to task stop")
	ai32equal(t, intervalTask, 3, "unexpected calls to interval task")
	ai32equal(t, beforeStart, 2, "unexpected calls to before start")
	ai32equal(t, afterStop, 1, "unexpected calls to after stop")
}

func TestRun_AppendAfterRunningError(t *testing.T) {
	taskStart := new(atomic.Int32)
	taskStop := new(atomic.Int32)

	g := New()

	var err error
	go func() {
		err = g.AppendTask(NewTask(
			counterHook(t, "task start", taskStart),
			counterHook(t, "task stop", taskStop),
		))
	}()

	time.AfterFunc(time.Second, func() {
		if err := g.Stop(); err != nil {
			t.Fatal(err)
		}
	})

	if err := g.Run(); err != nil {
		t.Fatal(err)
	}

	if !errors.Is(err, ErrAlreadyRunning) {
		t.Errorf("invalid error result, got: %v, want: %v", err, ErrAlreadyRunning)
	}

	ai32equal(t, taskStart, 0, "unexpected calls to task start")
	ai32equal(t, taskStop, 0, "unexpected calls to task stop")
}

func failIfError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("got error: %v", err)
	}
}
