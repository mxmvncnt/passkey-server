package routes

import (
	"net/http"
	"passkey-server/utils"

	"github.com/google/uuid"
)

func (handler *RoutesHandler) GetUser(w http.ResponseWriter, r *http.Request) error {
	userIDParam := r.PathValue("userID")

	userId, err := uuid.Parse(userIDParam)
	if err != nil {
		return err
	}

	user, err := handler.db.GetUserFromID(r.Context(), userId)
	if err != nil {
		return err
	}

	utils.SendJsonResponse(w, http.StatusOK, user)
	return nil
}
