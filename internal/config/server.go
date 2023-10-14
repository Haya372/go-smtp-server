package config

import "time"

type ServerConfig struct {
	Port              int           `yaml:"port"`
	MaxConnection     int           `yaml:"maxConnection"`
	ConnectionTimeout time.Duration `yaml:"connectionTimeout"`
}

func NewServerConfig(conf *Config) *ServerConfig {
	return conf.Server
}
