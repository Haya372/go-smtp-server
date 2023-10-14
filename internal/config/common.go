package config

import "time"

type Config struct {
	Server *ServerConfig `yaml:"server"`
	Smtp   *SmtpConfig   `yaml:"smtp"`
}

func NewDefaultConfig() *Config {
	return &Config{
		Server: &ServerConfig{
			Port:              25,
			MaxConnection:     10,
			ConnectionTimeout: 30 * time.Second,
		},
		Smtp: &SmtpConfig{
			EnablePipelining: true,
			Enable8BitMime:   true,
			EnableSize:       true,
			EnableStartTls:   false,

			MaxMailSize: 1048576,
		},
	}
}
