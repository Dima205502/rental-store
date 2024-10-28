package config

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB           `json:"db"`
	Server       `json:"server"`
	Broker_addrs []string `json:"broker_addrs"`
	Topic        string   `json:"topic"`
}

type DB struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Dbname   string `json:"dbname"`
}

type Server struct {
	Host string
	Port int
}

func Init() *Config {
	var cfg Config

	err := cleanenv.ReadConfig("/home/kinder/rental-store/auth/config.json", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Debug(fmt.Sprintf("Configuration: %+v\n", cfg))

	return &cfg
}
