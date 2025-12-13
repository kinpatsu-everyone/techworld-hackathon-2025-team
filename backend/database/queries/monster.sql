-- name: GetMonster :one
SELECT * FROM Monster
WHERE MonsterId = ? LIMIT 1;

-- name: GetMonsterDetail :one
SELECT
    m.MonsterId,
    m.Nickname,
    m.OriginalTrashBinImageUrl,
    m.GeneratedMonsterImageUrl,
    m.Latitude,
    m.Longitude,
    m.CreatedAt,
    m.UpdatedAt,
    ma.AttributeName,
    ma.ColorCode
FROM Monster m
LEFT JOIN MonsterAttribute ma ON m.MonsterId = ma.MonsterId
WHERE m.MonsterId = ? LIMIT 1;

-- name: ListMonsters :many
SELECT * FROM Monster
ORDER BY CreatedAt DESC;

-- name: ListMonstersWithAttribute :many
SELECT
    m.MonsterId,
    m.Nickname,
    m.OriginalTrashBinImageUrl,
    m.GeneratedMonsterImageUrl,
    m.Latitude,
    m.Longitude,
    m.CreatedAt,
    m.UpdatedAt,
    ma.AttributeName,
    ma.ColorCode
FROM Monster m
LEFT JOIN MonsterAttribute ma ON m.MonsterId = ma.MonsterId
ORDER BY m.CreatedAt DESC;

-- name: CreateMonster :execresult
INSERT INTO Monster (MonsterId, Nickname, OriginalTrashBinImageUrl, GeneratedMonsterImageUrl, Latitude, Longitude)
VALUES (?, ?, ?, ?, ?, ?);

-- name: UpdateMonster :execresult
UPDATE Monster
SET Nickname = ?, OriginalTrashBinImageUrl = ?, GeneratedMonsterImageUrl = ?, Latitude = ?, Longitude = ?
WHERE MonsterId = ?;

-- name: DeleteMonster :exec
DELETE FROM Monster
WHERE MonsterId = ?;
