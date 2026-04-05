package webauthn_util

import "github.com/go-webauthn/webauthn/webauthn"

type SessionEntry struct {
	User    User
	Session webauthn.SessionData
}

var ExistingUserSessionStore = map[string]*SessionEntry{}
var NewUserSessionStore = map[string]*SessionEntry{}
var LoginSessionStore = map[string]*webauthn.SessionData{}
