package main

import (
	"fmt"
	"github.com/ihatemodels/dumputils/internal/log"
	"github.com/ihatemodels/dumputils/pkg/postgres"
	"github.com/rotisserie/eris"
	"os"

	"github.com/ihatemodels/dumputils/internal/config"
)

func main() {
	if err := config.Init(os.Getenv("DUMPUTILS_CONFIG_PATH")); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "can not build application config: %v", eris.ToString(err, true))
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}

	if config.App.Log.Type == "human" {
		log.SetHumanLogger(config.App.Log.Level)
	}

	log.Infof("pgtools started")

	for _, instance := range config.App.Databases {
		fmt.Println(instance.ExcludeDatabasesSlice)
		db := postgres.Database{
			Name:              instance.Name,
			Host:              instance.Host,
			Password:          instance.Password,
			Port:              instance.Port,
			Username:          instance.Username,
			Database:          instance.Database,
			IsServer:          instance.DumpServer,
			DumpAll:           instance.DumpAll,
			Version:           instance.Version,
			Verbose:           instance.Verbose,
			ExcludedDatabases: instance.ExcludeDatabasesSlice,
		}

		if err := db.Dump(); err != nil {
			log.Errorf(err, "failed to dump database with name %s and host %s", instance.Name, instance.Host)
			continue
		}
	}
}
