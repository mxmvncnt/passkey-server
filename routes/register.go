package routes

import (
	"encoding/json"
	"net/http"
	"passkey-server/database"
	"passkey-server/utils"
	"passkey-server/utils/logger"
	webauthn_util "passkey-server/webauthn_util"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
)

func (handler *RoutesHandler) BeginRegistrationForNewUser(w http.ResponseWriter, r *http.Request) error {
	var requestBody struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		logger.Errorf("[register.go - BeginRegistrationForNewUser] failed to decode JSON body: %s", err)
		return err
	}

	user := webauthn_util.User{
		ID:          uuid.New(),
		Name:        requestBody.Email,
		DisplayName: requestBody.Email,
	}

	opts := []webauthn.RegistrationOption{
		webauthn.WithResidentKeyRequirement(protocol.ResidentKeyRequirementRequired),
		webauthn.WithExclusions(webauthn.Credentials(user.WebAuthnCredentials()).CredentialDescriptors()),
		webauthn.WithExtensions(map[string]any{"credProps": true}),
		webauthn.WithAuthenticatorSelection(protocol.AuthenticatorSelection{
			AuthenticatorAttachment: protocol.Platform,
			ResidentKey:             protocol.ResidentKeyRequirementRequired,
			UserVerification:        protocol.VerificationRequired,
		}),
	}

	options, session, err := handler.wa.BeginMediatedRegistration(user, protocol.MediationConditional, opts...)
	if err != nil {
		return err
	}

	webauthn_util.NewUserSessionStore[user.ID.String()] = &webauthn_util.SessionEntry{
		User:    user,
		Session: *session,
	}

	utils.SendJsonResponse(w, http.StatusOK, options.Response)
	return nil
}

func (handler *RoutesHandler) FinishRegistrationForNewUser(w http.ResponseWriter, r *http.Request) error {
	userIDBase64 := r.FormValue("user_id")
	parsedUuid := utils.GetUuidFromBase64(userIDBase64)

	data := webauthn_util.NewUserSessionStore[parsedUuid.String()]

	credential, err := handler.wa.FinishRegistration(
		data.User,
		data.Session,
		r,
	)
	if err != nil {
		return err
	}

	err = handler.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:   data.User.ID,
		Name: data.User.Name,
	})
	if err != nil {
		return err
	}

	transports := make([]string, len(credential.Transport))
	for i, c := range credential.Transport {
		transports[i] = string(c)
	}

	err = handler.db.CreateCredential(r.Context(), database.CreateCredentialParams{
		ID:                 credential.ID,
		UserID:             data.User.ID,
		Nickname:           "",
		PublicKey:          credential.PublicKey,
		AttestationType:    credential.AttestationType,
		Aaguid:             credential.Authenticator.AAGUID,
		SignCount:          int64(credential.Authenticator.SignCount),
		Transports:         transports,
		UserPresentFlag:    credential.Flags.UserPresent,
		UserVerifiedFlag:   credential.Flags.UserVerified,
		BackupEligibleFlag: credential.Flags.BackupEligible,
		BackupStateFlag:    credential.Flags.BackupState,
		CloneWarning:       credential.Authenticator.CloneWarning,
	})
	if err != nil {
		return err
	}

	delete(webauthn_util.NewUserSessionStore, parsedUuid.String())
	return nil
}
