-- name: CreateCredential :exec
INSERT INTO webauthn_credentials (id, user_id, nickname, public_key, attestation_type, aaguid, sign_count, transports, user_present_flag, user_verified_flag, backup_eligible_flag, backup_state_flag, clone_warning)
VALUES (@id::bytea, @user_id::uuid, @nickname::text, @public_key::bytea, @attestation_type::text, @aaguid::uuid, @sign_count::bigint, @transports::text[], @user_present_flag::boolean, @user_verified_flag::boolean, @backup_eligible_flag::boolean, @backup_state_flag::boolean, @clone_warning::boolean);

-- name: ListCredentialsForUser :many
SELECT *
FROM webauthn_credentials
WHERE user_id = @user_id::uuid;

-- name: DeleteCredential :exec
DELETE FROM webauthn_credentials
WHERE user_id = @user_id::uuid AND id = @id::uuid;

-- name: IsEmailExists :one
SELECT EXISTS(SELECT 1 FROM users WHERE name = @name::text);

-- name: CreateUser :exec
INSERT INTO users (id, name) VALUES (@id::uuid, @name::text);

-- name: UpdateSignCountForCredential :exec
UPDATE webauthn_credentials
SET sign_count = @sign_count::bigint, last_used_at = NOW()
WHERE id = @id::bytea;

-- name: GetUserFromID :one
SELECT * FROM users WHERE id = @id::uuid;

-- name: CreateSession :one
INSERT INTO session (created_at_ip, token, user_id, expires_at, is_long)
VALUES (@created_at_ip::text, @token::text, @user_id::uuid, @expires_at::timestamptz, @is_long::bool)
RETURNING *;

-- name: GetUserFromToken :one
SELECT sqlc.embed(session), sqlc.embed(users) FROM session
JOIN users ON session.user_id = users.id
WHERE session.token = @token::text
LIMIT 1;