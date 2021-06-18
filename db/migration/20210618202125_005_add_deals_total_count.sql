-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE filfox_messages_counts
ADD COLUMN deals_total_count BIGINT;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE filfox_messages_counts
DROP COLUMN deals_total_count;
