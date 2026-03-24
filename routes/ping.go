package routes

import (
	"net/http"
	"passkey-server/utils"
)

func (handler *RoutesHandler) Ping(w http.ResponseWriter, r *http.Request) error {
	utils.SendJsonResponse(w, http.StatusOK, "pong!")
	return nil
}
