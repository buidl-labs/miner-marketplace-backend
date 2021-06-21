-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE IF NOT EXISTS market_deal_proposals
(
    id BIGINT NOT NULL,
    height BIGINT,
    "timestamp" BIGINT,
    piece_cid TEXT,
    piece_size BIGINT NOT NULL,
    verified_deal BOOLEAN NOT NULL,
    "provider" TEXT NOT NULL,
    client TEXT NOT NULL,
    label TEXT,
    start_epoch BIGINT NOT NULL,
    end_epoch BIGINT NOT NULL,
    start_timestamp BIGINT,
    end_timestamp BIGINT,
    storage_price TEXT NOT NULL,
    storage_price_per_epoch TEXT,
    provider_collateral TEXT,
    client_collateral TEXT,
    PRIMARY KEY(id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE market_deal_proposals;
