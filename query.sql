-- name: AddTodo :exec
INSERT INTO todo (uuid, value)
    VALUES (?, ?);

-- name: ListTodos :many
SELECT uuid, value, deadline, done_at, created_at
    FROM todo
    ORDER BY created_at DESC;