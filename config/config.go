package config

import (
	"github.com/kelseyhightower/envconfig"
)

type APIConfig struct {
	AddrPort  string `endconfig:"ADDR_PORT" default:":8000"`
	DBConnStr string `envconfig:"DB_URL" required:"true"`
}

func GetAPIConfig() (APIConfig, error) {
	var cfg APIConfig
	if err := envconfig.Process("", &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}
