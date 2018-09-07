package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	defaultHost = ""
	defaultPort = "8080"
	defaultDSN  = "postgres://hellofresh:hellofresh@localhost:5432/hellofresh?sslmode=disable"
)

var commandlineConfig = newCommandLineConfig()

func newCommandLineConfig() *viper.Viper {
	v := viper.New()
	defineFlags()
	if !pflag.Parsed() {
		pflag.Parse()
	}
	loadCommandLineConfig(v)
	return v
}

func defineFlags() {
	pflag.String("host", defaultHost, "host that the http service binds to")
	pflag.String("port", defaultPort, "port that the http service listens to")
	pflag.String("dsn", defaultDSN, "postgreSQL database connection string")
}

func loadCommandLineConfig(v *viper.Viper) {
	v.BindPFlags(pflag.CommandLine)
}

type applicationConfig struct {
	host string
	port string
	dsn  string
}

func newApplicationConfig() *applicationConfig {
	return &applicationConfig{
		host: defaultHost,
		port: defaultPort,
		dsn:  defaultDSN,
	}
}

func (c *applicationConfig) bind(v *viper.Viper) {
	c.host = v.GetString("host")
	c.port = v.GetString("port")
	c.dsn = v.GetString("dsn")
}
