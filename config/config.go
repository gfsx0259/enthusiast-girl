package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Telegram `yaml:"telegram"`
		Sdlc     `yaml:"sdlc"`
		Quay     `yaml:"quay"`
		Stash    `yaml:"stash"`
		Alert    `yaml:"alert"`
	}

	Telegram struct {
		Maintainers []string `yaml:"maintainers"`
		Token       string   `env-required:"true" yaml:"token" env:"TELEGRAM_APITOKEN"`
	}

	Sdlc struct {
		User     string `env-required:"true" yaml:"user" env:"SDLC_USER"`
		Password string `env-required:"true" yaml:"password" env:"SDLC_PASSWORD"`
	}

	Quay struct {
		User     string `env-required:"true" yaml:"user" env:"QUAY_USER"`
		Password string `env-required:"true" yaml:"password" env:"QUAY_PASSWORD"`
	}

	Stash struct {
		User  string `yaml:"user"`
		Email string `yaml:"email"`
	}

	Alert struct {
		HookPort string `env-required:"true" yaml:"hookPort" env:"ALERT_HOOK_PORT"`
		ChatId   int    `env-required:"true" yaml:"chatId" env:"ALERT_CHAT_ID"`
	}
)

func New() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig("./config/config.yaml", cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
