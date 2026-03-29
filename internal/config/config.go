package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"local"`
	Database   Database   `yaml:"database"`
	HttpServer HttpServer `yaml:"http_server"`
}

type Database struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5555"`
	User     string `yaml:"user" env-default:"postgres_user"`
	Password string `yaml:"password" env-default:"postgres_pass"`
	Dbname   string `yaml:"dbname" env-default:"url_shortener_db"`
	Sslmode  string `yaml:"sslmode" env-default:"disable"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8090"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &config
}
