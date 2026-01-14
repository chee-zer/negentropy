-- name: StartSession :one
INSERT INTO sessions (start_time, task_id)
VALUES (?, ?)
RETURNING *;

-- name: EndSession :one
UPDATE sessions
SET end_time = ?
WHERE task_id = ?
RETURNING *;

-- name: EndSessionAsEntropy :one
UPDATE sessions
SET
end_time = ?,
task_id = 0

WHERE task_id = ?
RETURNING *;

-- name: GetDailyTaskDurations :many
SELECT task_id, SUM(duration_seconds) as total_seconds
FROM (
    -- the session that spans two days will be stored in db as a single session, but will be divided between 2 days inside the app
    -- also this is for anyday, not just today. will make wrappers for this query inside app instead.
    -- '?' arg is the queried date
    -- 1. sessions spanning single day, day x
    SELECT s.task_id,
    strftime('%s', s.end_time) - strftime('%s', s.start_time) AS duration_seconds
    FROM sessions AS s
    WHERE s.start_time >= sqlc.arg(query_date)
    AND s.end_time < strftime('%s', sqlc.arg(query_date), '+1 day', 'start of day')
    UNION ALL
    -- 2. sessions spanning two days, started on day x (not accounting for sessions spanning more than 2 days)
    SELECT s.task_id,
    strftime('%s', sqlc.arg(query_date), '+1 day', 'start of day') - strftime('%s', start_time) AS duration_seconds
    FROM sessions AS s
    WHERE date(s.start_time) = date(sqlc.arg(query_date))
    AND date(s.end_time) = date(sqlc.arg(query_date), '+1 day')
    OR s.end_time IS NULL
    UNION ALL

    --3. sessions spanning two days, ending on day x
    SELECT s.task_id,
    strftime('%s', s.end_time) - strftime('%s', sqlc.arg(query_date), 'start of day') AS duration_seconds
    FROM sessions AS s
    WHERE date(s.start_time) = date(sqlc.arg(query_date), '-1 day')
    AND date(s.end_time) = date(sqlc.arg(query_date))
) AS daily_sessions -- https://github.com/sqlc-dev/sqlc/issues/3963
GROUP BY task_id

-- TODO: Get weekly, monthly, yearly, total duration
