package tasks

import (
	"context"
	"donation-mgmt/src/dal"
	"errors"
)

var ErrRetryable = errors.New("[retryable]")

type TaskHandler interface {
	// HandleTask performs a task and returns an error if the task failed. If the error should be retried,
	// it should wrap the error with ErrRetryable.
	HandleTask(ctx context.Context, task *dal.Task) error
}
