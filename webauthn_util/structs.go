package webauthn_util

import (
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
)

// User implements webauthn.User
type User struct {
	ID          uuid.UUID
	DisplayName string
	Name        string
	Credentials []webauthn.Credential
}

func (u User) WebAuthnID() []byte {
	return u.ID[:]
}

func (u User) WebAuthnName() string {
	return u.Name
}

func (u User) WebAuthnDisplayName() string {
	return u.DisplayName
}

func (u User) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}
