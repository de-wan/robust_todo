-- name: AddTodo :exec
INSERT INTO todo (uuid, value)
    VALUES (?, ?);

-- name: ListTodos :many
SELECT uuid, value, done_at
    FROM todo
    WHERE archived_at IS NULL
    ORDER BY created_at DESC;

-- name: GetTodo :one
SELECT uuid, value, done_at FROM todo WHERE uuid = ?;

-- name: ToggleTodo :exec
UPDATE todo SET
    done_at =
        CASE WHEN done_at IS NULL
            THEN NOW()
            ELSE NULL
        END
WHERE uuid = ?;

-- name: EditTodo :exec
UPDATE todo SET value = ? WHERE uuid = ?;

-- name: ArchiveTodo :exec
UPDATE todo SET archived_at = NOW() WHERE uuid = ?;

-- name: ListArchivedTodos :many
SELECT uuid, value, done_at
    FROM todo
    WHERE archived_at IS NOT NULL
    ORDER BY created_at DESC;

-- name: RestoreTodo :exec
UPDATE todo SET archived_at = NULL WHERE uuid = ?;