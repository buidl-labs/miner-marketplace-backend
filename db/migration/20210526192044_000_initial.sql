-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE IF NOT EXISTS miners
(
    id TEXT NOT NULL,
    claimed BOOLEAN,
    region TEXT ,
    country TEXT,
    worker_id TEXT,
    worker_address TEXT,
    owner_id TEXT,
    owner_address TEXT,
    quality_adjusted_power TEXT,
    storage_ask_price TEXT,
    verified_ask_price TEXT,
    retrieval_ask_price TEXT,
    reputation_score INTEGER,
    transparency_score INTEGER,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS miner_personal_infos
(
    id TEXT NOT NULL,
    name TEXT,
    bio TEXT,
    email TEXT,
    website TEXT,
    twitter TEXT,
    slack TEXT,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS miner_services
(
    id TEXT NOT NULL,
    storage BOOLEAN,
    retrieval BOOLEAN,
    repair BOOLEAN,
    data_transfer_online BOOLEAN,
    data_transfer_offline BOOLEAN,
    PRIMARY KEY (id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE miner_services;
DROP TABLE miner_personal_infos;
DROP TABLE miners;
