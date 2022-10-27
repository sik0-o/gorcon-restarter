package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Servers map[string]ServerConf
}

type ServerConf struct {
	Host     string
	Port     int
	Password string
	Restart  RestartConf
}

type RestartConf struct {
	Period        time.Duration
	ServerLock    time.Duration
	Announcements Announcements
}

type Announcements struct {
	At  string
	Min string
	Sec string
}

func Read() (*Config, error) {
	conf := new(Config)
	b, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(b, conf); err != nil {
		return nil, err
	}

	return conf, nil
}
