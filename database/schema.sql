CREATE TABLE users
(
    id    UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE
);

CREATE TABLE webauthn_credentials
(
    id                   BYTEA PRIMARY KEY,
    user_id              UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    public_key           BYTEA       NOT NULL,
    attestation_type     TEXT,
    aaguid               BYTEA,
    sign_count           BIGINT      NOT NULL DEFAULT 0,
    transports           TEXT[],
    user_present_flag    BOOLEAN     NOT NULL DEFAULT false,
    user_verified_flag   BOOLEAN     NOT NULL DEFAULT false,
    backup_eligible_flag BOOLEAN     NOT NULL DEFAULT false,
    backup_state_flag    BOOLEAN     NOT NULL DEFAULT false,
    clone_warning        BOOLEAN     NOT NULL DEFAULT false,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used_at         TIMESTAMPTZ
);

CREATE INDEX idx_webauthn_credentials_user_id ON webauthn_credentials (user_id);