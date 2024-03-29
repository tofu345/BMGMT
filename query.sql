-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :many
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY email;

-- name: CreateUser :one
INSERT INTO users (email, first_name, last_name, PASSWORD, is_superuser)
    VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateUsers :one
UPDATE users SET
    email = $2,
    first_name = $3,
    last_name = $4
WHERE id = $1 RETURNING *;

-- name: GetLocations :many
SELECT * FROM location;

-- name: GetLocationRooms :many 
SELECT room.* 
FROM location
    JOIN room ON room.location_id = location.id
WHERE location.id = $1;

-- name: CreateLocation :one
INSERT INTO location (name, address) VALUES ($1, $2)
RETURNING *;
