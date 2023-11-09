-- name: GetOrganizationByID :one
SELECT * FROM organizations
WHERE id = sqlc.arg('OrganizationID')
  AND archived_at IS NULL;

-- name: InsertOrganization :one
INSERT INTO organizations(name, slug)
VALUES(sqlc.arg('Name'), sqlc.arg('Slug'))
RETURNING *;
