package main

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Configuration struct {
	Dbcfg     DBconfig `envconfig:"DB"`
	Servercfg ServerConfig
}

type ServerConfig struct {
	Port string `envconfig:"PORT" default:"8000"`
}
type DBconfig struct {
	User     string `required:"true" split_words:"true"`
	Name     string `required:"true" split_words:"true"`
	Password string `required:"true" split_words:"true"`
	SSLMode  string ` required:"true" default:"disable"`
}

func LoadConfig() Configuration {
	log.Println("loading configrations from enviroment variables")
	var cfg Configuration

	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf("failed to load enviroment variables with error %s", err.Error())
	}

	return cfg

}
