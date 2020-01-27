package env

import (
	"github.com/kelseyhightower/envconfig"
)

// SentinelEnv is the web configuration object that is populated at launch.
type SentinelEnv struct {
	Env      string `default:"local" required:"true"`
	Port     string `default:"3000" required:"true"`
	LogLevel string `split_words:"true" default:"INFO"`
}

// LoadSentinelEnv loads the configuration object from .env
func LoadSentinelEnv() (*SentinelEnv, error) {
	var se SentinelEnv
	err := envconfig.Process("sen", &se)
	if err != nil {
		return nil, err
	}
	return &se, nil
}
