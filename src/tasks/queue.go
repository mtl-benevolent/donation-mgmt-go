package tasks

import (
	"context"
	"donation-mgmt/src/dal"
	"donation-mgmt/src/libs/logger"
	"donation-mgmt/src/ptr"
	"donation-mgmt/src/system/logging"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type TaskHandlerMap map[dal.TaskType]TaskHandler

type Queue struct {
	l *slog.Logger

	db dal.Querier

	queueName      string
	handlers       TaskHandlerMap
	workerSlots    int
	pollInterval   time.Duration
	supportedTypes []dal.TaskType
	taskChan       chan *dal.Task
	lockDuration   time.Duration

	busyWorkers *atomic.Int64
	wg          sync.WaitGroup
}

type QueueConfig struct {
	QueueName string

	WorkerSlots  int
	WorkHandlers TaskHandlerMap
	PollInterval time.Duration
	LockDuration time.Duration
}

func NewQueue(db dal.Querier, config QueueConfig) (*Queue, error) {
	if config.QueueName == "" {
		config.QueueName = fmt.Sprintf("default-%s", uuid.New().String())
	}

	if config.WorkerSlots <= 0 {
		config.WorkerSlots = 1
	}

	if config.PollInterval <= 0 {
		config.PollInterval = 5 * time.Second
	}

	if config.LockDuration <= 0 {
		config.LockDuration = 10 * time.Second
	}

	if len(config.WorkHandlers) == 0 {
		return nil, fmt.Errorf("no work handlers provided")
	}

	supportedTypes := make([]dal.TaskType, 0, len(config.WorkHandlers))
	for taskType := range config.WorkHandlers {
		supportedTypes = append(supportedTypes, taskType)
	}

	return &Queue{
		l: logger.ForComponent(fmt.Sprintf("tasks.Queue.%s", config.QueueName)),

		db:             db,
		queueName:      config.QueueName,
		handlers:       config.WorkHandlers,
		workerSlots:    config.WorkerSlots,
		pollInterval:   config.PollInterval,
		supportedTypes: supportedTypes,
		taskChan:       make(chan *dal.Task, config.WorkerSlots),
		lockDuration:   config.LockDuration,
		busyWorkers:    &atomic.Int64{},
		wg:             sync.WaitGroup{},
	}, nil
}

func (q *Queue) Start(ctx context.Context) {
	// Start workers
	for i := 0; i < q.workerSlots; i++ {
		go q.worker(ctx)
	}

	// Start polling loop
	ticker := time.NewTicker(q.pollInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			close(q.taskChan)
			q.wg.Wait() // Wait for all in-progress tasks to finish
			return
		case <-ticker.C:
			q.dequeue(ctx)
		}
	}
}

func (q *Queue) dequeue(ctx context.Context) {
	idleWorkers := int64(q.workerSlots) - q.busyWorkers.Load()
	if idleWorkers <= 0 {
		return
	}

	tasks, err := q.db.PickTasks(ctx, dal.PickTasksParams{
		SupportedTaskTypes: q.supportedTypes,
		WorkerSlots:        int32(idleWorkers),
		LockDuration:       pgtype.Interval{Microseconds: q.lockDuration.Microseconds(), Valid: true},
		ProcessName:        ptr.Wrap(q.queueName),
	})
	if err != nil {
		q.l.Error("failed to pick tasks", logging.ErrorKey, err)
		return
	}

	if len(tasks) == 0 {
		q.l.Debug("no tasks to pick")
		return
	}

	for _, task := range tasks {
		q.taskChan <- &task
	}
}

func (q *Queue) worker(ctx context.Context) {
	for task := range q.taskChan {
		lockExpiration := time.Now().Add(q.lockDuration)
		if task.LockedUntil != nil {
			lockExpiration = *task.LockedUntil
		}

		taskCtx, cancel := context.WithDeadline(ctx, lockExpiration)
		q.processTaskWithRecovery(taskCtx, task)
		cancel()
	}
}

func (q *Queue) processTaskWithRecovery(ctx context.Context, task *dal.Task) {
	q.wg.Add(1)
	defer q.wg.Done()

	q.busyWorkers.Add(1)
	defer q.busyWorkers.Add(-1)

	defer func() {
		if r := recover(); r != nil {
			q.l.Error("panic in task handler", "task_id", task.ID, "panic", r)
			q.nackTask(ctx, task, fmt.Errorf("panic: %v", r))
		}
	}()

	q.processTask(ctx, task)
}

func (q *Queue) processTask(ctx context.Context, task *dal.Task) {
	handler, ok := q.handlers[task.Type]
	if !ok {
		q.nackTask(ctx, task, fmt.Errorf("unknown task type"))
		return
	}
	err := handler.HandleTask(ctx, task)
	if err == nil {
		if err := q.ackTask(ctx, task); err != nil {
			q.l.Error("failed to ack task", logging.ErrorKey, err)
			q.nackTask(ctx, task, err)
		}
		return
	}

	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		q.l.Error("task lock expired", "task_id", task.ID)

		// We don't NACK here, because the task is already potentially being processed by another worker.
		return
	}

	q.l.Error("error processing task", "task_id", task.ID, logging.ErrorKey, err, "attempt", task.Attempt, "max_retries", task.MaxRetries, "retriable", errors.Is(err, ErrRetryable))
	q.nackTask(ctx, task, err)
}

func (q *Queue) ackTask(ctx context.Context, task *dal.Task) error {
	_, err := q.db.AckTasks(ctx, []int64{task.ID})
	if err != nil {
		return fmt.Errorf("%w failed to ackowledge task: %w", ErrRetryable, err)
	}

	return nil
}

func (q *Queue) nackTask(ctx context.Context, task *dal.Task, err error) {
	taskStatus := dal.TaskStatusERRORUNRETRYABLE
	retryIn := pgtype.Interval{
		Valid: false,
	}

	if errors.Is(err, ErrRetryable) {
		taskStatus = dal.TaskStatusERRORRETRYABLE
		retryIn = pgtype.Interval{
			Valid:        true,
			Microseconds: int64(task.Attempt) * 5 * int64(time.Second),
		}
	}

	_, nackErr := q.db.NackTask(ctx, dal.NackTaskParams{
		TaskStatus:   taskStatus,
		ErrorMessage: ptr.Wrap(err.Error()),
		RetryIn:      retryIn,
		ProcessName:  ptr.Wrap(q.queueName),
	})

	if nackErr != nil {
		q.l.Error("failed to nack task", logging.ErrorKey, err)
	}
}
