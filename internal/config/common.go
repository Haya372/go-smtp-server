package config

import "time"

type Config struct {
	Server *ServerConfig `yaml:"server"`
	Smtp   *SmtpConfig   `yaml:"smtp"`
	Tls    *TlsConfig    `yaml:"tls"`
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
			EnableStartTls:   true,

			MaxMailSize: 1048576,
		},
		Tls: &TlsConfig{
			CertFilePath: "server.crt",
			KeyFilePath:  "server.key",
		},
	}
}
