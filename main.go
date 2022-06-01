package main

import (
	"fmt"
	"github.com/ihatemodels/pgtools/app"
	"github.com/ihatemodels/pgtools/internal/log"
	"os"

	"github.com/ihatemodels/pgtools/internal/config"
)

func main() {
	if err := config.Init(os.Getenv("PGTOOLS_CONFIG_PATH")); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "can not build application config: %v", err)
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}

	if config.App.Log.Type == "human" {
		log.SetHumanLogger(config.App.Log.Level)
	}

	log.Infof("pgtools started")

	app.Run()
}
