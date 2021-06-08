package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Bitstamp struct {
		URL       string            `yaml:"url"`
		ID        string            `yaml:"id"`
		Key       string            `yaml:"key"`
		Secret    string            `yaml:"secret"`
		Endpoints map[string]string `yaml:"endpoints"`
	} `yaml:"bitstamp"`
	Binance struct {
		URL       string            `yaml:"url"`
		Key       string            `yaml:"key"`
		Secret    string            `yaml:"secret"`
		Endpoints map[string]string `yaml:"endpoints"`
	} `yaml:"binance"`
}

func LoadConfig() (Config, error) {
	f, err := os.Open("config.yaml")
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err = decoder.Decode(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
