package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	defaultHost = ""
	defaultPort = "8080"
	defaultDSN  = ""
)

const noDefaultValue = ""

func newCommandLineConfig() *viper.Viper {
	v := viper.New()
	defineFlags()
	if !pflag.Parsed() {
		pflag.Parse()
	}
	loadCommandLineFlag(v, pflag.CommandLine)
	return v
}

func defineFlags() {
	pflag.String("host", noDefaultValue, "host that the http service binds to")
	pflag.String("port", noDefaultValue, "port that the http service listens to")
	pflag.String("dsn", noDefaultValue, "postgreSQL database connection string")
}

func loadCommandLineFlag(v *viper.Viper, flagSet *pflag.FlagSet) {
	flagSet.VisitAll(func(flag *pflag.Flag) {
		if !flag.Changed {
			return
		}
		if err := v.BindPFlag(flag.Name, flag); err != nil {
			return
		}
	})
}

func newEnvironmentVariableConfig() *viper.Viper {
	v := viper.New()
	loadEnvironmentVariables(v)
	return v
}

func loadEnvironmentVariables(v *viper.Viper) {
	envs := []string{"HOST", "PORT", "DSN"}
	for _, env := range envs {
		if err := v.BindEnv(env); err != nil {
			panic(err)
		}
	}
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
	if v.IsSet("host") {
		c.host = v.GetString("host")
	}
	if v.IsSet("port") {
		c.port = v.GetString("port")
	}
	if v.IsSet("dsn") {
		c.dsn = v.GetString("dsn")
	}
}
