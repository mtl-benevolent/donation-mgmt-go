
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

-- name: GetOrganizationWithSettings :one
SELECT 
	o.*, 
	COALESCE(os.timezone, 'America/Toronto') as timezone
FROM organizations o
LEFT OUTER JOIN organization_settings os
	ON os.organization_id = o.id
WHERE o.id = sqlc.arg('OrganizationID')
	AND o.archived_at IS NULL
	AND os.environment = sqlc.arg('Environment');

-- name: InsertOrganization :one
INSERT INTO organizations(name, slug)
VALUES(sqlc.arg('Name'), LOWER(sqlc.arg('Slug')))
RETURNING *;

-- name: ListOrganizationFiscalYears :many
SELECT DISTINCT fiscal_year FROM donations d
WHERE d.organization_id = sqlc.arg('OrganizationID')
	AND d.archived_at IS NULL
	AND d.environment = sqlc.Arg('Environment')
ORDER BY fiscal_year DESC;

-- name: UpsertOrganizationSettings :one
INSERT INTO organization_settings(organization_id, environment, timezone, email_provider, email_provider_settings)
VALUES(
	sqlc.arg('OrganizationID'),
	sqlc.arg('Environment'),
	sqlc.narg('Timezone'),
	sqlc.narg('EmailProvider'),
	COALESCE(sqlc.narg('EmailProviderSettings')::JSONB, '{}'::JSONB)
)
ON CONFLICT (organization_id, environment)
DO UPDATE
	SET timezone = COALESCE(sqlc.narg('Timezone'), EXCLUDED.timezone),
	email_provider = COALESCE(sqlc.narg('EmailProvider'), EXCLUDED.email_provider),
	email_provider_settings = COALESCE(sqlc.narg('EmailProviderSettings')::JSONB, EXCLUDED.email_provider_settings)
RETURNING *;
