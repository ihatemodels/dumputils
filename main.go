/*
 * Copyright (c) 2022.  Dumputils Authors
 *
 * https://opensource.org/licenses/MIT
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
 *
 */

package main

import (
	"fmt"
	"github.com/ihatemodels/dumputils/internal/config"
	"github.com/ihatemodels/dumputils/internal/log"
	"github.com/ihatemodels/dumputils/pkg/postgres"
	"github.com/rotisserie/eris"
	"os"
)

func main() {
	if err := config.Init(os.Getenv("DUMPUTILS_CONFIG_PATH")); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "can not build application config: %v", eris.ToString(err, true))
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}

	log.Configure()

	log.Infof("dumputils: configuration validated")

	//var notifierChannels []notifiers.Notifier
	//
	//if config.App.Notifiers.Slack.Enabled {
	//	slackNotifier, err := slack.New(config.App.Notifiers.Slack.BotToken, config.App.Notifiers.Slack.Channel)
	//	if err != nil {
	//		log.Errorf(err, "slack notifier failed to initialize")
	//	} else {
	//		notifierChannels = append(notifierChannels, slackNotifier)
	//		err := slackNotifier.Notify(notifiers.Input{
	//			InstanceName: "test",
	//			InstanceType: "PostgreSQL test",
	//			Duration:     1 * time.Hour,
	//			State:        notifiers.Failed,
	//		})
	//		if err != nil {
	//			log.Errorf(err, "notifier failed to send notification")
	//		}
	//	}
	//}

	for _, instance := range config.App.Databases {
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
			log.Errorf(err, "failed to dump instance %s running on host %s", instance.Name, instance.Host)
			continue
		}
	}
}
