
-- name: GetOrganizationTemplates :one
SELECT * FROM organization_templates
WHERE organization_id = sqlc.arg('OrganizationID')
	AND environment = sqlc.arg('Environment')
LIMIT 1;

-- name: UpsertNoMailingAddrEmailTemplate :one
INSERT INTO organization_templates(
    organization_id,
    environment,
    no_mail_addr_email_title,
    no_mailing_addr_email,
    updated_at
) VALUES (
    sqlc.arg('OrganizationID'),
    sqlc.arg('Environment'),
    sqlc.narg('NoMailingAddrEmailTitle'),
    sqlc.narg('NoMailingAddrEmail'),
    NOW()
)
ON CONFLICT (organization_id, environment) DO UPDATE SET
    no_mail_addr_email_title = sqlc.narg('NoMailingAddrEmailTitle'),
    no_mailing_addr_email = sqlc.narg('NoMailingAddrEmail'),
    updated_at = NOW()
RETURNING *;

-- name: UpsertNoMailingAddrReminderEmailTemplate :one
INSERT INTO organization_templates(
    organization_id,
    environment,
    no_mailing_addr_reminder_email,
    no_mailing_addr_reminder_email_title,
    updated_at
) VALUES (
    sqlc.arg('OrganizationID'),
    sqlc.arg('Environment'),
    sqlc.narg('NoMailingAddrReminderEmail'),
    sqlc.narg('NoMailingAddrReminderEmailTitle'),
    NOW()
)
ON CONFLICT (organization_id, environment) DO UPDATE SET
    no_mailing_addr_reminder_email = sqlc.narg('NoMailingAddrReminderEmail'),
    no_mailing_addr_reminder_email_title = sqlc.narg('NoMailingAddrReminderEmailTitle'),
    updated_at = NOW()
RETURNING *;

-- name: UpsertReceiptPdfTemplate :one
INSERT INTO organization_templates(
    organization_id,
    environment,
    receipt_pdf,
    updated_at
) VALUES (
    sqlc.arg('OrganizationID'),
    sqlc.arg('Environment'),
    sqlc.narg('ReceiptPdf'),
    NOW()
)
ON CONFLICT (organization_id, environment) DO UPDATE SET
    receipt_pdf = sqlc.narg('ReceiptPdf'),
    updated_at = NOW()
RETURNING *;

-- name: UpsertReceiptEmailTemplate :one
INSERT INTO organization_templates(
    organization_id,
    environment,
    receipt_email,
    receipt_email_title,
    updated_at
) VALUES (
    sqlc.arg('OrganizationID'),
    sqlc.arg('Environment'),
    sqlc.narg('ReceiptEmail'),
    sqlc.narg('ReceiptEmailTitle'),
    NOW()
)
ON CONFLICT (organization_id, environment) DO UPDATE SET
    receipt_email = sqlc.narg('ReceiptEmail'),
    receipt_email_title = sqlc.narg('ReceiptEmailTitle'),
    updated_at = NOW()
RETURNING *;
