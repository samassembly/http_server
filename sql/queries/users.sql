-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(), -- id
    NOW(),
    NOW(),
    $1, -- email from first parameter from query
    $2 -- hash from second parameter from query
)
RETURNING *;

-- name: ResetUsers :one
DELETE FROM users
RETURNING *;

-- name: LoginUser :one
SELECT * FROM users
WHERE email = $1;

-- name: UpdateUser :one
UPDATE users SET email = $2, hashed_password = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;