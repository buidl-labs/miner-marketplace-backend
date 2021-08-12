-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE miners
ADD COLUMN onboarded BOOLEAN;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE miners
DROP COLUMN onboarded;
