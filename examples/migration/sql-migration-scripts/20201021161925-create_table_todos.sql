-- +migrate Up
-- create_table_todos
CREATE TABLE IF NOT EXISTS todos (
    id SERIAL PRIMARY KEY
);

-- +migrate Down
-- create_table_todos
DROP TABLE IF EXISTS todos;