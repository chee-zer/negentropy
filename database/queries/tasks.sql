-- name: CreateTask :one
INSERT INTO tasks (name, color_hex, daily_target, completed)
VALUES (?, ?, ?, FALSE)
RETURNING *;

-- name: GetTasks :many
SELECT *
FROM tasks
ORDER BY id;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = ?;