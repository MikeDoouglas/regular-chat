package config

import "github.com/kelseyhightower/envconfig"

type Environment string

const (
	Dev  Environment = "dev"
	Prod Environment = "prod"
)

type Config struct {
	Environment Environment `envconfig:"ENV" required:"true"`
	Port        int         `envconfig:"PORT" required:"true"`
}

func LoadConfig() (*Config, error) {
	e := new(Config)
	if err := envconfig.Process("", e); err != nil {
		return nil, err
	}

	return e, nil
}

func (c *Config) IsProd() bool {
	return c.Environment == Prod
}
