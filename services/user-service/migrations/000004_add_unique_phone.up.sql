-- Add unique constraint on phone
ALTER TABLE users ADD CONSTRAINT users_phone_unique UNIQUE (phone);
