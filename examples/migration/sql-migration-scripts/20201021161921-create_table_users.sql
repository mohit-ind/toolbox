-- +migrate Up
-- create_table_users
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY
);

-- +migrate Down
-- create_table_users
DROP TABLE IF EXISTS users;