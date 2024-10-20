package server

import (
	"auth_service/internal/controller"
	"auth_service/internal/repository"
	"auth_service/internal/service"
	"log/slog"
	"net/http"
)

func Start() error {
	storage := repository.NewStorage()

	userManager := service.NewUserManager(storage)
	sessionManager := service.NewSessionManager(storage)
	emailManager := service.NewEmailManager(storage)

	authManager := controller.NewAuthManager(userManager, sessionManager, emailManager)

	slog.Debug("AuthManager created")

	http.HandleFunc("/signup", controller.Signup(authManager))
	http.HandleFunc("/signin", controller.Signin(authManager))
	http.HandleFunc("/logout", controller.Logout(authManager))
	http.HandleFunc("/verify-email", controller.VerifyEmail(authManager))

	slog.Debug("starting the server")

	return http.ListenAndServe(":8080", nil)
}
