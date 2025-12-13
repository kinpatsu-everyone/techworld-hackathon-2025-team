-- name: GetMonsterTrashCategory :one
SELECT * FROM MonsterTrashCategory
WHERE MonsterTrashCategoryId = ? LIMIT 1;

-- name: ListMonsterTrashCategories :many
SELECT * FROM MonsterTrashCategory
WHERE MonsterId = ?
ORDER BY TrashCategory;

-- name: ListMonstersByTrashCategory :many
SELECT * FROM MonsterTrashCategory
WHERE TrashCategory = ?
ORDER BY CreatedAt DESC;

-- name: CreateMonsterTrashCategory :execresult
INSERT INTO MonsterTrashCategory (MonsterTrashCategoryId, MonsterId, TrashCategory)
VALUES (?, ?, ?);

-- name: DeleteMonsterTrashCategory :exec
DELETE FROM MonsterTrashCategory
WHERE MonsterTrashCategoryId = ?;

-- name: DeleteMonsterTrashCategoriesByMonsterId :exec
DELETE FROM MonsterTrashCategory
WHERE MonsterId = ?;
