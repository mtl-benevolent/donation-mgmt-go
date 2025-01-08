-- name: CreateTask :one
INSERT INTO task_queue(
    type, body, max_retries 
) VALUES(sqlc.Arg('Type'), sqlc.Arg('Body'), sqlc.Arg('MaxRetries'))
RETURNING *;
