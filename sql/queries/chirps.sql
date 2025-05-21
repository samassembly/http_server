-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(), -- id
    NOW(),
    NOW(),
    $1, -- body from first parameter of query
    $2 -- user_id from second parameter of query
)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: GetChirp :one
SELECT * from chirps
WHERE id = $1;


-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;
