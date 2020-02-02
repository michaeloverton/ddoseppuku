package env

import (
	"github.com/kelseyhightower/envconfig"
)

// LaserEnv is the web configuration object that is populated at launch.
type LaserEnv struct {
	Env          string `default:"local" required:"true"`
	Port         string `default:"3000" required:"true"`
	LogLevel     string `split_words:"true" default:"INFO"`
	RedisAddress string `split_words:"true" required:"true"`
}

// LoadLaserEnv loads the configuration object from .env
func LoadLaserEnv() (*LaserEnv, error) {
	var le LaserEnv
	err := envconfig.Process("lsr", &le)
	if err != nil {
		return nil, err
	}
	return &le, nil
}
