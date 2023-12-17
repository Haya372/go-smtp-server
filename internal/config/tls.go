package config

import "crypto/tls"

type TlsConfig struct {
	CertFilePath string `yaml:"certFilePath"`
	KeyFilePath  string `yaml:"keyFilePath"`

	TlsConfig *tls.Config
}

func NewTlsConfig(config *Config) *TlsConfig {
	tlsConf := config.Tls
	if tlsConf == nil {
		return nil
	}

	cer, err := tls.LoadX509KeyPair(tlsConf.CertFilePath, tlsConf.KeyFilePath)
	if err != nil {
		panic(err)
	}

	tlsConf.TlsConfig = &tls.Config{Certificates: []tls.Certificate{cer}}
	return tlsConf
}
