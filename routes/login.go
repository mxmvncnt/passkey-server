package routes

import (
	"errors"
	"net/http"
	"passkey-server/database"
	"passkey-server/utils"
	"passkey-server/utils/apierror"
	"passkey-server/webauthn_util"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
)

func (handler *RoutesHandler) BeginLogin(w http.ResponseWriter, r *http.Request) error {
	options, session, err := handler.wa.BeginDiscoverableMediatedLogin(protocol.MediationDefault)
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

		dbCredentials, err := handler.db.ListCredentialsForUser(r.Context(), uid)
		if err != nil {
			return webauthn_util.User{}, err
		}

		credentials := make([]webauthn.Credential, len(dbCredentials))
		for i, c := range dbCredentials {
			transports := make([]protocol.AuthenticatorTransport, len(c.Transports))
			for j, t := range c.Transports {
				transports[j] = protocol.AuthenticatorTransport(t)
			}

			aaguid, err := c.Aaguid.MarshalBinary()
			if err != nil {
				return webauthn_util.User{}, err
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
					AAGUID:       aaguid,
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
		if protocolErr, ok := errors.AsType[*protocol.Error](err); ok {
			switch protocolErr.Error() {
			case "Unable to find the credential for the returned credential ID":
				return apierror.NewApiError(
					http.StatusForbidden,
					"no_credential_found",
					"The credential could not be found or was null",
					protocolErr.Error())
			}
		}
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
