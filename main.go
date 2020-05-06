package main

import (
	"context"
	"log"

	"github.com/ironcore864/sap-bt/config"
	"github.com/ironcore864/sap-bt/server"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	// load config from env vars
	var conf config.Config
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Fatal(err.Error())
	}

	// background context
	ctx := context.Background()

	// start server
	err = server.Start(ctx, &conf)
	log.Fatal(err)
}
