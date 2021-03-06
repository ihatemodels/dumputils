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

package postgres

import (
	"bufio"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ihatemodels/dumputils/pkg/utils"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ihatemodels/dumputils/internal/log"
	"github.com/rotisserie/eris"

	_ "github.com/lib/pq"
)

type Database struct {
	Name              string
	Host              string
	Port              int
	Database          string
	Username          string
	Password          string
	Version           int
	DumpServer        bool
	DumpAll           bool
	Verbose           bool
	ExcludedDatabases []string
	connectionString  string
	pgDumpBin         string
	pgDumpAllBin      string
}

const binPath = "/usr/lib/postgresql/%d/bin/%s"

var ErrCanNotDeleteFile = eris.New("can not delete uncompressed file in pg_dumpall")

func (d *Database) buildConnectionString() {
	d.connectionString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		d.Host,
		d.Port,
		d.Username,
		d.Password,
		d.Database)
}

// Probe if the database is reachable and needed tools exist.
func (d *Database) Probe() error {

	d.buildConnectionString()

	if d.DumpServer {
		d.pgDumpAllBin = fmt.Sprintf(binPath, d.Version, "pg_dumpall")
		_, err := os.Stat(d.pgDumpAllBin)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return eris.Wrapf(err, "file %s doesn't not exist", d.pgDumpBin)
			}
			return eris.Wrapf(err, "checking pg_dumpall tool exist failed for instance %s", d.Name)
		}
	} else {
		d.pgDumpBin = fmt.Sprintf(binPath, d.Version, "pg_dump")
		_, err := os.Stat(d.pgDumpBin)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return eris.Wrapf(err, "file %s doesn't not exist", d.pgDumpBin)
			}
			return eris.Wrapf(err, "checking pg_dump tool exist failed for instance %s", d.Name)
		}
	}

	db, err := sql.Open("postgres", d.connectionString)

	if err != nil {
		return eris.Wrapf(err, "can not open database connection to instance: %s - host: %s", d.Name, d.Host)
	}

	if err := db.Ping(); err != nil {
		return eris.Wrapf(err, "database ping failed for instance: %s host: %s", d.Name, d.Host)
	}

	db.Close()

	return nil
}

func (d *Database) Dump() error {

	if err := d.Probe(); err != nil {
		return eris.Wrapf(err, "failed to probe database with name %s and host %s", d.Name, d.Host)
	}

	switch {
	case d.DumpAll:
		db, err := sql.Open("postgres", d.connectionString)

		if err != nil {
			return eris.Wrapf(err, "can not open database connection to instance: %s - host: %s", d.Name, d.Host)
		}

		rows, err := db.Query("SELECT datname FROM pg_database;")

		if err != nil {
			return eris.Wrapf(err, "failed to get all databases in dumpAll mode for instance: %s", d.Name)
		}

		var databases []string

		for rows.Next() {
			var currentDatabase string
			if err := rows.Scan(&currentDatabase); err != nil {
				return eris.Wrapf(err, "failed to scan all databases in dumpAll mode for instance: %s", d.Name)
			}

			if utils.Contains(d.ExcludedDatabases, currentDatabase) {
				log.Debugf("skipping database: %s in dump_all for instance: %s", currentDatabase, d.Name)
				continue
			}

			databases = append(databases, currentDatabase)
		}

		if err := rows.Err(); err != nil {
			return eris.Wrapf(err, "failed to scan all databases in dumpAll mode in instance: %s", d.Name)
		}

		rows.Close()

		log.Debugf("starting to dump: %v databases in instance: %s", databases, d.Name)

		for _, db := range databases {
			if err := d.pgDump(db); err != nil {
				return eris.Wrapf(err, "failed to pg_dump database: %s in instance: %s", d.Name)
			}
		}
	case d.DumpServer:
		if err := d.pgDumpServer(); err != nil {
			return eris.Wrapf(err, "failed to pg_dumpall instance %s", d.Name)
		}
	default:
		if err := d.pgDump(d.Database); err != nil {
			return eris.Wrapf(err, "failed to pg_dump database: %s in instance: %s", d.Name)
		}
	}

	return nil
}

// PgDump executes pg_dump with maximum level of compression.
func (d *Database) pgDump(database string) error {

	fileName := fmt.Sprintf("%s-%s-%s.dump", d.Name, database, time.Now().Format("2006-01-02-15-04-05"))

	args := []string{"-h", d.Host, "-p", strconv.Itoa(d.Port), "-U", d.Username, "-O", "-Fc", "-Z", "9", database, "-f", fileName}

	if d.Verbose {
		args = append([]string{"-v"}, args...)
	}

	cmd := exec.Command(d.pgDumpBin, args...)

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", d.Password))

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	log.Infof("now executing: %s", cmd.String())

	if err := cmd.Run(); err != nil {
		return eris.Wrapf(err, "command %s failed: %s", cmd.String(), out.String())
	}

	if d.Verbose {
		text := strings.Split(out.String(), "\n")
		for _, line := range text {
			log.Infof("pg_dump: %s", line)
		}
	}

	return nil
}

func (d *Database) pgDumpServer() error {

	fileName := fmt.Sprintf("%s-%s.server.dump", d.Name, time.Now().Format("2006-01-02-15-04-05"))

	args := []string{"-h", d.Host, "-p", strconv.Itoa(d.Port), "-U", d.Username, "-f", fileName}

	if d.Verbose {
		args = append([]string{"-v"}, args...)
	}

	cmd := exec.Command(d.pgDumpAllBin, args...)

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("PGPASSWORD=%s", d.Password))

	// capture only stderr cuz pgtools outputs there always
	stderr, err := cmd.StderrPipe()

	if err != nil {
		return eris.Wrapf(err, "command %s failed: %s", cmd.String())
	}

	if err := cmd.Start(); err != nil {
		return eris.Wrapf(err, "command %s failed: %s", cmd.String())
	}

	if d.Verbose {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			log.Infof("pg_dumpall: %s", scanner.Text())
		}
	}

	if err := utils.Gzip(fileName); err != nil {
		return eris.Wrapf(err, "can not compress file: %s", fileName)
	}

	if err := os.Remove(fileName); err != nil {
		return ErrCanNotDeleteFile
	}
	_ = cmd.Wait()

	return nil
}
