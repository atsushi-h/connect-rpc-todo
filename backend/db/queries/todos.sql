-- name: ListTodosByUser :many
SELECT * FROM todos WHERE user_id = $1 ORDER BY created_at DESC;
-- name: CreateTodo :one
INSERT INTO todos (user_id, title) VALUES ($1, $2) RETURNING *;
-- name: UpdateTodo :one
UPDATE todos SET title = $1, completed = $2, updated_at = now()
WHERE id = $3 AND user_id = $4 RETURNING *;
-- name: DeleteTodo :exec
DELETE FROM todos WHERE id = $1 AND user_id = $2;
