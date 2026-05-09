CREATE TABLE IF NOT EXISTS transactions (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id          UUID            NOT NULL REFERENCES users(id),
    target_user_id   UUID            REFERENCES users(id),
    amount           DECIMAL(15, 2)  NOT NULL,
    remarks          TEXT            NOT NULL DEFAULT '',
    balance_before   DECIMAL(15, 2)  NOT NULL,
    balance_after    DECIMAL(15, 2)  NOT NULL,
    status           VARCHAR(20)     NOT NULL DEFAULT 'SUCCESS',
    transaction_type VARCHAR(10)     NOT NULL,
    category         VARCHAR(20)     NOT NULL,
    created_at       TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_transactions_user_id    ON transactions (user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_created_at ON transactions (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_transactions_status     ON transactions (status);
