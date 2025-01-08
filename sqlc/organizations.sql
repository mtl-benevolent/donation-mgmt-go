
-- name: ListOrganizations :many
SELECT * FROM organizations
WHERE archived_at IS NULL
OFFSET sqlc.arg('Offset')
LIMIT sqlc.arg('Limit');

-- name: ListAuthorizedOrganizations :many
SELECT o.* FROM organizations o
INNER JOIN scoped_user_roles sur ON sur.organization_id = o.id
WHERE sur.subject = sqlc.arg('Subject')
	AND o.archived_at IS NULL
OFFSET sqlc.arg('Offset')
LIMIT sqlc.arg('Limit');

-- name: CountOrganizations :one
SELECT COUNT(*) AS total FROM organizations
WHERE archived_at IS NULL;

-- name: CountAuthorizedOrganizations :one
SELECT COUNT(*) FROM organizations o
INNER JOIN scoped_user_roles sur ON sur.organization_id = o.id
WHERE sur.subject = sqlc.arg('Subject')
	AND o.archived_at IS NULL;

-- name: GetOrganizationByID :one
SELECT * from organizations
WHERE id = sqlc.arg('OrganizationID')
	AND archived_at IS NULL;

-- name: GetOrganizationBySlug :one
SELECT * from organizations
WHERE slug = LOWER(sqlc.arg('Slug'))
  AND archived_at IS NULL;

-- name: GetOrganizationIDBySlug :one
SELECT id from organizations
WHERE slug = LOWER(sqlc.arg('Slug'))
	AND archived_at IS NULL;

-- name: InsertOrganization :one
INSERT INTO organizations(name, slug, timezone)
VALUES(sqlc.arg('Name'), LOWER(sqlc.arg('Slug')), sqlc.arg('TimeZone'))
RETURNING *;

-- name: UpdateOrganizationBySlug :one
UPDATE organizations
SET name = sqlc.arg('Name'), timezone = sqlc.arg('TimeZone')
WHERE slug = LOWER(sqlc.arg('Slug'))
RETURNING *;

-- name: ListOrganizationFiscalYears :many
SELECT DISTINCT fiscal_year FROM donations d
WHERE d.organization_id = sqlc.arg('OrganizationID')
	AND d.archived_at IS NULL
	AND d.environment = sqlc.Arg('Environment')
ORDER BY fiscal_year DESC;
