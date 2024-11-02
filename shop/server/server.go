package server

import (
	"log/slog"
	"net/http"
	"shop/config"
	"strconv"
)

func Start(cfg *config.Config) error {
	slog.Debug("AuthManager created")

	http.HandleFunc("/add-thing")
	http.HandleFunc("/remuve-thing")
	http.HandleFunc("/buy-thing")
	http.HandleFunc("/show-all-things")
	http.HandleFunc("/show-rental-things")
	http.HandleFunc("/show-sale-things")

	slog.Debug("Starting the server")

	return http.ListenAndServe(cfg.Server.Host+":"+strconv.Itoa(cfg.Server.Port), nil)
}
