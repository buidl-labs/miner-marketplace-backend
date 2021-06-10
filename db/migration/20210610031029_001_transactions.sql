-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE IF NOT EXISTS transactions
(
    id TEXT NOT NULL,
    miner_id TEXT NOT NULL,
    height BIGINT NOT NULL,
    transaction_type TEXT NOT NULL,
    method_name TEXT NOT NULL,
    "value" TEXT NOT NULL,
    miner_fee TEXT NOT NULL,
    burn_fee TEXT NOT NULL,
    "from" TEXT NOT NULL,
    "to" TEXT NOT NULL,
    exit_code INTEGER NOT NULL,
    deals INTEGER[],
    PRIMARY KEY(id, miner_id, height)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE transactions;
