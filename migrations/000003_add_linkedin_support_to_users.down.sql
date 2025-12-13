-- Reverting the changes from the 'up' migration.

-- Before adding the NOT NULL constraint, you might need to handle any users
-- with a NULL password to avoid errors. For this rollback, we assume this is handled.
ALTER TABLE "users" ALTER COLUMN "password" SET NOT NULL;

-- Drop the 'linkedin_id' column.
ALTER TABLE "users" DROP COLUMN "linkedin_id";
