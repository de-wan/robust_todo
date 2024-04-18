-- name: AddTodo :exec
INSERT INTO todo (uuid, value)
    VALUES (?, ?);

-- name: ListTodos :many
SELECT uuid, value
    FROM todo
    ORDER BY created_at DESC;