-- name: StartSession :one
INSERT INTO sessions (start_time, task_id)
VALUES (?, ?)
RETURNING *;

-- name: GetDailyTaskDurations :many
SELECT task_id, SUM(duration_seconds) as total_duration
FROM sessions
WHERE DATE(start_time) = CURRENT_DATE
GROUP BY task_id;

-- name: GetTaskDurationForPeriod :one
SELECT SUM(duration_seconds) as total_duration
FROM sessions
WHERE task_id = ?
AND start_time >= ?
AND start_time <= ?;