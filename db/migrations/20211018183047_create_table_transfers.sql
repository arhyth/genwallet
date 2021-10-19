-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS transfers (
    id serial PRIMARY KEY,
    "from" text REFERENCES accounts (id),
    "to" text REFERENCES accounts (id),
    -- validation (greater than 0) delegated to service
    amount real,
    created_at timestamp with time zone DEFAULT now()
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE IF EXISTS transfers
