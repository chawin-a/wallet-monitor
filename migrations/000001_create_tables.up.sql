CREATE TYPE wallet_transaction_type AS ENUM ('INCOMING', 'OUTGOING', 'UNKNOWN');

CREATE TABLE IF NOT EXISTS "wallets" (
    "wallet" text NOT NULL PRIMARY KEY,
    "worker_id" int,
    "last_block" int8
);

CREATE TABLE IF NOT EXISTS "wallet_transactions" (
    "wallet" text NOT NULL, 
    "tx_hash" text NOT NULL,
    "block_number" int8,
    "block_hash" text,
    "from" text,
    "to" text,
    "gas" numeric,
    "gas_price" numeric,
    "nonce" int,
    "input" bytea,
    "transaction_index" int,
    "value" numeric,
    "type" wallet_transaction_type,
    "status" int,
    PRIMARY KEY("wallet", "tx_hash")
);

CREATE INDEX wallet_oldest_last_block_by_worker_id ON wallets (
    worker_id, 
    last_block
);