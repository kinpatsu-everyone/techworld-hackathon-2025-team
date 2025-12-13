-- name: GetMonsterAttribute :one
SELECT * FROM MonsterAttribute
WHERE MonsterId = ? LIMIT 1;

-- name: CreateMonsterAttribute :execresult
INSERT INTO MonsterAttribute (MonsterId, AttributeName, ColorCode)
VALUES (?, ?, ?);

-- name: UpdateMonsterAttribute :execresult
UPDATE MonsterAttribute
SET AttributeName = ?, ColorCode = ?
WHERE MonsterId = ?;

-- name: DeleteMonsterAttribute :exec
DELETE FROM MonsterAttribute
WHERE MonsterId = ?;
