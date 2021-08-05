-- +goose Up
-- SQL in this section is executed when the migration is applied.

ALTER TABLE filfox_miner_messages_counts
ADD COLUMN till_height BIGINT;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

ALTER TABLE filfox_miner_messages_counts
DROP COLUMN till_height;
