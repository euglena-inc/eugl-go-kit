package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

func Load[T any]() (T, error) {
	var cfg T
	if err := env.Parse(&cfg); err != nil {
		return cfg, fmt.Errorf("parse env config: %w", err)
	}
	return cfg, nil
}
