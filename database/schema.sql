CREATE TABLE users
(
    id           UUID PRIMARY KEY,
    name         TEXT,
    display_name TEXT
);

CREATE TABLE webauthn_credentials
(
    id                   BYTEA PRIMARY KEY,
    user_id              UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    nickname             TEXT,
    public_key           BYTEA       NOT NULL,
    attestation_type     TEXT,
    aaguid               UUID,
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

create table if not exists session
(
    id                      uuid DEFAULT gen_random_uuid()  not null constraint session_pk primary key,
    created_at              timestamp with time zone DEFAULT now() NOT NULL,
    last_used_at            timestamp with time zone,
    created_at_ip           text                            NOT NULL,
    created_at_user_agent   text,
    device_nickname         text,
    token                   text                            NOT NULL,
    user_id                 uuid                            NOT NULL,
    expires_at              timestamp with time zone,
    is_long                 boolean DEFAULT false           NOT NULL
);

CREATE INDEX idx_session_token on session(token);
CREATE INDEX idx_user_id on session(user_id);