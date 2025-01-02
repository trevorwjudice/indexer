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

CREATE TABLE raydium_v4.initialize2(
    slot BIGINT NOT NULL,
    transaction_index INTEGER NOT NULL,
    instruction_index SMALLINT NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    signature BYTEA NOT NULL,
    pool BYTEA NOT NULL,
    minter BYTEA NOT NULL,
    coin_mint BYTEA NOT NULL,
    pool_coin_token_account BYTEA NOT NULL,
    pc_mint BYTEA NOT NULL,
    pool_pc_token_account BYTEA NOT NULL,
    lp_mint BYTEA NOT NULL,
    nonce smallint NOT NULL,
    init_pc_amount BIGINT NOT NULL,
    init_coin_amount BIGINT NOT NULL,
    lp_amount BIGINT NOT NULL,
    CONSTRAINT raydium_v4_initialize2_pkey PRIMARY KEY (slot, transaction_index, instruction_index)
);

CREATE TABLE raydium_v4.add_liquidity(
    slot BIGINT NOT NULL,
    transaction_index INTEGER NOT NULL,
    instruction_index SMALLINT NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    signature BYTEA NOT NULL,
    pool BYTEA NOT NULL,
    minter BYTEA NOT NULL,
    amount_base BIGINT NOT NULL,
    amount_quote BIGINT NOT NULL,
    lp_token_amount BIGINT NOT NULL,
    CONSTRAINT raydium_v4_add_liquidity_pkey PRIMARY KEY (slot, transaction_index, instruction_index)
);

CREATE TABLE raydium_v4.remove_liquidity(
    slot BIGINT NOT NULL,
    transaction_index INTEGER NOT NULL,
    instruction_index SMALLINT NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    signature BYTEA NOT NULL,
    pool BYTEA NOT NULL,
    owner BYTEA NOT NULL,
    amount_base BIGINT NOT NULL,
    amount_quote BIGINT NOT NULL,
    lp_token_amount BIGINT NOT NULL,
    CONSTRAINT raydium_v4_remove_liquidity_pkey PRIMARY KEY (slot, transaction_index, instruction_index)
);

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
    CONSTRAINT raydium_v4_swaps_pkey PRIMARY KEY (slot, transaction_index, instruction_index)
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
    CONSTRAINT pump_fun_create_pkey PRIMARY KEY (slot, transaction_index, instruction_index)
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
    CONSTRAINT pump_fun_swaps_pkey PRIMARY KEY (slot, transaction_index, instruction_index)
);

CREATE TABLE pump_fun.set_params(
    slot BIGINT NOT NULL,
    transaction_index INTEGER NOT NULL,
    instruction_index SMALLINT NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    signature BYTEA NOT NULL,
    fee_recipient BYTEA NOT NULL,
    initial_virtual_token_reserves BIGINT NOT NULL,
    initial_virtual_sol_reserves BIGINT NOT NULL,
    initial_real_token_reserves BIGINT NOT NULL,
    token_total_supply BIGINT NOT NULL,
    fee_basis_points BIGINT NOT NULL,
    CONSTRAINT pump_fun_set_params_pkey PRIMARY KEY (slot, transaction_index, instruction_index)
);


drop schema if exists spl cascade;

CREATE SCHEMA spl;

CREATE TABLE spl.initialize_account(
    slot BIGINT NOT NULL,
    transaction_index INTEGER NOT NULL,
    instruction_index SMALLINT NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    signature BYTEA NOT NULL,
    owner BYTEA NOT NULL,
    mint BYTEA NOT NULL,
    account BYTEA NOT NULL,
    CONSTRAINT spl_initialize_account_pkey PRIMARY KEY (slot, transaction_index, instruction_index)
);


CREATE TABLE spl.transfer(
    slot BIGINT NOT NULL,
    transaction_index INTEGER NOT NULL,
    instruction_index SMALLINT NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    signature BYTEA NOT NULL,
    mint BYTEA NOT NULL,
    authority BYTEA NOT NULL,
    source BYTEA NOT NULL,
    destination BYTEA NOT NULL,
    amount BIGINT NOT NULL,
    CONSTRAINT spl_transfer_pkey PRIMARY KEY (slot, transaction_index, instruction_index)
);

CREATE TABLE spl.burn(
    slot BIGINT NOT NULL,
    transaction_index INTEGER NOT NULL,
    instruction_index SMALLINT NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    signature BYTEA NOT NULL,
    mint BYTEA NOT NULL,
    account BYTEA NOT NULL,
    owner BYTEA NOT NULL,
    amount BIGINT NOT NULL,
    CONSTRAINT spl_burn_pkey PRIMARY KEY (slot, transaction_index, instruction_index)
);

CREATE TABLE spl.associated_token_account_create(
    slot BIGINT NOT NULL,
    transaction_index INTEGER NOT NULL,
    instruction_index SMALLINT NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    signature BYTEA NOT NULL,
    account BYTEA NOT NULL,
    mint BYTEA NOT NULL,
    source BYTEA NOT NULL,
    wallet BYTEA NOT NULL,
    CONSTRAINT spl_associated_token_account_create_pkey PRIMARY KEY (slot, transaction_index, instruction_index)
);

