package config

import (
	"time"

	args "github.com/alexflint/go-arg"
	log "github.com/tarent/go-log-middleware/logging"
)

type Config struct {
	ValidationHost    string        `arg:"--validation-host,env" help:"Default is empty"`
	ValidationPort    int           `arg:"--validation-port,env" help:"Default is 9999"`
	ValidationTimeout time.Duration `arg:"--validation-timeout,env" help:"Default is 5s"`
}

func NewConfig() Config {
	config := Config{}
	config.ValidationPort = 9999
	timeout, err := time.ParseDuration("5s")
	if err != nil {
		log.Logger.Fatal(err)
	}
	config.ValidationTimeout = timeout
	args.MustParse(&config)
	return config
}
