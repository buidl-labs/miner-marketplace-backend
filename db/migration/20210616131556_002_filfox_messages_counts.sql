-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE IF NOT EXISTS filfox_miner_messages_counts
(
    id TEXT NOT NULL,
    miner_messages_total_count BIGINT,
    miner_transfers_reward_total_count BIGINT,
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS filfox_messages_counts
(
    id TEXT NOT NULL,
    publish_storage_deals_messages_total_count BIGINT,
    withdraw_balance_market_messages_total_count BIGINT,
    add_balance_messages_total_count BIGINT,
    PRIMARY KEY(id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE filfox_miner_messages_counts;
DROP TABLE filfox_messages_counts;
