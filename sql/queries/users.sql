-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid(), -- id
    NOW(),
    NOW(),
    $1 -- email from first parameter from query
)
RETURNING *;

-- name: ResetUsers :one
DELETE FROM users
RETURNING *;