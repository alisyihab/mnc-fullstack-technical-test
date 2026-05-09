CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name  VARCHAR(100)        NOT NULL DEFAULT '',
    last_name   VARCHAR(100)        NOT NULL DEFAULT '',
    phone_number VARCHAR(20)        NOT NULL UNIQUE,
    address     TEXT                NOT NULL DEFAULT '',
    pin         VARCHAR(255)        NOT NULL,
    balance     DECIMAL(15, 2)      NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_users_phone_number ON users (phone_number);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at   ON users (deleted_at);
