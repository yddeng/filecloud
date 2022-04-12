package service

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	WebAddr          string `toml:"WebAddr"`
	WebIndex         string `toml:"WebIndex"`
	FilePath         string `toml:"FilePath"`
	FileDiskTotal    uint64 `toml:"FileDiskTotal"`
	SaveFileMultiple bool   `toml:"SaveFileMultiple"`
}

var config *Config

func LoadConfig(path string) *Config {
	conf := &Config{}
	_, err := toml.DecodeFile(path, conf)
	if err != nil {
		panic(err)
	}
	config = conf
	logger.Info(config)
	return config
}
