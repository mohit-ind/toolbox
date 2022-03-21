-- +migrate Up
-- add_email_to_users
ALTER TABLE users ADD COLUMN email VARCHAR(64) NOT NULL;

-- +migrate Down
-- add_email_to_users
ALTER TABLE users DROP COLUMN email;
