/*
  Warnings:

  - You are about to drop the column `is_valid` on the `organization_settings` table. All the data in the column will be lost.
  - You are about to drop the column `is_valid` on the `organization_templates` table. All the data in the column will be lost.

*/
-- AlterTable
ALTER TABLE "organization_settings" DROP COLUMN "is_valid";

-- AlterTable
ALTER TABLE "organization_templates" DROP COLUMN "is_valid";
