package ch1

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Path string `envconfig:"PATH" required:"true"`
}

func envVars() {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		return
	}
	log.Println(c)
}
