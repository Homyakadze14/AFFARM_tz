package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env            string `env:"ENV" env-default:"local"`
	Database       DatabaseConfig
	HTTP           HTTPConfig
	MigrationsPath string
}

type HTTPConfig struct {
	Port string `env:"HTTP_PORT" env-required:"true"`
}

type DatabaseConfig struct {
	URL     string `env:"PG_URL" env-required:"true"`
	PoolMax int    `env:"PG_POOL_MAX" env-required:"true"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath)
}

func MustLoadPath(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
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
