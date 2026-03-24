package routes

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"passkey-server/database"
)

// Filename starts with 0 to be at the top of the list

type RoutesHandler struct {
	db *database.Queries
	wa *webauthn.WebAuthn
}

func NewRoutesHandler(pool *database.Queries, webAuthn *webauthn.WebAuthn) *RoutesHandler {
	return &RoutesHandler{db: pool, wa: webAuthn}
}
