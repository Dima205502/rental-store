package main

import (
	"auth_service/logger"
	"auth_service/server"
	"log/slog"
)

func main() {
	logger.Init()

	if err := server.Start(); err != nil {
		slog.Error("server not started", "error", err.Error())
	}
}
