-- name: GetOrganizations :many
SELECT * from organizations
WHERE archived_at IS NULL;

-- name: GetOrganizationBySlug :one
SELECT * from organizations
WHERE slug = LOWER(sqlc.arg('Slug'))
  AND archived_at IS NULL;

-- name: InsertOrganization :one
INSERT INTO organizations(name, slug)
VALUES(sqlc.arg('Name'), LOWER(sqlc.arg('Slug')))
RETURNING *;
