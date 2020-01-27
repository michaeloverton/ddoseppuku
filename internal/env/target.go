package env

import (
	"github.com/kelseyhightower/envconfig"
)

// TargetEnv is the web configuration object that is populated at launch.
type TargetEnv struct {
	Env      string `default:"local" required:"true"`
	Port     string `default:"3000" required:"true"`
	LogLevel string `split_words:"true" default:"INFO"`
}

// LoadTargetEnv loads the configuration object from .env
func LoadTargetEnv() (*TargetEnv, error) {
	var te TargetEnv
	err := envconfig.Process("tgt", &te)
	if err != nil {
		return nil, err
	}
	return &te, nil
}
