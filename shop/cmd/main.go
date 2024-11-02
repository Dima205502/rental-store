package main

import (
	"log/slog"
	"shop/config"
	"shop/logger"
	"shop/server"
)

func main() {
	logger.Init()

	cfg := config.Init()

	if err := server.Start(cfg); err != nil {
		slog.Error("server not started", "error", err.Error())
	}
}
