CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE users (
    id UUID PRIMARY KEY, -- generated in Go

    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,

    phone_number VARCHAR(20),
    wallet_address VARCHAR(255),

    subscribed BOOLEAN NOT NULL DEFAULT false,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    deleted_at TIMESTAMPTZ
);

-- Uniqueness constraints 
CREATE UNIQUE INDEX idx_users_email ON users (email);
CREATE UNIQUE INDEX idx_users_wallet_address ON users (wallet_address);

-- Soft delete support
CREATE INDEX idx_users_deleted_at ON users (deleted_at);