-- name: GetUser :one
SELECT * FROM User
WHERE UserId = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM User
ORDER BY CreatedAt DESC;

-- name: CreateUser :execresult
INSERT INTO User (UserId, Nickname)
VALUES (?, ?);

-- name: UpdateUser :execresult
UPDATE User
SET Nickname = ?
WHERE UserId = ?;

-- name: DeleteUser :exec
DELETE FROM User
WHERE UserId = ?;
