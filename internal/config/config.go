package config

import (
    "github.com/ilyakaznacheev/cleanenv"
    "log"
    "os"
)

type Config struct {
    Env  string `yaml:"env" json:"env"`
    Addr string `yaml:"addr" json:"addr"`
    Port string `yaml:"port" json:"port"`

    Backends []Backend `yaml:"backends" json:"backends"`
}

type Backend struct {
    Url string `yaml:"url" json:"url"`
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
