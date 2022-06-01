package app

import (
	"github.com/ihatemodels/pgtools/internal/config"
	"github.com/ihatemodels/pgtools/internal/log"
	"github.com/ihatemodels/pgtools/internal/pgtool"
)

func Run() {
	for _, instance := range config.App.Databases {
		db := pgtool.Database{
			Name:     instance.Name,
			Host:     instance.Host,
			Password: instance.Password,
			Port:     instance.Port,
			Username: instance.Username,
			Database: instance.Database,
			Version:  instance.Version,
			Verbose:  instance.Verbose,
		}

		if err := db.Probe(); err != nil {
			log.Errorf(err, "failed to probe database with name %s and host %s", instance.Name, instance.Host)
			continue
		}

		if _, err := db.PgDump(); err != nil {
			log.Errorf(err, "failed to create backup of database with name %s on host %s", instance.Name, instance.Host)
			continue
		}
	}

}
