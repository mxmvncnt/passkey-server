CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE
);

CREATE TABLE webauthn_credentials (
    id BYTEA PRIMARY KEY,
    user_id UUID NOT NULL,
    public_key BYTEA NOT NULL,
    attestation_type TEXT,
    aaguid BYTEA,
    sign_count BIGINT,
    transports TEXT[]
);