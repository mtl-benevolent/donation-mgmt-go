-- CreateTable
CREATE TABLE "organization_settings" (
    "organization_id" BIGINT NOT NULL,
    "environment" "Environment" NOT NULL,
    "timezone" TEXT NOT NULL DEFAULT 'America/Toronto',
    "email_provider_settings" TEXT NOT NULL,
    "is_valid" BOOLEAN NOT NULL DEFAULT false,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "organization_settings_pkey" PRIMARY KEY ("organization_id","environment")
);

-- CreateTable
CREATE TABLE "organization_templates" (
    "organization_id" BIGINT NOT NULL,
    "environment" "Environment" NOT NULL,
    "is_valid" BOOLEAN NOT NULL DEFAULT false,
    "no_mailing_addr_email" TEXT,
    "no_mail_addr_email_title" TEXT,
    "no_mailing_addr_reminder_email" TEXT,
    "no_mailing_addr_reminder_email_title" TEXT,
    "receipt_pdf" TEXT,
    "receipt_email" TEXT,
    "receipt_email_title" TEXT,
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "organization_templates_pkey" PRIMARY KEY ("organization_id","environment")
);

-- AddForeignKey
ALTER TABLE "organization_settings" ADD CONSTRAINT "organization_settings_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organizations"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "organization_templates" ADD CONSTRAINT "organization_templates_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organizations"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
