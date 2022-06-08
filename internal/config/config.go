/*
 * Copyright (c) 2022. Dumputils Authors
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

package config

import (
	"fmt"
	"gopkg.in/dealancer/validate.v2"
	"os"
	"strings"

	"github.com/rotisserie/eris"
	"gopkg.in/yaml.v2"
)

// New returns validated Settings instance. Errors on file opening or configuration miss sense.
func New(filePath string) (Settings, error) {
	var out Settings
	if len(filePath) == 0 {
		return out, eris.New("internal/config: filePath is empty")
	}

	f, err := os.Open(filePath)

	if err != nil {
		return out, eris.Wrap(err, "internal/config: can not open file")
	}

	defer f.Close()

	d := yaml.NewDecoder(f)

	if err := d.Decode(&out); err != nil {
		return out, eris.Wrap(err, "internal/config: can not decode file to yaml struct")
	}

	if err := validate.Validate(&out); err != nil {
		return out, eris.Wrap(err, "internal/config: the provided configuration is invalid")
	}

	for i, instance := range out.Databases {
		if instance.DumpAll && instance.DumpServer {
			return out, eris.New(fmt.Sprintf("internal/config: dumpAll and dumpServer "+
				"flags can not be used together in database name: %s", instance.Name))
		}
		if !instance.DumpAll && !instance.DumpServer {
			if instance.Database == "" {
				return out, eris.New(fmt.Sprintf("internal/config: Database field can not "+
					"be empty in single dump mode for instance: %s", instance.Name))
			}
		}

		if instance.DumpAll {

			// Do not dump the system templates.
			// See: https://www.postgresql.org/docs/current/manage-ag-templatedbs.html
			out.Databases[i].ExcludeDatabasesSlice = []string{"template0", "template1"}

			split := strings.Split(instance.ExcludeDatabases, ",")

			if len(split) > 0 {
				for _, db := range split {
					out.Databases[i].ExcludeDatabasesSlice = append(
						out.Databases[i].ExcludeDatabasesSlice,
						strings.ReplaceAll(db, " ", ""),
					)
				}
			}
		}
	}

	return out, nil
}

type Settings struct {
	Log struct {
		Type  string `yaml:"type"  validate:"one_of=human,json"`
		Level string `yaml:"level" validate:"one_of=debug,info,warning,error"`
	} `yaml:"log"`

	Databases []struct {
		Host                  string `yaml:"host" validate:"empty=false"`
		Enabled               bool   `yaml:"enabled"`
		Name                  string `yaml:"name" validate:"empty=false"`
		Port                  int    `yaml:"port" validate:"ne=0"`
		Database              string `yaml:"database"`
		Username              string `yaml:"username" validate:"empty=false"`
		Password              string `yaml:"password" validate:"empty=false"`
		Version               int    `yaml:"version"  validate:"one_of=10,11,12,13,14"`
		Verbose               bool   `yaml:"verbose"`
		DumpAll               bool   `yaml:"dumpAll"`
		ExcludeDatabases      string `yaml:"excludeDatabases"`
		ExcludeDatabasesSlice []string
		DumpServer            bool `yaml:"dumpServer"`
	} `yaml:"databases"`

	Outputs struct {
		TmpDir string `yaml:"tmpDir" validate:"empty=false"`
		Minio  struct {
			Enabled         bool   `yaml:"enabled"`
			Endpoint        string `yaml:"endpoint"`
			AccessKeyID     string `yaml:"accessKeyID"`
			BucketName      string `yaml:"bucketName"`
			SecretAccessKey string `yaml:"secretAccessKey"`
		} `yaml:"minio"`
		Sftp struct {
			Enabled   bool   `yaml:"enabled"`
			Host      string `yaml:"host"`
			Port      int    `yaml:"port"`
			User      string `yaml:"user"`
			Password  string `yaml:"password"`
			Directory string `yaml:"directory"`
		} `yaml:"sftp"`
		Filesystem struct {
			Enabled bool   `yaml:"enabled"`
			Path    string `yaml:"path"`
		} `yaml:"filesystem"`
	} `yaml:"outputs"`

	Notifiers struct {
		Email struct {
			Enabled       bool   `yaml:"enabled"`
			SMTP          string `yaml:"smtp"`
			Port          int    `yaml:"port"`
			Sender        string `yaml:"sender"`
			Password      string `yaml:"password"`
			SendOnSuccess bool   `yaml:"sendOnSuccess"`
		} `yaml:"email"`
		Slack struct {
			Verbose       bool   `yaml:"verbose"`
			Enabled       bool   `yaml:"enabled"`
			BotToken      string `yaml:"botToken" validate:"empty=false"`
			Channel       string `yaml:"channel"  validate:"empty=false"`
			SendOnSuccess bool   `yaml:"sendOnSuccess"`
		} `yaml:"slack"`
	} `yaml:"notifiers"`
}
