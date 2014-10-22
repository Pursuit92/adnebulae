package main

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Main MainConfig
	Openstack OSConfig
	Chef ChefConfig
}

type MainConfig struct {
	SSL bool
	SSLCert string `toml:"sslcrt"`
	SSLKey string
	Listen string
	Files,ConfBase string
}

type OSConfig struct {
	AuthUrl string `toml:"auth-url"`
}

type ChefConfig struct {
	ServerUrl string `toml:"server-url"`
	ValidatorKey string `toml:"validator-key"`
	ValidatorName string `toml:"validator-name"`
	ClientConfig string `toml:"client-config"`
}

func readConfig(path string) (*Config,error) {
	var config Config
	_,err := toml.DecodeFile(path,&config)
	if err != nil {
		return nil,err
	}

	return &config,err
}


