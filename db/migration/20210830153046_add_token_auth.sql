-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE IF NOT EXISTS token_auth
(
    "id" TEXT NOT NULL,
    "digest" TEXT NOT NULL,
    PRIMARY KEY (id)
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.

DROP TABLE token_auth;