package main

import (
	"auth_service/config"
	"auth_service/logger"
	"auth_service/server"
	"auth_service/utils"
	"log/slog"
)

func main() {
	logger.Init()

	cfg := config.Init()

	utils.Init(cfg)

	if err := server.Start(cfg); err != nil {
		slog.Error("server not started", "error", err.Error())
	}
}
