/*
 * Copyright (c) 2022. Dumputils Authors
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package dumputils

import (
	"github.com/ihatemodels/dumputils/internal/config"
	"github.com/ihatemodels/dumputils/internal/log"
	"github.com/ihatemodels/dumputils/pkg/dump/postgres"
	"github.com/rotisserie/eris"
	"os"
)

func Run() error {

	cfg, err := config.New(os.Getenv("DUMPUTILS_CONFIG_PATH"))

	if err != nil {
		return err
	}

	log.Init(cfg.Log.Level, cfg.Log.Type)

	log.Infof("dumputils started... configuration validated")

	for _, instance := range cfg.Databases {
		if !instance.Enabled {
			log.Warnf("instance %s is defined in the configuration but not enabled", instance.Name)
			continue
		}
		db := postgres.Database{
			Name:              instance.Name,
			Host:              instance.Host,
			Password:          instance.Password,
			Port:              instance.Port,
			Username:          instance.Username,
			Database:          instance.Database,
			DumpServer:        instance.DumpServer,
			DumpAll:           instance.DumpAll,
			Version:           instance.Version,
			Verbose:           instance.Verbose,
			ExcludedDatabases: instance.ExcludeDatabasesSlice,
		}

		if err := db.Dump(); err != nil {
			if eris.Is(err, postgres.ErrCanNotDeleteFile) {
				continue
			}
			log.Errorf(err, "failed to dump instance %s running on host %s", instance.Name, instance.Host)
			continue
		}
	}
	return nil
}
