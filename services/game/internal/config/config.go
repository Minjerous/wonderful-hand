package config

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"os"
)

type Config struct {
	Network struct {
		GameAddr string
		GRPCAddr string
	}
	Server struct {
		Name      string
		Multicore bool
		ReusePort bool
	}

	Etcd struct {
		Endpoints []string
	}

	GRPC struct {
		User  string
		Chess string
		Room  string
	}
}

func DefaultConfig() Config {
	c := Config{}
	c.Network.GameAddr = ":8080"
	c.Server.Name = "Wonderful Hand game"
	c.Server.Multicore = true
	c.Server.ReusePort = true
	return c
}

func ReadConfig() (Config, error) {
	c := DefaultConfig()
	if _, err := os.Stat("./etc/config.toml"); os.IsNotExist(err) {
		data, err := toml.Marshal(c)
		if err != nil {
			return c, fmt.Errorf("failed encoding default config: %v", err)
		}
		if err := ioutil.WriteFile("config.toml", data, 0644); err != nil {
			return c, fmt.Errorf("failed creating config: %v", err)
		}
		return c, nil
	}
	data, err := ioutil.ReadFile("./etc/config.toml")
	if err != nil {
		return c, fmt.Errorf("error reading config: %v", err)
	}
	if err := toml.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("error decoding config: %v", err)
	}
	return c, nil
}
