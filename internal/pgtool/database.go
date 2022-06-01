package pgtool

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/ihatemodels/pgtools/internal/log"
	"github.com/rotisserie/eris"

	_ "github.com/lib/pq"
)

type Database struct {
	Name         string
	Host         string
	Port         int
	Database     string
	Username     string
	Password     string
	Version      int
	IsServer     bool
	Verbose      bool
	pgDumpBin    string
	pgDumpAllBin string
}

const binPath = "/usr/lib/postgresql/%d/bin/%s"

// Probe if the database is reachable and needed tools exist.
func (d *Database) Probe() error {

	if d.IsServer {
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

	connection := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s",
		d.Host,
		d.Port,
		d.Username,
		d.Password,
		d.Database)

	db, err := sql.Open("postgres", connection)

	if err != nil {
		return eris.Wrapf(err, "can not open database connection to instance: %s - host: %s", d.Name, d.Host)
	}

	if err := db.Ping(); err != nil {
		return eris.Wrapf(err, "database ping failed for instance: %s host: %s", d.Name, d.Host)
	}

	db.Close()

	return nil
}

// PgDump executes pg_dump with maximum level of compression.
// Returns pg_dump stdout+stderr and errors if pg_dump returns exit code != 0
func (d *Database) PgDump() (string, error) {

	fileName := fmt.Sprintf("%s-%s-%s.dump", d.Name, d.Database, time.Now().Format("2006-01-02-15-04-05-000000"))

	args := []string{"-h", d.Host, "-p", strconv.Itoa(d.Port), "-U", d.Username, "-O", "-Fc", "-Z", "9", d.Database, "-f", fileName}

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
		return "", eris.Wrapf(err, "command %s failed: %s", cmd.String(), out.String())
	}

	if d.Verbose {
		log.Infof("output %s", out.String())
	}

	return out.String(), nil
}

func (d *Database) PgDumpAll() error {
	return nil
}
