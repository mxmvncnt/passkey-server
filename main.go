package main

import (
	"context"
	"net/http"
	"passkey-server/config"
	"passkey-server/database"
	"passkey-server/middleware"
	"passkey-server/routes"
	"passkey-server/utils/logger"

	"github.com/go-webauthn/webauthn/webauthn"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/cors"
)

func main() {
	router := http.NewServeMux()

	logger.Infof("Loaded config `%s`", config.ConfigName)

	dbPool, dbErr := pgxpool.New(context.Background(), config.DatabaseURL)
	if dbErr != nil {
		logger.Fatalf("Failed to initialize database: %v", dbErr)
	}
	defer dbPool.Close()
	queries := database.New(dbPool)
	dbErr = dbPool.Ping(context.Background())
	if dbErr != nil {
		logger.Fatalf("Failed to connect to database: %v", dbErr)
	}

	var err error

	webAuthn, err := webauthn.New(&webauthn.Config{
		RPDisplayName: config.RPDisplayName,
		RPID:          config.RPID,
		RPOrigins:     config.RPOrigins,
	})
	if err != nil {
		logger.Fatalf("Failed to initialize webauthn_util: %v", err)
	}

	routesHandler := routes.NewRoutesHandler(queries, webAuthn)
	router.HandleFunc("GET /ping", middleware.Combined(routesHandler.Ping))
	router.HandleFunc("POST /passkey/register/begin", middleware.Combined(routesHandler.BeginRegistrationForNewUser))
	router.HandleFunc("POST /passkey/register/finish", middleware.Combined(routesHandler.FinishRegistrationForNewUser))

	router.HandleFunc("POST /passkey/add/begin", middleware.Combined(routesHandler.BeginRegistrationForExistingUser))
	router.HandleFunc("POST /passkey/add/finish", middleware.Combined(routesHandler.FinishRegistrationForExistingUser))

	router.HandleFunc("POST /passkey/login/begin", middleware.Combined(routesHandler.BeginLogin))
	router.HandleFunc("POST /passkey/login/finish", middleware.Combined(routesHandler.FinishLogin))

	router.HandleFunc("GET /credentials/{userID}/list", middleware.Combined(routesHandler.GetCredentialsList))
	router.HandleFunc("GET /users/{userID}", middleware.Combined(routesHandler.GetUser))

	logger.Info("Server started on http://" + config.ServerHostname + ":" + config.ServerPort)
	handler := cors.AllowAll().Handler(router)
	err = http.ListenAndServe(config.ServerHostname+":"+config.ServerPort, handler)
	if err != nil {
		logger.Fatalf("Failed to start server: %v", err)
		return
	}
}
