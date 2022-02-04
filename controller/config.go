package controller

import (
	"log"
	"os"

	"github.com/spf13/pflag"
)

var (
	rosAddr   = pflag.StringP("api", "r", "", "RouterOS REST API address")
	rosUser   = pflag.StringP("user", "u", "admin", "RouterOS username")
	rosPasswd = pflag.StringP("password", "p", "", "RouterOS password")

	listFile = pflag.StringP("list", "f", "", "path to list-domain file")

	bindIP   = pflag.StringP("ip", "b", "0.0.0.0", "bind IP address")
	bindPort = pflag.IntP("port", "t", 5514, "bind port")
)

type Config struct {
	RouterOSAddr   string
	RouterOSUser   string
	RouterOSPasswd string

	// path to list domain file
	ListFile string

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
	if c.ListFile == "" {
		log.Fatal("Missing List-domain file")
	} else {
		if _, err := os.Stat(c.ListFile); err != nil {
			log.Fatal(err)
		}
	}
}

func FromParams() *Config {
	pflag.Parse()
	c := &Config{
		RouterOSAddr:      *rosAddr,
		RouterOSUser:      *rosUser,
		RouterOSPasswd:    *rosPasswd,
		ListFile:          *listFile,
		LogServerBindIP:   *bindIP,
		LogServerBindPort: *bindPort,
	}

	c.Validate()
	return c
}
