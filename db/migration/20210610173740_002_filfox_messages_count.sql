-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE IF NOT EXISTS filfox_messages_counts
(
    id TEXT NOT NULL,
    miner_messages_total_count BIGINT,
    miner_transfers_reward_total_count BIGINT,
    market_actor_messages_total_count BIGINT,
    PRIMARY KEY(id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DELETE TABLE filfox_messages_counts;
