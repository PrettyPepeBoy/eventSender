package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env     string     `yaml:"env"`
	Storage string     `yaml:"storage"`
	Server  HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Port        string        `yaml:"port"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
	User        string        `yaml:"user"`
	Password    string        `yaml:"password"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("file does not exist :%s", configPath)
	}

	var config Config

	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		log.Fatalf("failed to read config file :%s", configPath)
	}

	return &config
}
