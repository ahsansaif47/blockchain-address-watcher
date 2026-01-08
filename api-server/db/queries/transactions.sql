-- name: CreateUser :one
INSERT INTO users (
    id,
    email,
    password_hash,
    phone_number,
    wallet_address,
    subscribed,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, NOW(), NOW()
)
RETURNING
    id;

-- name: SignInUser :one
SELECT
    id,
    email,
    password_hash,
    phone_number,
    wallet_address,
    subscribed,
    created_at,
    updated_at,
    deleted_at
FROM users
WHERE email = $1 AND deleted_at IS NULL;

-- name: SoftDeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: HardDeleteUser :exec
DELETE FROM users
WHERE id = $1;