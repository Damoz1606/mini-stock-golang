-- SQLC query definitions for category CRUD operations.
-- Column names follow the PostgreSQL schema: id, name, "createdAt", "updatedAt".
-- All queries use parameterized placeholders ($1, $2, ...) for SQL injection safety.

-- name: CreateCategory :one
INSERT INTO categories (id, name, "createdAt", "updatedAt")
VALUES ($1, $2, $3, $4)
RETURNING id, name, "createdAt", "updatedAt";

-- name: FindCategoryByID :one
SELECT id, name, "createdAt", "updatedAt"
FROM categories
WHERE id = $1;

-- name: FindAllCategories :many
SELECT id, name, "createdAt", "updatedAt"
FROM categories
WHERE ($1::text IS NULL OR name ILIKE '%' || $1 || '%')
ORDER BY
  CASE WHEN $2 = 'name' AND $3 = 'asc' THEN name END ASC,
  CASE WHEN $2 = 'name' AND $3 = 'desc' THEN name END DESC,
  CASE WHEN $2 = 'createdAt' AND $3 = 'asc' THEN "createdAt" END ASC,
  CASE WHEN $2 = 'createdAt' AND $3 = 'desc' THEN "createdAt" END DESC
LIMIT $4 OFFSET $5;

-- name: CountCategories :one
SELECT COUNT(*) AS count FROM categories
WHERE ($1::text IS NULL OR name ILIKE '%' || $1 || '%');

-- name: UpdateCategory :one
UPDATE categories
SET name = $1, "updatedAt" = $2
WHERE id = $3
RETURNING id, name, "createdAt", "updatedAt";

-- name: DeleteCategory :exec
DELETE FROM categories WHERE id = $1;

-- name: ExistsCategoryByName :one
SELECT EXISTS(SELECT 1 FROM categories WHERE name = $1 AND id != $2) AS "exists";
