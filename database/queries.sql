-- name: CreateCredential :exec
INSERT INTO webauthn_credentials (id, user_id, public_key, attestation_type, aaguid, sign_count, transports)
VALUES (@id::bytea, @user_id::uuid, @public_key::bytea, @attestation_type::text, @aaguid::bytea, @sign_count::bigint, @transports::text[]);

-- name: ListCredentialsByUser :many
SELECT *
FROM webauthn_credentials
WHERE user_id = @user_id::uuid;

-- name: IsEmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = @email::text);

-- name: CreateUser :exec
INSERT INTO users (id, email) VALUES (@id::uuid, @email::text);