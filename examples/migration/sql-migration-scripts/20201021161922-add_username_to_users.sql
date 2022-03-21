-- +migrate Up
-- add_username_to_users
ALTER TABLE users ADD COLUMN username VARCHAR(32) NOT NULL;

-- +migrate Down
-- add_username_to_users
ALTER TABLE users DROP COLUMN username;
