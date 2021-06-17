-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE IF NOT EXISTS miner_storage_deal_stats
(
    id TEXT NOT NULL,
    average_price TEXT,
    data_stored TEXT,
    fault_terminated INTEGER,
    no_penalties INTEGER,
    slashed INTEGER,
    success_rate TEXT,
    terminated INTEGER,
    total INTEGER,
    PRIMARY KEY (id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE miner_storage_deal_stats;
