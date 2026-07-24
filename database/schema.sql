CREATE TABLE users
(
    id            UUID PRIMARY KEY,
    project_id    UUID,       -- FK added later via ALTER TABLE
    name          TEXT,
    display_name  TEXT
);

CREATE TABLE projects
(
    id         UUID PRIMARY KEY,
    owner_id   UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    name       TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_projects_owner_id ON projects (owner_id);

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

ALTER TABLE users ADD CONSTRAINT fk_users_project_id FOREIGN KEY (project_id) REFERENCES projects (id) ON DELETE CASCADE;
CREATE UNIQUE INDEX idx_users_project_name ON users (project_id, name); -- do not allow duplicate users.name in the same project
