package config

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Sender       string   `json:"sender"`
	AppPassword  string   `json:"app_password"`
	SmtpHost     string   `json:"smtp_host"`
	SmtpPort     int      `json:"smtp_port"`
	Broker_addrs []string `json:"broker_addrs"`
	Topic        string   `json:"topic"`
}

func Init() *Config {
	var cfg Config

	err := cleanenv.ReadConfig("/home/kinder/rental-store/notifier/config.json", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Debug(fmt.Sprintf("Configuration: %+v\n", cfg))

	return &cfg
}
