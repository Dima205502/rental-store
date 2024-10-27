package server

import (
	"auth_service/config"
	"auth_service/internal/controller"
	"auth_service/internal/repository"
	"auth_service/internal/service"
	"log/slog"
	"net/http"
	"strconv"
)

func Start(cfg *config.Config) error {
	storage := repository.NewStorage(cfg.DB)

	userManager := service.NewUserManager(storage)
	sessionManager := service.NewSessionManager(storage)
	emailManager := service.NewEmailManager(storage)

	authManager := controller.NewAuthManager(userManager, sessionManager, emailManager)

	slog.Debug("AuthManager created")

	http.HandleFunc("/signup", controller.Signup(authManager))
	http.HandleFunc("/signin", controller.Signin(authManager))
	http.HandleFunc("/logout", controller.Logout(authManager))
	http.HandleFunc("/verify-email", controller.VerifyEmail(authManager))

	slog.Debug("Starting the server")

	return http.ListenAndServe(cfg.Server.Host+":"+strconv.Itoa(cfg.Server.Port), nil)
}
