-- name: CreateCredential :exec
INSERT INTO webauthn_credentials (id, user_id, public_key, attestation_type, aaguid, sign_count, transports, user_present_flag, user_verified_flag, backup_eligible_flag, backup_state_flag, clone_warning)
VALUES (@id::bytea, @user_id::uuid, @public_key::bytea, @attestation_type::text, @aaguid::bytea, @sign_count::bigint, @transports::text[], @user_present_flag::boolean, @user_verified_flag::boolean, @backup_eligible_flag::boolean, @backup_state_flag::boolean, @clone_warning::boolean);
-- name: ListCredentialsByUser :many
SELECT *
FROM webauthn_credentials
WHERE user_id = @user_id::uuid;

-- name: IsEmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = @email::text);

-- name: CreateUser :exec
INSERT INTO users (id, email) VALUES (@id::uuid, @email::text);

-- name: UpdateSignCountForCredential :exec
UPDATE webauthn_credentials
SET sign_count = @sign_count::bigint
WHERE id = @id::bytea;

-- name: GetUserFromID :one
SELECT * FROM users WHERE id = @id::uuid;