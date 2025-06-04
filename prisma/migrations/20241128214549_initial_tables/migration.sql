-- CreateEnum
CREATE TYPE "Environment" AS ENUM ('SANDBOX', 'LIVE');

-- CreateEnum
CREATE TYPE "DonationType" AS ENUM ('ONE_TIME', 'RECURRENT');

-- CreateEnum
CREATE TYPE "DonationSource" AS ENUM ('PAYPAL', 'CHEQUE', 'DIRECT_DEPOSIT', 'STOCKS', 'OTHER');

-- CreateTable
CREATE TABLE "organizations" (
    "id" BIGSERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "slug" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "archived_at" TIMESTAMPTZ,

    CONSTRAINT "organizations_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "roles" (
    "id" BIGSERIAL NOT NULL,
    "name" TEXT NOT NULL,
    "capabilities" TEXT[],
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "archived_at" TIMESTAMPTZ,

    CONSTRAINT "roles_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "scoped_user_roles" (
    "id" BIGSERIAL NOT NULL,
    "subject" TEXT NOT NULL,
    "role_id" BIGINT NOT NULL,
    "organization_id" BIGINT NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "scoped_user_roles_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "global_user_roles" (
    "id" BIGSERIAL NOT NULL,
    "subject" TEXT NOT NULL,
    "role_id" BIGINT NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT "global_user_roles_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "donations" (
    "id" BIGSERIAL NOT NULL,
    "slug" TEXT NOT NULL,
    "organization_id" BIGINT NOT NULL,
    "external_id" TEXT,
    "environment" "Environment" NOT NULL,
    "fiscal_year" SMALLINT NOT NULL,
    "reason" TEXT,
    "type" "DonationType" NOT NULL,
    "source" "DonationSource" NOT NULL,
    "donor_firstname" TEXT,
    "donor_lastname_or_orgName" TEXT NOT NULL,
    "donor_email" TEXT,
    "donor_address" JSONB,
    "emit_receipt" BOOLEAN NOT NULL,
    "send_by_email" BOOLEAN NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ,
    "archived_at" TIMESTAMPTZ,

    CONSTRAINT "donations_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "donation_payments" (
    "id" BIGSERIAL NOT NULL,
    "external_id" TEXT,
    "donation_id" BIGINT NOT NULL,
    "amount_in_cents" BIGINT NOT NULL,
    "receipt_amount_in_cents" BIGINT NOT NULL,
    "received_at" TIMESTAMPTZ NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "archived_at" TIMESTAMPTZ,

    CONSTRAINT "donation_payments_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "donation_comments" (
    "id" BIGSERIAL NOT NULL,
    "comment" TEXT NOT NULL,
    "author" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "archived_at" TIMESTAMPTZ,
    "donation_id" BIGINT NOT NULL,

    CONSTRAINT "donation_comments_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "organizations_slug_key" ON "organizations"("slug");

-- CreateIndex
CREATE UNIQUE INDEX "roles_name_key" ON "roles"("name");

-- CreateIndex
CREATE UNIQUE INDEX "scoped_user_roles_subject_organization_id_key" ON "scoped_user_roles"("subject", "organization_id");

-- CreateIndex
CREATE UNIQUE INDEX "global_user_roles_subject_key" ON "global_user_roles"("subject");

-- CreateIndex
CREATE UNIQUE INDEX "donations_slug_key" ON "donations"("slug");

-- CreateIndex
CREATE INDEX "donations_organization_id_environment_fiscal_year_idx" ON "donations"("organization_id", "environment", "fiscal_year");

-- CreateIndex
CREATE INDEX "donations_organization_id_environment_fiscal_year_external__idx" ON "donations"("organization_id", "environment", "fiscal_year", "external_id", "source");

-- CreateIndex
CREATE UNIQUE INDEX "donations_organization_id_environment_slug_key" ON "donations"("organization_id", "environment", "slug");

-- AddForeignKey
ALTER TABLE "scoped_user_roles" ADD CONSTRAINT "scoped_user_roles_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "roles"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "scoped_user_roles" ADD CONSTRAINT "scoped_user_roles_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organizations"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "global_user_roles" ADD CONSTRAINT "global_user_roles_role_id_fkey" FOREIGN KEY ("role_id") REFERENCES "roles"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "donations" ADD CONSTRAINT "donations_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organizations"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "donation_payments" ADD CONSTRAINT "donation_payments_donation_id_fkey" FOREIGN KEY ("donation_id") REFERENCES "donations"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "donation_comments" ADD CONSTRAINT "donation_comments_donation_id_fkey" FOREIGN KEY ("donation_id") REFERENCES "donations"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
