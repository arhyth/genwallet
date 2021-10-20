-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE transfers
ADD COLUMN currency text NOT NULL;
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE transfers
DROP COLUMN IF EXISTS currency;