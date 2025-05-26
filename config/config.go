package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("[config] no .env file found, using env vars")
	}
}

var Values = struct {
	KeyDB struct {
		Address  string `env:"KEYDB_ADDR,required"`
		Password string `env:"KEYDB_PASSWORD,required"`
	}
}{}

func LoadConfig() {
	if err := env.Parse(&Values); err != nil {
		log.Fatalf("[LoadConfig] failed to parse config: %v", err)
	}
}
