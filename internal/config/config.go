package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Env             string `yaml:"env" env:"ENV" env-default:"local"`
	HTTPServer      `yaml:"http_server"`
	Clients         ClientsConfig `yaml:"clients"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"SHUTDOWN_TIMEOUT" env-default:"10s"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env:"HTTP_ADDRESS" env-default:"localhost:5000"`
	Timeout     time.Duration `yaml:"timeout" env:"HTTP_TIMEOUT" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT" env-default:"60s"`
}

type ClientsConfig struct {
	User struct {
		Address      string        `yaml:"address" env:"USER_SERVICE_ADDRESS"`
		Timeout      time.Duration `yaml:"timeout" env:"USER_SERVICE_TIMEOUT"`
		RetriesCount int           `yaml:"retries_count" env:"USER_SERVICE_RETRIES_COUNT"`
	} `yaml:"user"`
	Club struct {
		Address      string        `yaml:"address" env:"CLUB_SERVICE_ADDRESS"`
		Timeout      time.Duration `yaml:"timeout" env:"CLUB_SERVICE_TIMEOUT"`
		RetriesCount int           `yaml:"retries_count" env:"CLUB_SERVICE_RETRIES_COUNT"`
	} `yaml:"club"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		return MustLoadFromEnv()
	}

	return MustLoadByPath(path)
}

func MustLoadByPath(configPath string) *Config {
	cfg, err := LoadByPath(configPath)
	if err != nil {
		panic(err)
	}

	return cfg
}

func LoadByPath(configPath string) (*Config, error) {
	var cfg Config

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("there is no config file: %w", err)
	}

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	return &cfg, nil
}

func MustLoadFromEnv() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("Env empty")
	}
	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
