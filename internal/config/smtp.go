package config

type SmtpConfig struct {
	// ESMTP extensions
	EnablePipelining bool `yaml:"enablePipelining"`
	Enable8BitMime   bool `yaml:"enable8BitMime"`
	EnableSize       bool `yaml:"enableSize"`
	EnableStartTls   bool `yaml:"enableStartTls"`

	MaxMailSize int `yaml:"maxMailSize"`
}

func NewSmtpConfig(conf *Config) *SmtpConfig {
	return conf.Smtp
}
