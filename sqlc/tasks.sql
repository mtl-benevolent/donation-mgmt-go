-- name: CreateTask :one
INSERT INTO task_queue(
    type, body, max_retries
) VALUES(sqlc.Arg('Type'), sqlc.Arg('Body'), sqlc.Arg('MaxRetries'))
RETURNING *;

-- name: PickTasks :many
WITH tasks_to_process AS (
    SELECT * FROM task_queue tq
	WHERE tq."type" = sqlc.slice('SupportedTaskTypes')::"TaskType"[]
		AND tq.status IN ('CREATED', 'IN_PROGRESS', 'ERROR_RETRYABLE')
		AND (tq.locked_until IS NULL OR tq.locked_until <= NOW())
		AND tq.attempt < tq.max_retries
	ORDER BY tq.created_at ASC
	LIMIT sqlc.arg('WorkerSlots')
	FOR UPDATE SKIP LOCKED
)
UPDATE task_queue tq
SET
    status = 'IN_PROGRESS',
    last_picked_up_at = NOW(),
    locked_until = NOW() + sqlc.Arg('LockDuration')::interval,
    locked_by = sqlc.Arg('ProcessName'),
    attempt = tq.attempt + 1
FROM tasks_to_process tp
WHERE tq.id = tp.id
RETURNING tq.*;

-- name: AckTasks :execrows
UPDATE task_queue tq
SET
	status = 'COMPLETED',
	last_error_message = NULL,
	completed_at = NOW(),
	locked_until = NULL,
	locked_by = NULL
WHERE tq.id = sqlc.slice('TaskIDs');

-- name: NackTask :execrows
-- NackTask
WITH tasks_to_process AS (
	SELECT 
		*, 
		(tq.attempt < tq.max_retries AND sqlc.arg('TaskStatus')::"TaskStatus" NOT IN ('ERROR_UNRETRYABLE')) AS is_retryable 
	FROM task_queue tq
	WHERE tq.id = 1
)
UPDATE task_queue tq
SET
	status = CASE
		WHEN tp.is_retryable THEN sqlc.arg('TaskStatus')::"TaskStatus"
		ELSE 'ERROR_UNRETRYABLE'::"TaskStatus"
	END,
	last_error_message = CASE
		WHEN tq.attempt < tq.max_retries THEN sqlc.narg('ErrorMessage')
		WHEN (COALESCE(sqlc.narg('ErrorMessage'), '') = '') THEN 'Max attempts reached'
        ELSE CONCAT('Max attempts reached:', ' ', sqlc.narg('ErrorMessage'))
	END,
	locked_until = CASE
		WHEN tp.is_retryable AND NOT ISNULL(sqlc.narg('RetryIn')::INTERVAL) THEN NOW() + sqlc.narg('RetryIn')::INTERVAL
		ELSE NULL
	END,
	locked_by = CASE
		WHEN tp.is_retryable THEN sqlc.narg('ProcessName')
		ELSE NULL
	END
FROM tasks_to_process tp
WHERE tq.id = tp.id;
