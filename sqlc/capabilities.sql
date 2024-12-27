
-- name: HasCapabilitiesForOrgBySlug :one
with sur as (
	select ur.id, ur.organization_id, ur.role_id
	from scoped_user_roles ur 
	inner join organizations o on o.id = ur.organization_id
	where o.slug = lower(sqlc.arg('OrganizationSlug'))
		and o.archived_at is null
		and ur.subject = sqlc.arg('Subject')
)
select r.id, r.name, r.created_at, sur.organization_id, (sur.id is not null)::boolean as is_user_role, (gur.id is not null)::boolean as is_global_role
from "roles" r
left outer join sur on r.id = sur.role_id
left outer join global_user_roles gur on r.id = gur.role_id and gur.subject = sqlc.arg('Subject')
where r.archived_at is null
	and r.capabilities @> sqlc.arg('Capabilities')::text[]
	and (sur.id is not null or gur.id is not null)
limit 1;

-- name: HasCapabilitiesForOrgByID :one
with sur as (
	select ur.id, ur.organization_id, ur.role_id
	from scoped_user_roles ur 
	inner join organizations o on o.id = ur.organization_id
	where o.id = sqlc.arg('OrganizationID')
		and o.archived_at is null
		and ur.subject = sqlc.arg('Subject')
)
select r.id, r.name, r.created_at, sur.organization_id, (sur.id is not null)::boolean as is_user_role, (gur.id is not null)::boolean as is_global_role
from "roles" r
left outer join sur on r.id = sur.role_id
left outer join global_user_roles gur on r.id = gur.role_id and gur.subject = sqlc.arg('Subject')
where r.archived_at is null
	and r.capabilities @> sqlc.arg('Capabilities')::text[]
	and (sur.id is not null or gur.id is not null)
limit 1;

-- name: HasGlobalCapabilities :one
select r.id, r.name, r.created_at
from "roles" r
inner join global_user_roles gur on r.id = gur.role_id and gur.subject = sqlc.arg('Subject')
where r.archived_at is null
	and r.capabilities @> sqlc.arg('Capabilities')::text[]
limit 1;

-- name: GetScopedRoles :many
select r.id, r.name, r.capabilities, ur.created_at as granted_on
from roles r
inner join scoped_user_roles ur on r.id = ur.role_id
where ur.subject = sqlc.arg('Subject')
	and ur.organization_id  = sqlc.arg('OrganizationID')
	and r.archived_at is null;

-- name: GrantScopedRole :one
INSERT INTO scoped_user_roles(subject, role_id, organization_id)
VALUES(sqlc.arg('Subject'), sqlc.arg('RoleID'), sqlc.arg('OrganizationID'))
RETURNING *;

-- name: RevokeScopedRoles :exec
DELETE FROM scoped_user_roles
WHERE subject = sqlc.arg('Subject')
  AND organization_id = sqlc.arg('OrganizationID');
