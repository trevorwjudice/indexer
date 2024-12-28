drop schema if exists indexer cascade;

CREATE SCHEMA indexer;

CREATE TABLE indexer.progress(
    slot_start BIGINT NOT NULL,
    slot_end BIGINT NOT NULL,
    status smallint NOT NULL,
    block_count integer NOT NULL,
    time_taken integer NOT NULL,
    CONSTRAINT indexer_progress_pkey PRIMARY KEY (slot_start, slot_end)
);

drop schema if exists raydium_v4 cascade;

CREATE SCHEMA raydium_v4;

CREATE TABLE raydium_v4.swaps(
    slot BIGINT NOT NULL,
    transaction_index INTEGER NOT NULL,
    instruction_index SMALLINT NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    signature BYTEA NOT NULL,
    pool BYTEA NOT NULL,
    maker BYTEA NOT NULL,
    amount_base BIGINT NOT NULL,
    amount_quote BIGINT NOT NULL,
    CONSTRAINT raydium_v4_swaps_pkey PRIMARY KEY (block, transaction_index, instruction_index)
);

drop schema if exists pump_fun cascade;

CREATE SCHEMA pump_fun;

CREATE TABLE pump_fun.create(
    slot BIGINT NOT NULL,
    transaction_index INTEGER NOT NULL,
    instruction_index SMALLINT NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    signature BYTEA NOT NULL,
    mint BYTEA NOT NULL,
    deployer BYTEA NOT NULL,
    bonding_curve BYTEA NOT NULL,
    associated_bonding_curve BYTEA NOT NULL,
    metadata_slot BYTEA NOT NULL,
    name TEXT NOT NULL,
    symbol TEXT NOT NULL,
    metadata_uri TEXT NOT NULL,
    CONSTRAINT pump_fun_create_pkey PRIMARY KEY (block, transaction_index, instruction_index)
);

CREATE TABLE pump_fun.swaps(
    slot BIGINT NOT NULL,
    transaction_index INTEGER NOT NULL,
    instruction_index SMALLINT NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    signature BYTEA NOT NULL,
    mint BYTEA NOT NULL,
    maker_token_account BYTEA NOT NULL,
    maker BYTEA NOT NULL,
    token_amount BIGINT NOT NULL,
    sol_amount BIGINT NOT NULL,
    fee BIGINT NOT NULL,
    CONSTRAINT pump_fun_swaps_pkey PRIMARY KEY (block, transaction_index, instruction_index)
);
