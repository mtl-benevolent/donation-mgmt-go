-- name: GetOrganizations :many
SELECT * from organizations
WHERE archived_at IS NULL
OFFSET sqlc.arg('Offset')
LIMIT sqlc.arg('Limit');

-- name: GetOrganizationsCount :one
SELECT COUNT(*) AS total from organizations
WHERE archived_at IS NULL;

-- name: GetOrganizationBySlug :one
SELECT * from organizations
WHERE slug = LOWER(sqlc.arg('Slug'))
  AND archived_at IS NULL;

-- name: InsertOrganization :one
INSERT INTO organizations(name, slug)
VALUES(sqlc.arg('Name'), LOWER(sqlc.arg('Slug')))
RETURNING *;

-- name: UpdateOrganizationBySlug :one
UPDATE organizations
SET name = sqlc.arg('Name')
WHERE slug = LOWER(sqlc.arg('Slug'))
RETURNING *;
