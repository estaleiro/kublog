package config

import "github.com/BurntSushi/toml"

// BaseConfig defines basic options
type BaseConfig struct {
	Kubeconfig          string
	LogFilename         string `toml:"log_filename"`
	Period              int
	ElasticSearchConfig ElasticSearchConfig `toml:"elasticsearch"`
}

// ElasticSearchConfig defines ES options
type ElasticSearchConfig struct {
	Hosts               []string
	Indexname           string
	Timeout             string
	EnableSniffer       bool   `toml:"enable_sniffer"`
	HealthCheckInterval string `toml:"health_check_interval"`
}

// ReadConfig reads config file
func ReadConfig(configfile string) (*BaseConfig, error) {
	var config BaseConfig
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
