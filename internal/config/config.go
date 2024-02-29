package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string `yaml:"env" env-default:"local"`
	Version        string `yaml:"version" env-default:"unknown"`
	Port           int    `yaml:"port" env-default:""`
	TTL            int    `yaml:"ttl" env-default:"60"`
	ClearInterval  int    `yaml:"clear_interval" env-default:"3"`
	UseTracing     bool   `yaml:"use_tracing"`
	TracingAddress string `yaml:"tracing_address" env-default:""`
}

func MustLoad(name string) *Config {
	configPath := os.Getenv(strings.ToUpper(name) + "_CONFIG_PATH")
	if configPath == "" {
		configPath = "config/" + name + "/local.yaml"
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	portEnv := os.Getenv(strings.ToUpper(name) + "_PORT")
	if portEnv != "" {
		newPort, err := strconv.Atoi(portEnv)
		if err == nil {
			cfg.Port = newPort
		}
	}

	return &cfg
}
