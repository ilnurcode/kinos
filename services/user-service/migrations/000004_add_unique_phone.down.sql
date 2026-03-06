-- Remove unique constraint on phone
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_phone_unique;
