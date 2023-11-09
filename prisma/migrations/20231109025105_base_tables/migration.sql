-- CreateEnum
CREATE TYPE "Enviroment" AS ENUM ('SANDBOX', 'LIVE');

-- CreateEnum
CREATE TYPE "DonationType" AS ENUM ('ONE_TIME', 'RECURRENT');

-- CreateEnum
CREATE TYPE "DonationSource" AS ENUM ('PAYPAL', 'CHEQUE', 'DIRECT_DEPOSIT', 'STOCKS', 'OTHER');

-- CreateTable
CREATE TABLE "organizations" (
    "id" UUID NOT NULL,
    "name" STRING NOT NULL,
    "slug" STRING NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "archived_at" TIMESTAMPTZ,

    CONSTRAINT "organizations_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "donations" (
    "id" UUID NOT NULL,
    "organization_id" UUID NOT NULL,
    "external_id" STRING,
    "environment" "Enviroment" NOT NULL,
    "fiscal_year" INT2 NOT NULL,
    "reason" STRING,
    "type" "DonationType" NOT NULL,
    "source" "DonationSource" NOT NULL,
    "donor_firstname" STRING,
    "donor_lastname_or_orgName" STRING NOT NULL,
    "donor_email" STRING,
    "donor_address" JSONB,
    "emit_receipt" BOOL NOT NULL,
    "send_by_email" BOOL NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMPTZ,
    "archived_at" TIMESTAMPTZ,

    CONSTRAINT "donations_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "donation_payments" (
    "id" UUID NOT NULL,
    "external_id" STRING,
    "donation_id" UUID NOT NULL,
    "amount" INT8 NOT NULL CHECK(amount >= 0),
    "receipt_amount" INT8 NOT NULL CHECK(receipt_amount >= 0),
    "received_at" TIMESTAMPTZ NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "archived_at" TIMESTAMPTZ,

    CONSTRAINT "donation_payments_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "ReceiptAmount_LessThanOrEq_To_Amount" CHECK(receipt_amount <= amount)
);

-- CreateTable
CREATE TABLE "donation_events" (
    "id" UUID NOT NULL,
    "event_name" STRING NOT NULL,
    "details" JSONB,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "donation_id" UUID NOT NULL,

    CONSTRAINT "donation_events_pkey" PRIMARY KEY ("id")
);

-- CreateTable
CREATE TABLE "donation_comments" (
    "id" UUID NOT NULL,
    "comment" STRING NOT NULL,
    "author" STRING NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "archived_at" TIMESTAMPTZ,
    "donation_id" UUID NOT NULL,

    CONSTRAINT "donation_comments_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE UNIQUE INDEX "organizations_slug_key" ON "organizations"("slug");

-- CreateIndex
CREATE INDEX "donations_organization_id_environment_fiscal_year_idx" ON "donations"("organization_id", "environment", "fiscal_year");

-- AddForeignKey
ALTER TABLE "donations" ADD CONSTRAINT "donations_organization_id_fkey" FOREIGN KEY ("organization_id") REFERENCES "organizations"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "donation_payments" ADD CONSTRAINT "donation_payments_donation_id_fkey" FOREIGN KEY ("donation_id") REFERENCES "donations"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "donation_events" ADD CONSTRAINT "donation_events_donation_id_fkey" FOREIGN KEY ("donation_id") REFERENCES "donations"("id") ON DELETE RESTRICT ON UPDATE CASCADE;

-- AddForeignKey
ALTER TABLE "donation_comments" ADD CONSTRAINT "donation_comments_donation_id_fkey" FOREIGN KEY ("donation_id") REFERENCES "donations"("id") ON DELETE RESTRICT ON UPDATE CASCADE;
