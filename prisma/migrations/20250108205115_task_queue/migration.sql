-- CreateEnum
CREATE TYPE "TaskType" AS ENUM ('GENERATE_RECEIPT');

-- CreateEnum
CREATE TYPE "TaskStatus" AS ENUM ('CREATED', 'IN_PROGRESS', 'COMPLETED', 'ERROR_RETRYABLE', 'ERROR_UNRETRYABLE');

-- CreateTable
CREATE TABLE "task_queue" (
    "id" BIGSERIAL NOT NULL,
    "body" JSONB,
    "type" "TaskType" NOT NULL,
    "status" "TaskStatus" NOT NULL DEFAULT 'CREATED',
    "comment" TEXT,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "completed_at" TIMESTAMPTZ,
    "last_picked_up_at" TIMESTAMP(3),
    "locked_until" TIMESTAMPTZ,
    "locked_by" TEXT NOT NULL,
    "max_retries" INTEGER NOT NULL,
    "retry_count" INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT "task_queue_pkey" PRIMARY KEY ("id")
);

-- CreateIndex
CREATE INDEX "task_queue_locked_until_idx" ON "task_queue"("locked_until");
