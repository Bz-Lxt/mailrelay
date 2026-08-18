package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Dir  string
	Addr string
	Web  string
}

func FromEnv(base Config) Config {
	if v := os.Getenv(envPrefix() + "_DIR"); v != "" {
		base.Dir = v
	}
	if v := os.Getenv(envPrefix() + "_ADDR"); v != "" {
		base.Addr = v
	}
	if v := os.Getenv(envPrefix() + "_WEB"); v != "" {
		base.Web = v
	}
	return base
}

func Normalize(c Config) (Config, error) {
	if strings.TrimSpace(c.Dir) == "" {
		c.Dir = "data"
	}
	if strings.TrimSpace(c.Addr) == "" {
		c.Addr = ":8080"
	}
	if strings.TrimSpace(c.Web) == "" {
		c.Web = "web"
	}
	if !strings.Contains(c.Addr, ":") {
		return c, fmt.Errorf("addr must include port")
	}
	return c, nil
}

func envPrefix() string {
	return strings.ToUpper(strings.ReplaceAll("mailrelay", "-", "_"))
}
