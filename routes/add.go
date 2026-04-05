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

func (handler *RoutesHandler) BeginRegistrationForExistingUser(w http.ResponseWriter, r *http.Request) error {
	var requestBody struct {
		UserID string `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		logger.Errorf("[register.go - BeginRegistrationForNewUser] failed to decode JSON body: %s", err)
		return err
	}

	parsedUuid, err := uuid.Parse(requestBody.UserID)
	if err != nil {
		return err
	}

	user, err := handler.db.GetUserFromID(r.Context(), parsedUuid)
	if err != nil {
		return err
	}

	webauthnUser := webauthn_util.User{
		ID:          user.ID,
		DisplayName: user.DisplayName.String,
		Name:        user.Name.String,
		Credentials: nil,
	}

	opts := []webauthn.RegistrationOption{
		webauthn.WithResidentKeyRequirement(protocol.ResidentKeyRequirementRequired),
		webauthn.WithExclusions(webauthn.Credentials(webauthnUser.WebAuthnCredentials()).CredentialDescriptors()),
		webauthn.WithExtensions(map[string]any{"credProps": true}),
		webauthn.WithAuthenticatorSelection(protocol.AuthenticatorSelection{
			AuthenticatorAttachment: protocol.Platform,
			ResidentKey:             protocol.ResidentKeyRequirementRequired,
			UserVerification:        protocol.VerificationRequired,
		}),
	}

	options, session, err := handler.wa.BeginMediatedRegistration(webauthnUser, protocol.MediationConditional, opts...)
	if err != nil {
		return err
	}

	webauthn_util.ExistingUserSessionStore[user.ID.String()] = &webauthn_util.SessionEntry{
		User:    webauthnUser,
		Session: *session,
	}

	utils.SendJsonResponse(w, http.StatusOK, options.Response)
	return nil
}

func (handler *RoutesHandler) FinishRegistrationForExistingUser(w http.ResponseWriter, r *http.Request) error {
	userIDBase64 := r.FormValue("user_id")
	parsedUuid := utils.GetUuidFromBase64(userIDBase64)

	data := webauthn_util.ExistingUserSessionStore[parsedUuid.String()]

	credential, err := handler.wa.FinishRegistration(
		data.User,
		data.Session,
		r,
	)
	if err != nil {
		return err
	}
	delete(webauthn_util.ExistingUserSessionStore, parsedUuid.String())

	transports := make([]string, len(credential.Transport))
	for i, c := range credential.Transport {
		transports[i] = string(c)
	}

	aaguidAsUuid, err := uuid.FromBytes(credential.Authenticator.AAGUID)

	err = handler.db.CreateCredential(r.Context(), database.CreateCredentialParams{
		ID:                 credential.ID,
		UserID:             data.User.ID,
		Nickname:           "",
		PublicKey:          credential.PublicKey,
		AttestationType:    credential.AttestationType,
		Aaguid:             aaguidAsUuid,
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

	return nil
}
