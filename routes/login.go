package routes

import (
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
	"net/http"
	"passkey-server/database"
	"passkey-server/utils"
	"passkey-server/utils/apierror"
	"passkey-server/webauthn_util"
)

func (handler *RoutesHandler) BeginLogin(w http.ResponseWriter, r *http.Request) error {
	options, session, err := handler.wa.BeginDiscoverableLogin()
	if err != nil {
		return err
	}

	sessionID := uuid.New().String()
	webauthn_util.LoginSessionStore[sessionID] = session

	response := map[string]interface{}{
		"options":    options,
		"session_id": sessionID,
	}

	utils.SendJsonResponse(w, http.StatusOK, response)
	return nil
}

func (handler *RoutesHandler) FinishLogin(w http.ResponseWriter, r *http.Request) error {
	sessionID := r.FormValue("session_id")

	userHandler := func(rawID, userHandle []byte) (webauthn.User, error) {
		uid, err := uuid.FromBytes(userHandle)
		if err != nil {
			return webauthn_util.User{}, err
		}

		dbCredentials, err := handler.db.ListCredentialsByUser(r.Context(), uid)
		if err != nil {
			return webauthn_util.User{}, err
		}

		credentials := make([]webauthn.Credential, len(dbCredentials))
		for i, c := range dbCredentials {
			transports := make([]protocol.AuthenticatorTransport, len(c.Transports))
			for j, t := range c.Transports {
				transports[j] = protocol.AuthenticatorTransport(t)
			}

			credentials[i] = webauthn.Credential{
				ID:        c.ID,
				PublicKey: c.PublicKey,
				Transport: transports,
				Flags: webauthn.CredentialFlags{
					UserPresent:    c.UserPresentFlag,
					UserVerified:   c.UserVerifiedFlag,
					BackupEligible: c.BackupEligibleFlag,
					BackupState:    c.BackupStateFlag,
				},
				Authenticator: webauthn.Authenticator{
					AAGUID:       c.Aaguid,
					SignCount:    uint32(c.SignCount),
					CloneWarning: c.CloneWarning,
				},
			}
		}

		return webauthn_util.User{
			ID:          uid,
			Name:        "",
			DisplayName: "",
			Credentials: credentials,
		}, nil
	}

	session := webauthn_util.LoginSessionStore[sessionID]
	if session == nil {
		return apierror.NewApiError(
			http.StatusNotFound,
			"no_passkey_session_found",
			"Please try re-starting the passkey login process",
			"There was no passkey login session found for this ID")
	}

	webauthnUser, credential, err := handler.wa.FinishPasskeyLogin(userHandler, *session, r)
	if err != nil {
		return err
	}

	err = handler.db.UpdateSignCountForCredential(r.Context(), database.UpdateSignCountForCredentialParams{
		ID:        credential.ID,
		SignCount: int64(credential.Authenticator.SignCount),
	})
	if err != nil {
		return err
	}

	parsedUuid, _ := uuid.FromBytes(webauthnUser.WebAuthnID())

	user, err := handler.db.GetUserFromID(r.Context(), parsedUuid)

	delete(webauthn_util.LoginSessionStore, sessionID)
	utils.SendJsonResponse(w, http.StatusOK, user)
	return nil
}
