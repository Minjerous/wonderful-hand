package config

import (
	"errors"
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
)

var cfg *Config

type Config struct {
	Server struct {
		// Name is the name of the server as it shows up in the logger.
		Name string
		// Address is the address on which the server should listen.
		Address string

		// HTTPCert specified the certs of http api server
		HTTPCert struct {
			Enable       bool
			CertFilePath string
			KeyFilePath  string
		}
	}

	Etcd struct {
		Endpoints []string
	}

	// GRPC tells server how to discover services
	GRPC struct {
		// User specified user rpc server identity
		User string
	}
}

func ReadConfig() (Config, error) {
	if cfg != nil {
		return *cfg, nil
	}
	c := Config{}
	if _, err := os.Stat("./etc/config.toml"); os.IsNotExist(err) {
		return c, errors.New("no config file found")
	}
	data, err := ioutil.ReadFile("./etc/config.toml")
	if err != nil {
		return c, fmt.Errorf("error reading config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("error decoding config: %v", err)
	}
	cfg = &c
	return c, nil
}
