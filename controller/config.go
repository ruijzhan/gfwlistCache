package controller

import (
	"log"

	"github.com/spf13/pflag"
)

var (
	rosAddr   = pflag.StringP("api", "r", "", "RouterOS REST API address")
	rosUser   = pflag.StringP("user", "u", "admin", "RouterOS username")
	rosPasswd = pflag.StringP("password", "p", "", "RouterOS password")

	bindIP   = pflag.StringP("ip", "b", "0.0.0.0", "bind IP address")
	bindPort = pflag.IntP("port", "t", 5514, "bind port")
)

type Config struct {
	RouterOSAddr   string
	RouterOSUser   string
	RouterOSPasswd string

	LogServerBindIP   string
	LogServerBindPort int
}

func (c *Config) Validate() {
	if c.RouterOSAddr == "" {
		log.Fatal("Missing RouterOS API address")
	}
	if c.RouterOSUser == "" {
		log.Fatal("Missing RouterOS Username")
	}
	if c.RouterOSPasswd == "" {
		log.Fatal("Missing RouterOS Password")
	}
}

func FromParams() *Config {
	pflag.Parse()
	c := &Config{
		RouterOSAddr:      *rosAddr,
		RouterOSUser:      *rosUser,
		RouterOSPasswd:    *rosPasswd,
		LogServerBindIP:   *bindIP,
		LogServerBindPort: *bindPort,
	}

	c.Validate()
	return c
}
