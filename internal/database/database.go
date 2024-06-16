package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/astrorick/seekret/internal/server"
	"github.com/astrorick/seekret/pkg/version"
)

/*
This function should be executed as a part of the server initialization procedure.
It is meant to run preliminary consistency checks on the provided database in order to initialize missing tables and set some default values.
*/
func (srv *server.Server) RunPreliminaryChecks() error {
	// checks
	var statsRow *sql.Row
	var usersRow *sql.Row

	// assess database status
	switch srv.Config.DatabaseType {
	case "sqlite3":
		statsRow = srv.Database.QueryRow("SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'stats'")
		usersRow = srv.Database.QueryRow("SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'users'")
	default:
		return errors.New("unsupported database type")
	}

	// check the 'stats' table
	if err := statsRow.Scan(); errors.Is(err, sql.ErrNoRows) {
		// generate empty 'stats' table
		if _, err := srv.Database.Exec("CREATE TABLE stats (id INTEGER PRIMARY KEY, version TEXT NOT NULL)"); err != nil {
			return err
		}

		// write default values
		if _, err := srv.Database.Exec("INSERT INTO stats (version) VALUES (?)", srv.Version.String()); err != nil {
			return err
		}
	} else {
		// read database version from 'stats' table
		var databaseVersionString string
		srv.Database.QueryRow("SELECT version FROM stats").Scan(&databaseVersionString)

		// parse to 'Version' object
		databaseVersion, err := version.New(databaseVersionString)
		if err != nil {
			return err
		}

		// check versions mismatch
		if srv.Version.OlderThan(databaseVersion) {
			// database was created with a newer server version
			return fmt.Errorf("outdated server version (%s) for the provided database (%s)", srv.Version.String(), databaseVersion.String())
		}
		if srv.Version.NewerThan(databaseVersion) {
			// TODO: proceed to migration
			//return nil

			fmt.Printf("Database updated from version %s to version %s\n", databaseVersion.String(), srv.Version.String())
		}
	}

	// check the 'users' table
	if err := usersRow.Scan(); errors.Is(err, sql.ErrNoRows) {
		// generate empty 'users' table
		if _, err := srv.Database.Exec("CREATE TABLE users (id INTEGER PRIMARY KEY, username TEXT NOT NULL, salt TEXT NOT NULL, verifier TEXT NOT NULL)"); err != nil {
			return err
		}
	}

	return nil
}
