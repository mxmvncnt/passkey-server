package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"passkey-server/database"
	"passkey-server/utils/apierror"
	"passkey-server/utils/logger"
	"time"

	"github.com/google/uuid"
)

const SessionTokenLength = 64

func GetUserFromToken(db *database.Queries, token string) (database.User, error) {
	ctx := context.Background()

	result, err := db.GetUserFromToken(ctx, token)

	if err != nil && err.Error() == "no rows in result set" {
		return database.User{}, apierror.NewApiError(
			http.StatusUnauthorized,
			"invalid_token",
			"Make sure you are logged in correctly and that your session is not expired.",
			"The token you have provided does not exist")
	}

	if err != nil {
		// TODO: proper error handling
		return database.User{}, err
	}

	if result.Session.ExpiresAt.Time.Before(time.Now()) {
		return database.User{}, apierror.NewApiError(
			http.StatusUnauthorized,
			"invalid_token",
			"Make sure you are logged in correctly and that your session is not expired.",
			fmt.Sprintf("The token has expired on %v", result.Session.ExpiresAt))
	}

	logger.Debug("Getting user from Token")

	return result.User, nil
}

func CreateSession(db *database.Queries, userID uuid.UUID, isLongSession bool) (database.Session, error) {
	sessionToken, err := RandomString(SessionTokenLength)
	if err != nil {
		return database.Session{}, err
	}

	var expirationTime time.Time
	if isLongSession {
		expirationTime = time.Now().Add(time.Hour * 24 * 7)
	} else {
		expirationTime = time.Now().Add(time.Hour)
	}

	result, err := db.CreateSession(context.Background(), database.CreateSessionParams{
		CreatedAtIp: "12.34.56.78",
		Token:       sessionToken,
		ExpiresAt:   expirationTime,
		UserID:      userID,
		IsLong:      isLongSession,
	})
	if err != nil {
		// TODO: better error handling
		return database.Session{}, err
	}

	logger.Debugf("Created new session for user: %s", userID)

	return result, nil
}

func RandomString(length int) (string, error) {
	result, err := RandomBytes(length)
	if err != nil {
		logger.Errorf("[session.go - RandomString] Could not generate random string: %v", err)
		return "", err
	}
	return base64.StdEncoding.EncodeToString(result), nil
}

func RandomBytes(length int) ([]byte, error) {
	result := make([]byte, length)
	_, err := rand.Read(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
