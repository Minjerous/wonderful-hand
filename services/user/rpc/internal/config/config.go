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
	Network struct {
		Host string
		Addr string
	}

	Server struct {
		Name string
	}

	Mysql struct {
		DataSourceName  string
		MaxIdleConns    int
		MaxOpenConns    int
		ConnMaxLifetime int
	}

	Redis []struct {
		Host     string
		Password string
		Type     string
	}

	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	Etcd struct {
		Endpoints []string
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
