-- name: AddTodo :exec
INSERT INTO todo (uuid, value)
    VALUES (?, ?);