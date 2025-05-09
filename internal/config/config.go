package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env  string `yaml:"env" json:"env"`
	Host string `yaml:"host" json:"host"`

	Backends []Backend `yaml:"backends" json:"backends"`

	RateLimiter RateLimiter `yaml:"rate_limiter" json:"rate_limiter"`
}

type Backend struct {
	Url string `yaml:"url" json:"url"`
}

type RateLimiter struct {
	DefaultBucketCapacity int           `yaml:"default_bucket_capacity" json:"default_bucket_capacity"`
	DefaultRefill         int           `yaml:"default_refill" json:"default_refill"`
	DefaultRefillInterval time.Duration `yaml:"default_refill_interval" json:"default_refill_interval"`

	BucketSettingsDatabase BucketSettingsDatabase `yaml:"bucket_settings_database" json:"bucket_settings_database"`

	API API `yaml:"api" json:"api"`
}

type BucketSettingsDatabase struct {
	Host     string `yaml:"host" json:"host"`
	Port     string `yaml:"port" json:"port"`
	User     string `yaml:"user" json:"user"`
	Database string `yaml:"database" json:"database"`
	Password string `yaml:"password" json:"password"`
}

type API struct {
	Host string `yaml:"host" json:"host"`
}

const (
	configPathEnvVariableName = "CONFIG_PATH"
)

func MustLoad() *Config {
	configPath := os.Getenv(configPathEnvVariableName)
	if configPath == "" {
		log.Fatalf("%s env variable is not set", configPathEnvVariableName)
	}

	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		log.Fatalf("config file not found: %s", configPath)
	}

	var cfg Config

	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("cannot read config file: %s", err.Error())
	}

	return &cfg
}
