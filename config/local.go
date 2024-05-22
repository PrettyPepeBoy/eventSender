package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env        string       `yaml:"env"`
	Storage    string       `yaml:"storage"`
	Server     BrokerServer `yaml:"broker_server"`
	Postgres   Postgresql   `yaml:"postgresql"`
	MailSender MailSend     `yaml:"mail_sender"`
}

type BrokerServer struct {
	Port        string        `yaml:"port"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
	Timeout     time.Duration `yaml:"timeout"`
	User        string        `yaml:"user"`
	Password    string        `yaml:"password"`
}

type Postgresql struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type MailSend struct {
	Mail     string `yaml:"mail"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
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
