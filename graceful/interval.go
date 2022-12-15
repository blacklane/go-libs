package graceful

import (
	"context"
	"time"
)

type intervalTask struct {
	interval time.Duration
	task     func(context.Context) error
}

func NewIntervalTask(interval time.Duration, task func(context.Context) error) Task {
	return &intervalTask{
		interval: interval,
		task:     task,
	}
}

func (t *intervalTask) Start(ctx context.Context) error {
	for {
		select {
		case <-time.After(t.interval):
			if err := t.task(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (t *intervalTask) Stop(context.Context) error {
	return nil
}
