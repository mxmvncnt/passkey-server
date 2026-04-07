package routes

import (
	"encoding/json"
	"net/http"
	"passkey-server/database"
	"passkey-server/utils"

	"github.com/google/uuid"
)

func (handler *RoutesHandler) GetCredentialsList(w http.ResponseWriter, r *http.Request) error {
	userIDParam := r.PathValue("userID")

	userID, err := uuid.Parse(userIDParam)
	if err != nil {
		return err
	}

	credentials, err := handler.db.ListCredentialsForUser(r.Context(), userID)
	utils.SendJsonResponse(w, http.StatusOK, credentials)
	return nil
}

func (handler *RoutesHandler) DeleteCredential(w http.ResponseWriter, r *http.Request) error {
	var requestBody struct {
		UserID       string `json:"user_id"`
		CredentialID string `json:"credential_id"`
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		return err
	}

	parsedUserID, err := uuid.Parse(requestBody.UserID)
	if err != nil {
		return err
	}

	//credentialID := utils.GetByteArrayFromBase64(requestBody.CredentialID)
	credentialID := []byte(requestBody.CredentialID)

	err = handler.db.DeleteCredential(r.Context(), database.DeleteCredentialParams{
		UserID: parsedUserID,
		ID:     credentialID,
	})
	if err != nil {
		return err
	}

	return nil
}
