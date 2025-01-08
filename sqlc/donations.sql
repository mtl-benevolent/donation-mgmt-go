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

-- name: GetDonationBySlug :many
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
WHERE d.slug = sqlc.Arg('Slug')
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
	external_id, donation_id, amount_in_cents, receipt_amount_in_cents, received_at 
) VALUES(
	sqlc.Arg('ExternalID'), sqlc.Arg('DonationID'), sqlc.Arg('Amount'), sqlc.Arg('ReceiptAmount'), sqlc.Arg('ReceivedAt')
)
RETURNING *;

-- name: InsertPaymentToRecurrentDonation :one
INSERT INTO donation_payments(external_id, donation_id, amount_in_cents, receipt_amount_in_cents)
SELECT sqlc.Arg('ExternalID') as external_id, d.id, sqlc.Arg('AmountInCents') as amount, sqlc.Arg('ReceiptAmountInCents') as receipt_amount FROM donations d
WHERE d.type = 'RECURRENT'
	AND d.archived_at is null
	AND d.fiscal_year = sqlc.Arg('FiscalYear')
	AND d.organization_id = sqlc.Arg('OrganizationID')
	AND d.external_id = sqlc.Arg('ExternalID')
	AND d.environment = sqlc.Arg('Environment')
LIMIT 1
RETURNING id, donation_id;

-- name: UpdateDonationBySlug :execrows
UPDATE donations d 
SET donor_email = sqlc.arg('DonorEmail'), 
	donor_address = sqlc.arg('DonorAddress'),
	fiscal_year = sqlc.arg('FiscalYear')
where d.slug = sqlc.arg('Slug')
and d.environment = sqlc.arg('Environment')
and d.organization_id = sqlc.arg('OrganizationID')
and d.archived_at is null;
