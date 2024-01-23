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
	and r.capabilities @> sqlc.arg('Capabilities')::varchar[]
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
	and r.capabilities @> sqlc.arg('Capabilities')::varchar[]
	and (sur.id is not null or gur.id is not null)
limit 1;

-- name: HasGlobalCapabilities :one
select r.id, r.name, r.created_at
from "roles" r
inner join global_user_roles gur on r.id = gur.role_id and gur.subject = sqlc.arg('Subject')
where r.archived_at is null
	and r.capabilities @> sqlc.arg('Capabilities')::varchar[]
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

-- name: GetDonationByID :many
WITH comments_count AS (
	SELECT count(*) AS "comments_count", dc.donation_id FROM donation_comments dc 
	WHERE dc.archived_at IS NULL
		AND dc.donation_id = sqlc.Arg('ID')
	GROUP BY donation_id
)
SELECT d.*, coalesce(cc.comments_count, 0) AS "comments_count", dp.* FROM donations d
INNER JOIN donation_payments dp
	ON dp.donation_id = d.id 
LEFT OUTER JOIN comments_count cc
	ON cc.donation_id = d.id
WHERE d.id = sqlc.Arg('ID')
	AND d.archived_at IS NULL
	AND d.organization_id = sqlc.Arg('OrganizationID')
	AND d.environment = sqlc.Arg('Environment');

-- name: InsertDonation :one
INSERT INTO donations(
	slug, organization_id, external_id, environment, fiscal_year, reason, type, source, 
	donor_firstname, "donor_lastname_or_orgName", donor_email, donor_address, emit_receipt, send_by_email
) VALUES (
	sqlc.Arg('Slug'), sqlc.Arg('OrganizationID'), sqlc.Arg('ExternalID'), sqlc.Arg('Environment'), 
	sqlc.Arg('FiscalYear'), sqlc.Arg('Reason'), sqlc.Arg('Type'), sqlc.Arg('Source'),
	sqlc.Arg('DonorFirstname'), sqlc.Arg('DonorLastNameOrOrgName'), sqlc.Arg('DonorEmail'), sqlc.Arg('DonorAddress'),
	sqlc.Arg('EmitReceipt'), sqlc.Arg('SendByEmail')
)
RETURNING *;

-- name: InsertDonationPayment :one
INSERT INTO donation_payments(
	external_id, donation_id, amount, receipt_amount, received_at 
) VALUES(
	sqlc.Arg('ExternalID'), sqlc.Arg('DonationID'), sqlc.Arg('Amount'), sqlc.Arg('ReceiptAmount'), sqlc.Arg('ReceivedAt')
)
RETURNING *;

-- name: InsertPaymentToRecurrentDonation :one
INSERT INTO donation_payments(external_id, donation_id, amount, receipt_amount)
SELECT sqlc.Arg('ExternalID') as external_id, d.id, sqlc.Arg('Amount') as amount, sqlc.Arg('ReceiptAmount') as receipt_amount FROM donations d
WHERE d.type = 'RECURRENT'
	AND d.archived_at is null
	AND d.fiscal_year = sqlc.Arg('FiscalYear')
	AND d.organization_id = sqlc.Arg('OrganizationID')
	AND d.external_id = sqlc.Arg('ExternalID')
	AND d.environment = sqlc.Arg('Environment')
LIMIT 1
RETURNING id, donation_id;

-- name: UpdateDonationBySlug :exec
UPDATE donations d 
SET donor_email = sqlc.arg('DonorEmail'), 
	donor_address = sqlc.arg('DonorAddress'),
	fiscal_year = sqlc.arg('FiscalYear')
where d.slug = sqlc.arg('Slug')
and d.environment = sqlc.arg('Environment')
and d.organization_id = sqlc.arg('OrganizationID')
and d.archived_at is null;
