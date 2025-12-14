-- Add a nullable 'linkedin_id' column to store the user's unique ID from LinkedIn.
ALTER TABLE "users" ADD COLUMN "linkedin_id" TEXT;

-- Make the password column nullable to allow for social-only signups.
ALTER TABLE "users" ALTER COLUMN "password_hash" DROP NOT NULL;
