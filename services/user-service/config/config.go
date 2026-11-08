package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"gopkg.in/ini.v1"
)

const (
	defaultHTTPAddr     = ":8080"
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 10 * time.Second
)

type Config struct {
	HTTPAddr     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	DatabaseURL  string
}

func Load(path string) (Config, error) {
	file, err := ini.Load(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config %s: %w", path, err)
	}
	cfg := Config{
		HTTPAddr:     file.Section("server").Key("addr").MustString(defaultHTTPAddr),
		ReadTimeout:  file.Section("server").Key("read_timeout").MustDuration(defaultReadTimeout),
		WriteTimeout: file.Section("server").Key("write_timeout").MustDuration(defaultWriteTimeout),
		DatabaseURL:  file.Section("database").Key("url").String(),
	}
	if v := os.Getenv("DATABASE_URL"); v != "" {
		cfg.DatabaseURL = v
	}
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("config: database url is required")
	}
	return cfg, nil
}
