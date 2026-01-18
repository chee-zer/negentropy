-- name: CreateTask :one
INSERT INTO tasks (name, color_hex, daily_target, completed)
VALUES (?, ?, ?, FALSE)
RETURNING *;

-- name: GetTasks :many
SELECT *
FROM tasks
ORDER BY id;

-- name: GetHours :one
SELECT SUM(daily_target)
FROM tasks
WHERE daily_target != "ENTROPY";

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = ?;
