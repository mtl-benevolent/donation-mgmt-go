package tasks_test

import (
	"context"
	"donation-mgmt/src/config"
	"donation-mgmt/src/dal"
	dalmocks "donation-mgmt/src/dal/mocks"
	"donation-mgmt/src/libs/logger"
	"donation-mgmt/src/tasks"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockHandler struct {
	called   *atomic.Int32
	panicVal any
	err      error
	handler  func(ctx context.Context, task *dal.Task) error // Optional custom handler
}

func (m *mockHandler) HandleTask(ctx context.Context, task *dal.Task) error {
	m.called.Add(1)
	if m.panicVal != nil {
		panic(m.panicVal)
	}

	// If custom handler is provided, use it
	if m.handler != nil {
		return m.handler(ctx, task)
	}

	return m.err
}

func Test_WhenTaskIsProcessed_ShouldCallHandlerAndAck(t *testing.T) {
	logger.BootstrapLogger(&config.AppConfiguration{LogLevel: "ERROR"})

	h := &mockHandler{called: &atomic.Int32{}}

	mockQuerier := dalmocks.NewQuerier(t)
	mockQuerier.On("PickTasks", mock.Anything, mock.Anything).Return([]dal.Task{{ID: 1, Type: "TEST", Attempt: 0, MaxRetries: 1}}, nil)
	mockQuerier.On("AckTasks", mock.Anything, []int64{1}).Return(int64(1), nil)

	q, err := tasks.NewQueue(mockQuerier, tasks.QueueConfig{
		QueueName:    "test",
		WorkerSlots:  1,
		WorkHandlers: tasks.TaskHandlerMap{"TEST": h},
		PollInterval: 10 * time.Millisecond,
		LockDuration: 100 * time.Millisecond,
	})
	require.NoError(t, err, "failed to create queue")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	go q.Start(ctx)

	time.Sleep(50 * time.Millisecond)
	assert.Greater(t, h.called.Load(), int32(0), "handler was not called")
	mockQuerier.AssertExpectations(t)
}

func Test_WhenHandlerReturnsRetryableError_ShouldNack(t *testing.T) {
	logger.BootstrapLogger(&config.AppConfiguration{LogLevel: "ERROR"})

	h := &mockHandler{called: &atomic.Int32{}, err: tasks.ErrRetryable}

	mockQuerier := dalmocks.NewQuerier(t)
	mockQuerier.On("PickTasks", mock.Anything, mock.Anything).Return([]dal.Task{{ID: 2, Type: "TEST", Attempt: 0, MaxRetries: 1}}, nil)
	mockQuerier.On("NackTask", mock.Anything, mock.Anything).Return(int64(1), nil)

	q, err := tasks.NewQueue(mockQuerier, tasks.QueueConfig{
		QueueName:    "test",
		WorkerSlots:  1,
		WorkHandlers: tasks.TaskHandlerMap{"TEST": h},
		PollInterval: 10 * time.Millisecond,
		LockDuration: 100 * time.Millisecond,
	})
	require.NoError(t, err, "failed to create queue")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	go q.Start(ctx)

	time.Sleep(50 * time.Millisecond)
	assert.Greater(t, h.called.Load(), int32(0), "handler was not called")
	mockQuerier.AssertExpectations(t)
}

func Test_WhenHandlerPanics_ShouldNack(t *testing.T) {
	logger.BootstrapLogger(&config.AppConfiguration{LogLevel: "ERROR"})

	h := &mockHandler{called: &atomic.Int32{}, panicVal: "panic!"}

	mockQuerier := dalmocks.NewQuerier(t)
	mockQuerier.On("PickTasks", mock.Anything, mock.Anything).Return([]dal.Task{{ID: 3, Type: "TEST", Attempt: 0, MaxRetries: 1}}, nil)
	mockQuerier.On("NackTask", mock.Anything, mock.Anything).Return(int64(1), nil)

	q, err := tasks.NewQueue(mockQuerier, tasks.QueueConfig{
		QueueName:    "test",
		WorkerSlots:  1,
		WorkHandlers: tasks.TaskHandlerMap{"TEST": h},
		PollInterval: 10 * time.Millisecond,
		LockDuration: 100 * time.Millisecond,
	})
	require.NoError(t, err, "failed to create queue")

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	go q.Start(ctx)

	time.Sleep(50 * time.Millisecond)
	assert.Greater(t, h.called.Load(), int32(0), "handler was not called")
	mockQuerier.AssertExpectations(t)
}

func Test_WhenMultipleTasksAvailable_ShouldOnlyPullUpToWorkerSlots(t *testing.T) {
	logger.BootstrapLogger(&config.AppConfiguration{LogLevel: "ERROR"})

	workerSlots := 2
	totalTasks := 5

	// Create tasks for the mock to return
	taskList := make([]dal.Task, totalTasks)
	for i := 0; i < totalTasks; i++ {
		taskList[i] = dal.Task{ID: int64(i + 1), Type: "TEST", Attempt: 0, MaxRetries: 1}
	}

	mockQuerier := dalmocks.NewQuerier(t)
	mockQuerier.On("PickTasks", mock.Anything, mock.Anything).Return(taskList, nil)
	mockQuerier.On("AckTasks", mock.Anything, mock.Anything).Return(int64(1), nil)

	done := make(chan struct{})

	// Use a custom handler that holds the task until signaled
	handler := &mockHandler{
		called: &atomic.Int32{},
		handler: func(ctx context.Context, task *dal.Task) error {
			// Wait for the done channel to be closed, or the context to be done
			select {
			case <-done:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		},
	}

	q, err := tasks.NewQueue(mockQuerier, tasks.QueueConfig{
		QueueName:    "test",
		WorkerSlots:  workerSlots,
		WorkHandlers: tasks.TaskHandlerMap{"TEST": handler},
		PollInterval: 10 * time.Millisecond,
		LockDuration: 100 * time.Millisecond,
	})
	require.NoError(t, err, "failed to create queue")
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	go q.Start(ctx)
	time.Sleep(30 * time.Millisecond)

	close(done) // Signal all handlers to complete

	// Only up to workerSlots tasks should be processed at a time
	assert.LessOrEqual(t, int(handler.called.Load()), workerSlots, "should not process more tasks than worker slots at a time")
}
