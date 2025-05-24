-- CreateEnum
CREATE TYPE "TaskType" AS ENUM ('GENERATE_RECEIPT');

-- CreateEnum
CREATE TYPE "TaskStatus" AS ENUM ('CREATED', 'IN_PROGRESS', 'COMPLETED', 'ERROR_RETRYABLE', 'ERROR_UNRETRYABLE');

-- CreateTable
CREATE TABLE "tasks" (
    "id" BIGSERIAL NOT NULL,
    "body" JSONB,
    "type" "TaskType" NOT NULL,
    "status" "TaskStatus" NOT NULL DEFAULT 'CREATED',
    "last_error_message" TEXT,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "completed_at" TIMESTAMPTZ,
    "last_picked_up_at" TIMESTAMP(3),
    "locked_until" TIMESTAMPTZ,
    "locked_by" TEXT,
    "max_retries" INTEGER NOT NULL,
    "attempt" INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT "tasks_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "tasks_type_status_created_at_attempt_locked_until_idx" ON "tasks"("type", "status", "created_at", "attempt", "locked_until");
