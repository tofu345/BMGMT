-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT email, first_name, last_name, is_superuser FROM users
ORDER BY email;

-- name: CreateUser :one
INSERT INTO users (email, first_name, last_name, PASSWORD, is_superuser)
    VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetUserLocAdmins :many
SELECT location.*
FROM users
    JOIN loc_admin ON loc_admin.user_id = users.id
    JOIN location ON location.id = loc_admin.location_id
WHERE users.id = $1;

-- -- name: UpdateUsers :one
-- UPDATE users SET
--     email = $2,
--     first_name = $3,
--     last_name = $4
-- WHERE id = $1 RETURNING *;

-- name: GetLocations :many
SELECT * FROM location;

-- name: GetLocation :one
SELECT * FROM location WHERE id = $1 LIMIT 1;

-- name: GetLocationAdmins :many 
SELECT users.email 
FROM loc_admin 
    JOIN users ON loc_admin.user_id = users.id
WHERE loc_admin.location_id = $1;

-- name: GetLocationRooms :many 
SELECT room.name, users.email, users.first_name, users.last_name
FROM location
    JOIN room ON room.location_id = location.id
    LEFT JOIN users ON room.tenant_id = users.id
WHERE location.id = $1;

-- name: CreateLocation :one
INSERT INTO location (name, address) VALUES ($1, $2)
RETURNING *;

-- name: CreateRoom :one
INSERT INTO room (name, location_id) VALUES ($1, $2)
RETURNING *;

-- name: CreateLocationAdmin :one
INSERT INTO loc_admin (user_id, location_id) VALUES ($1, $2)
RETURNING *;
